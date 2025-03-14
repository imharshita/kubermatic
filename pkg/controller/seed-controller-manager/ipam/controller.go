/*
Copyright 2022 The Kubermatic Kubernetes Platform contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ipam

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	kubermaticv1 "k8c.io/kubermatic/v2/pkg/apis/kubermatic/v1"
	kubermaticv1helper "k8c.io/kubermatic/v2/pkg/apis/kubermatic/v1/helper"
	kuberneteshelper "k8c.io/kubermatic/v2/pkg/kubernetes"
	"k8c.io/kubermatic/v2/pkg/provider"
	"k8c.io/kubermatic/v2/pkg/version/kubermatic"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/record"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	ControllerName = "kkp-ipam-controller"
)

// Reconciler stores all components required for the IPAM controller.
type Reconciler struct {
	ctrlruntimeclient.Client

	workerName   string
	log          *zap.SugaredLogger
	configGetter provider.KubermaticConfigurationGetter
	recorder     record.EventRecorder
	versions     kubermatic.Versions
}

// Add creates a new IPAM controller.
func Add(
	mgr manager.Manager,
	log *zap.SugaredLogger,
	workerName string,
	configGetter provider.KubermaticConfigurationGetter,
	versions kubermatic.Versions,
) error {
	log = log.Named(ControllerName)

	reconciler := &Reconciler{
		Client:       mgr.GetClient(),
		workerName:   workerName,
		log:          log,
		recorder:     mgr.GetEventRecorderFor(ControllerName),
		configGetter: configGetter,
		versions:     versions,
	}

	c, err := controller.New(ControllerName, mgr, controller.Options{
		Reconciler:              reconciler,
		MaxConcurrentReconciles: 1, // this should be 1 to avoid race condition
	})
	if err != nil {
		return fmt.Errorf("failed to create controller: %w", err)
	}

	if err := c.Watch(&source.Kind{Type: &kubermaticv1.Cluster{}}, &handler.EnqueueRequestForObject{}); err != nil {
		return fmt.Errorf("failed to create watch for clusters: %w", err)
	}

	enqueueClustersForIPAMPool := handler.EnqueueRequestsFromMapFunc(func(a ctrlruntimeclient.Object) []reconcile.Request {
		ipamPool := a.(*kubermaticv1.IPAMPool)

		clusterList := &kubermaticv1.ClusterList{}
		if err := mgr.GetClient().List(context.Background(), clusterList); err != nil {
			utilruntime.HandleError(fmt.Errorf("failed to list Clusters: %w", err))
			log.Errorw("Failed to list clusters", zap.Error(err))
			return []reconcile.Request{}
		}

		requests := []reconcile.Request{}
		for _, cluster := range clusterList.Items {
			_, isClusterDCConfiguredInPool := ipamPool.Spec.Datacenters[cluster.Spec.Cloud.DatacenterName]
			if !isClusterDCConfiguredInPool {
				// The IPAM pool is not relevant to this cluster, so skip it
				continue
			}
			requests = append(requests, reconcile.Request{NamespacedName: types.NamespacedName{Name: cluster.Name}})
		}
		return requests
	})
	if err := c.Watch(&source.Kind{Type: &kubermaticv1.IPAMPool{}}, enqueueClustersForIPAMPool); err != nil {
		return fmt.Errorf("failed to create watch for IPAM Pools: %w", err)
	}

	return nil
}

func (r *Reconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.With("request", request)
	log.Debug("Processing")

	cluster := &kubermaticv1.Cluster{}
	if err := r.Get(ctx, request.NamespacedName, cluster); err != nil {
		if apierrors.IsNotFound(err) {
			log.Debug("Skipping because the cluster is already gone")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if cluster.Status.NamespaceName == "" {
		log.Debug("Cluster has no namespace name yet, skipping")
		return reconcile.Result{}, nil
	}

	if cluster.DeletionTimestamp != nil {
		log.Debug("Cluster is in deletion, skipping")
		return reconcile.Result{}, nil
	}

	// Add a wrapping here so we can emit an event on error
	result, err := kubermaticv1helper.ClusterReconcileWrapper(
		ctx,
		r.Client,
		r.workerName,
		cluster,
		r.versions,
		kubermaticv1.ClusterConditionIPAMControllerReconcilingSuccess,
		func() (*reconcile.Result, error) {
			return r.reconcile(ctx, cluster)
		},
	)
	if err != nil {
		log.Errorw("Reconciling failed", zap.Error(err))
		r.recorder.Event(cluster, corev1.EventTypeWarning, "ReconcilingError", err.Error())
	}
	if result == nil {
		result = &reconcile.Result{}
	}
	return *result, err
}

func (r *Reconciler) reconcile(ctx context.Context, cluster *kubermaticv1.Cluster) (*reconcile.Result, error) {
	// List IPAM Pools
	ipamPoolList := &kubermaticv1.IPAMPoolList{}
	err := r.Client.List(ctx, ipamPoolList)
	if err != nil {
		return nil, fmt.Errorf("failed to list IPAM pools: %w", err)
	}

	// Loop IPAM pools, considering only the relevant ones (i.e. IPAM pools with same cluster datacenter)
	for _, ipamPool := range ipamPoolList.Items {
		clusterDC := cluster.Spec.Cloud.DatacenterName

		dcIPAMPoolCfg, isClusterDCConfigured := ipamPool.Spec.Datacenters[clusterDC]

		ipamAllocation := &kubermaticv1.IPAMAllocation{}
		err = r.Client.Get(ctx, types.NamespacedName{Namespace: cluster.Status.NamespaceName, Name: ipamPool.Name}, ipamAllocation)
		if err == nil {
			if !isClusterDCConfigured {
				// There is an allocation for a datacenter that is not present
				// in the IPAM pool spec anymore, so delete it
				if err := r.Delete(ctx, ipamAllocation); err != nil && !apierrors.IsNotFound(err) {
					return nil, err
				}
			}
			// Skip because no allocation from this IPAM pool is needed for the cluster
			continue
		} else if !apierrors.IsNotFound(err) {
			return nil, err
		}

		if !isClusterDCConfigured {
			// This IPAM pool is not relevant to cluster, so skip it
			continue
		}

		dcIPAMPoolUsageMap, err := r.compileCurrentAllocationsForPoolInDatacenter(ctx, ipamPool.Name, clusterDC, dcIPAMPoolCfg)
		if err != nil {
			return nil, err
		}

		err = r.generateNewClusterAllocationForPool(ctx, cluster, &ipamPool, dcIPAMPoolCfg, dcIPAMPoolUsageMap)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Reconciler) compileCurrentAllocationsForPoolInDatacenter(ctx context.Context, ipamPoolName, dc string, dcIPAMPoolCfg kubermaticv1.IPAMPoolDatacenterSettings) (sets.Set[string], error) {
	dcIPAMPoolUsageMap := sets.New[string]()

	// Check for exclusions in the configuration to mark them as "not free"
	switch dcIPAMPoolCfg.Type {
	case kubermaticv1.IPAMPoolAllocationTypeRange:
		ipsToExclude, err := getIPsFromAddressRanges(dcIPAMPoolCfg.ExcludeRanges)
		if err != nil {
			return nil, err
		}
		for _, ipToExclude := range ipsToExclude {
			dcIPAMPoolUsageMap.Insert(ipToExclude)
		}
	case kubermaticv1.IPAMPoolAllocationTypePrefix:
		for _, subnetCIDRToExclude := range dcIPAMPoolCfg.ExcludePrefixes {
			dcIPAMPoolUsageMap.Insert(string(subnetCIDRToExclude))
		}
	}

	// List all IPAM allocations
	ipamAllocationList := &kubermaticv1.IPAMAllocationList{}
	err := r.Client.List(ctx, ipamAllocationList)
	if err != nil {
		return nil, fmt.Errorf("failed to list IPAM allocations: %w", err)
	}

	// Iterate current IPAM allocations to build a map of used IPs (for range allocation type)
	// or used subnets (for prefix allocation type) per datacenter pool
	for _, ipamAllocation := range ipamAllocationList.Items {
		if ipamAllocation.Name != ipamPoolName || ipamAllocation.Spec.DC != dc {
			// This allocation is not relevant for this IPAM Pool, so skip it
			continue
		}

		switch ipamAllocation.Spec.Type {
		case kubermaticv1.IPAMPoolAllocationTypeRange:
			currentAllocatedIPs, err := getIPsFromAddressRanges(ipamAllocation.Spec.Addresses)
			if err != nil {
				return nil, err
			}
			// check if the current allocation is compatible with the IPAMPool being applied
			err = checkRangeAllocation(currentAllocatedIPs, string(dcIPAMPoolCfg.PoolCIDR), dcIPAMPoolCfg.AllocationRange)
			if err != nil {
				return nil, err
			}
			for _, ip := range currentAllocatedIPs {
				dcIPAMPoolUsageMap.Insert(ip)
			}
		case kubermaticv1.IPAMPoolAllocationTypePrefix:
			// check if the current allocation is compatible with the IPAMPool being applied
			err := checkPrefixAllocation(string(ipamAllocation.Spec.CIDR), string(dcIPAMPoolCfg.PoolCIDR), dcIPAMPoolCfg.AllocationPrefix)
			if err != nil {
				return nil, err
			}
			dcIPAMPoolUsageMap.Insert(string(ipamAllocation.Spec.CIDR))
		}
	}

	return dcIPAMPoolUsageMap, nil
}

func (r *Reconciler) generateNewClusterAllocationForPool(ctx context.Context, cluster *kubermaticv1.Cluster, ipamPool *kubermaticv1.IPAMPool, dcIPAMPoolCfg kubermaticv1.IPAMPoolDatacenterSettings, dcIPAMPoolUsageMap sets.Set[string]) error {
	newClustersAllocation := &kubermaticv1.IPAMAllocation{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Status.NamespaceName,
			Name:      ipamPool.Name,
			Labels:    map[string]string{},
		},
		Spec: kubermaticv1.IPAMAllocationSpec{
			Type: dcIPAMPoolCfg.Type,
			DC:   cluster.Spec.Cloud.DatacenterName,
		},
	}
	kuberneteshelper.EnsureUniqueOwnerReference(newClustersAllocation, metav1.OwnerReference{
		APIVersion: kubermaticv1.SchemeGroupVersion.String(),
		Kind:       kubermaticv1.IPAMPoolKindName,
		UID:        ipamPool.GetUID(),
		Name:       ipamPool.Name,
	})

	switch dcIPAMPoolCfg.Type {
	case kubermaticv1.IPAMPoolAllocationTypeRange:
		addresses, err := findFirstFreeRangesOfPool(ipamPool.Name, string(dcIPAMPoolCfg.PoolCIDR), dcIPAMPoolCfg.AllocationRange, dcIPAMPoolUsageMap)
		if err != nil {
			return err
		}
		newClustersAllocation.Spec.Addresses = addresses
	case kubermaticv1.IPAMPoolAllocationTypePrefix:
		subnetCIDR, err := findFirstFreeSubnetOfPool(ipamPool.Name, string(dcIPAMPoolCfg.PoolCIDR), dcIPAMPoolCfg.AllocationPrefix, dcIPAMPoolUsageMap)
		if err != nil {
			return err
		}
		newClustersAllocation.Spec.CIDR = kubermaticv1.SubnetCIDR(subnetCIDR)
	}

	err := r.Create(ctx, newClustersAllocation)
	if err != nil {
		return fmt.Errorf("failed to create IPAM Pool Allocation for IPAM Pool %s in cluster %s: %w", ipamPool.Name, cluster.Name, err)
	}

	return nil
}
