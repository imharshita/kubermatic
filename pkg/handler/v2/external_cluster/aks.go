/*
Copyright 2021 The Kubermatic Kubernetes Platform contributors.

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

package externalcluster

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerservice/mgmt/containerservice"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"

	clusterv1alpha1 "github.com/kubermatic/machine-controller/pkg/apis/cluster/v1alpha1"
	providerconfig "github.com/kubermatic/machine-controller/pkg/providerconfig/types"
	apiv1 "k8c.io/kubermatic/v2/pkg/api/v1"
	apiv2 "k8c.io/kubermatic/v2/pkg/api/v2"
	kubermaticv1 "k8c.io/kubermatic/v2/pkg/apis/kubermatic/v1"
	"k8c.io/kubermatic/v2/pkg/handler/v1/common"
	kuberneteshelper "k8c.io/kubermatic/v2/pkg/kubernetes"
	"k8c.io/kubermatic/v2/pkg/provider"
	"k8c.io/kubermatic/v2/pkg/provider/cloud/aks"
	"k8c.io/kubermatic/v2/pkg/resources"
	"k8c.io/kubermatic/v2/pkg/util/errors"
)

const (
	AKSNodepoolNameLabel = "kubernetes.azure.com/agentpool"
	AgentPoolModeSystem  = "System"
)

func createNewAKSCluster(ctx context.Context, aksCloudSpec *apiv2.AKSCloudSpec) error {
	aksClient, err := aks.GetAKSClusterClient(resources.AKSCredentials{
		TenantID:       aksCloudSpec.TenantID,
		ClientID:       aksCloudSpec.ClientID,
		SubscriptionID: aksCloudSpec.SubscriptionID,
		ClientSecret:   aksCloudSpec.ClientSecret,
	})
	if err != nil {
		return err
	}

	clusterSpec := aksCloudSpec.ClusterSpec
	agentPoolProfiles := clusterSpec.MachineDeploymentSpec
	basicSettings := agentPoolProfiles.BasicSettings
	optionalSettings := agentPoolProfiles.OptionalSettings
	clusterToCreate := containerservice.ManagedCluster{
		Name:     to.StringPtr(aksCloudSpec.Name),
		Location: to.StringPtr(clusterSpec.Location),
		ManagedClusterProperties: &containerservice.ManagedClusterProperties{
			DNSPrefix:         to.StringPtr(aksCloudSpec.Name),
			KubernetesVersion: to.StringPtr(clusterSpec.KubernetesVersion),
			ServicePrincipalProfile: &containerservice.ManagedClusterServicePrincipalProfile{
				ClientID: to.StringPtr(aksCloudSpec.ClientID),
				Secret:   to.StringPtr(aksCloudSpec.ClientSecret),
			},
		},
	}
	agentPoolProfilesToBeCreated := containerservice.ManagedClusterAgentPoolProfile{
		Name:              to.StringPtr(agentPoolProfiles.Name),
		VMSize:            to.StringPtr(basicSettings.VMSize),
		Count:             to.Int32Ptr(basicSettings.Count),
		Mode:              (containerservice.AgentPoolMode)(basicSettings.Mode),
		OsDiskSizeGB:      to.Int32Ptr(basicSettings.OsDiskSizeGB),
		AvailabilityZones: &basicSettings.AvailabilityZones,
		EnableAutoScaling: to.BoolPtr(basicSettings.EnableAutoScaling),
		MaxCount:          to.Int32Ptr(basicSettings.ScalingConfig.MaxCount),
		MinCount:          to.Int32Ptr(basicSettings.ScalingConfig.MinCount),
		NodeLabels:        optionalSettings.NodeLabels,
	}

	clusterToCreate.AgentPoolProfiles = &[]containerservice.ManagedClusterAgentPoolProfile{agentPoolProfilesToBeCreated}

	_, err = aksClient.CreateOrUpdate(
		ctx,
		aksCloudSpec.ResourceGroup,
		aksCloudSpec.Name,
		clusterToCreate,
	)
	if err != nil {
		return err
	}

	return nil
}

func checkCreatePoolReqValidity(aksMD *apiv2.AKSMachineDeploymentCloudSpec) error {
	if aksMD == nil {
		return errors.NewBadRequest("AKS MachineDeploymentSpec cannot be nil")
	}
	basicSettings := aksMD.BasicSettings
	// check whether required fields for nodepool creation are provided
	fields := reflect.ValueOf(&basicSettings).Elem()
	for i := 0; i < fields.NumField(); i++ {
		yourjsonTags := fields.Type().Field(i).Tag.Get("required")
		if strings.Contains(yourjsonTags, "true") && fields.Field(i).IsZero() {
			return errors.NewBadRequest("required field is missing: %v", fields.Type().Field(i).Name)
		}
	}
	if basicSettings.EnableAutoScaling {
		maxCount := basicSettings.ScalingConfig.MaxCount
		minCount := basicSettings.ScalingConfig.MinCount
		if maxCount == 0 {
			return errors.NewBadRequest("InvalidParameter: value of maxCount for enabled autoscaling is invalid")
		}
		if minCount == 0 {
			return errors.NewBadRequest("InvalidParameter: value of minCount for enabled autoscaling is invalid")
		}
	}
	return nil
}

func checkCreateClusterReqValidity(aksCloudSpec *apiv2.AKSCloudSpec) error {
	if len(aksCloudSpec.ClusterSpec.Location) == 0 {
		return errors.NewBadRequest("required field is missing: Location")
	}
	agentPoolProfiles := aksCloudSpec.ClusterSpec.MachineDeploymentSpec
	if agentPoolProfiles == nil || agentPoolProfiles.BasicSettings.Mode != AgentPoolModeSystem {
		return errors.NewBadRequest("Must define at least one system pool!")
	}
	return checkCreatePoolReqValidity(agentPoolProfiles)
}

func createOrImportAKSCluster(ctx context.Context, name string, userInfoGetter provider.UserInfoGetter, project *kubermaticv1.Project, cloud *apiv2.ExternalClusterCloudSpec, clusterProvider provider.ExternalClusterProvider, privilegedClusterProvider provider.PrivilegedExternalClusterProvider) (*kubermaticv1.ExternalCluster, error) {
	// check whether required fields for cluster import are provided
	fields := reflect.ValueOf(cloud.AKS).Elem()
	for i := 0; i < fields.NumField(); i++ {
		yourjsonTags := fields.Type().Field(i).Tag.Get("required")
		if strings.Contains(yourjsonTags, "true") && fields.Field(i).IsZero() {
			return nil, errors.NewBadRequest("required field is missing: %v", fields.Type().Field(i).Name)
		}
	}

	if cloud.AKS.ClusterSpec != nil {
		if err := checkCreateClusterReqValidity(cloud.AKS); err != nil {
			return nil, err
		}
		if err := createNewAKSCluster(ctx, cloud.AKS); err != nil {
			return nil, err
		}
	}
	newCluster := genExternalCluster(name, project.Name)
	newCluster.Spec.CloudSpec = &kubermaticv1.ExternalClusterCloudSpec{
		AKS: &kubermaticv1.ExternalClusterAKSCloudSpec{
			Name:          cloud.AKS.Name,
			ResourceGroup: cloud.AKS.ResourceGroup,
		},
	}
	keyRef, err := clusterProvider.CreateOrUpdateCredentialSecretForCluster(ctx, cloud, project.Name, newCluster.Name)
	if err != nil {
		return nil, common.KubernetesErrorToHTTPError(err)
	}
	kuberneteshelper.AddFinalizer(newCluster, apiv1.CredentialsSecretsCleanupFinalizer)
	newCluster.Spec.CloudSpec.AKS.CredentialsReference = keyRef

	return createNewCluster(ctx, userInfoGetter, clusterProvider, privilegedClusterProvider, newCluster, project)
}

func patchAKSCluster(ctx context.Context, oldCluster, newCluster *apiv2.ExternalCluster, secretKeySelector provider.SecretKeySelectorValueFunc, cloud *kubermaticv1.ExternalClusterCloudSpec) (*apiv2.ExternalCluster, error) {
	clusterName := cloud.AKS.Name
	resourceGroup := cloud.AKS.ResourceGroup

	newVersion := newCluster.Spec.Version.Semver().String()

	cred, err := aks.GetCredentialsForCluster(*cloud, secretKeySelector)
	if err != nil {
		return nil, err
	}

	aksClient, err := aks.GetAKSClusterClient(cred)
	if err != nil {
		return nil, err
	}
	aksCluster, err := aks.GetAKSCluster(ctx, aksClient, cloud)
	if err != nil {
		return nil, err
	}

	location := aksCluster.Location

	updateCluster := containerservice.ManagedCluster{
		Location: location,
		ManagedClusterProperties: &containerservice.ManagedClusterProperties{
			KubernetesVersion: &newVersion,
		},
	}
	_, err = aksClient.CreateOrUpdate(ctx, resourceGroup, clusterName, updateCluster)
	if err != nil {
		return nil, err
	}

	return newCluster, nil
}

func getAKSNodePools(ctx context.Context, cluster *kubermaticv1.ExternalCluster, secretKeySelector provider.SecretKeySelectorValueFunc, clusterProvider provider.ExternalClusterProvider) ([]apiv2.ExternalClusterMachineDeployment, error) {
	cloud := cluster.Spec.CloudSpec

	cred, err := aks.GetCredentialsForCluster(*cloud, secretKeySelector)
	if err != nil {
		return nil, err
	}

	aksClient, err := aks.GetAKSClusterClient(cred)
	if err != nil {
		return nil, err
	}
	aksCluster, err := aks.GetAKSCluster(ctx, aksClient, cloud)
	if err != nil {
		return nil, err
	}

	poolProfiles := *aksCluster.ManagedClusterProperties.AgentPoolProfiles

	return getAKSMachineDeployments(poolProfiles, cluster, clusterProvider)
}

func getAKSMachineDeployments(poolProfiles []containerservice.ManagedClusterAgentPoolProfile, cluster *kubermaticv1.ExternalCluster, clusterProvider provider.ExternalClusterProvider) ([]apiv2.ExternalClusterMachineDeployment, error) {
	machineDeployments := make([]apiv2.ExternalClusterMachineDeployment, 0, len(poolProfiles))

	nodes, err := clusterProvider.ListNodes(cluster)
	if err != nil {
		return nil, common.KubernetesErrorToHTTPError(err)
	}

	for _, poolProfile := range poolProfiles {
		var readyReplicas int32
		for _, n := range nodes.Items {
			if n.Labels != nil {
				if n.Labels[AKSNodepoolNameLabel] == *poolProfile.Name {
					readyReplicas++
				}
			}
		}
		machineDeployments = append(machineDeployments, createMachineDeploymentFromAKSNodePoll(poolProfile, readyReplicas))
	}

	return machineDeployments, nil
}

func getAKSNodePool(ctx context.Context, cluster *kubermaticv1.ExternalCluster, nodePoolName string, secretKeySelector provider.SecretKeySelectorValueFunc, credentialsReference *providerconfig.GlobalSecretKeySelector, clusterProvider provider.ExternalClusterProvider) (*apiv2.ExternalClusterMachineDeployment, error) {
	cloud := cluster.Spec.CloudSpec

	cred, err := aks.GetCredentialsForCluster(*cloud, secretKeySelector)
	if err != nil {
		return nil, err
	}

	aksClient, err := aks.GetAKSClusterClient(cred)
	if err != nil {
		return nil, err
	}
	aksCluster, err := aks.GetAKSCluster(ctx, aksClient, cloud)
	if err != nil {
		return nil, err
	}

	var poolProfile containerservice.ManagedClusterAgentPoolProfile
	var flag bool = false
	for _, agentPoolProperty := range *aksCluster.ManagedClusterProperties.AgentPoolProfiles {
		if *agentPoolProperty.Name == nodePoolName {
			poolProfile = agentPoolProperty
			flag = true
			break
		}
	}
	if !flag {
		return nil, fmt.Errorf("no nodePool found with the name: %v", nodePoolName)
	}

	return getAKSMachineDeployment(poolProfile, cluster, clusterProvider)
}

func getAKSMachineDeployment(poolProfile containerservice.ManagedClusterAgentPoolProfile, cluster *kubermaticv1.ExternalCluster, clusterProvider provider.ExternalClusterProvider) (*apiv2.ExternalClusterMachineDeployment, error) {
	nodes, err := clusterProvider.ListNodes(cluster)
	if err != nil {
		return nil, common.KubernetesErrorToHTTPError(err)
	}

	var readyReplicas int32
	for _, n := range nodes.Items {
		if n.Labels != nil {
			if n.Labels[AKSNodepoolNameLabel] == *poolProfile.Name {
				readyReplicas++
			}
		}
	}
	md := createMachineDeploymentFromAKSNodePoll(poolProfile, readyReplicas)
	return &md, nil
}

func createMachineDeploymentFromAKSNodePoll(nodePool containerservice.ManagedClusterAgentPoolProfile, readyReplicas int32) apiv2.ExternalClusterMachineDeployment {
	name := to.String(nodePool.Name)
	md := apiv2.ExternalClusterMachineDeployment{
		NodeDeployment: apiv1.NodeDeployment{
			ObjectMeta: apiv1.ObjectMeta{
				ID:   name,
				Name: name,
			},
			Spec: apiv1.NodeDeploymentSpec{
				Replicas: to.Int32(nodePool.Count),
				Template: apiv1.NodeSpec{
					Versions: apiv1.NodeVersionInfo{
						Kubelet: to.String(nodePool.OrchestratorVersion),
					},
				},
			},
			Status: clusterv1alpha1.MachineDeploymentStatus{
				Replicas:      to.Int32(nodePool.Count),
				ReadyReplicas: readyReplicas,
			},
		},
		Cloud: &apiv2.ExternalClusterMachineDeploymentCloudSpec{
			AKS: &apiv2.AKSMachineDeploymentCloudSpec{},
		},
	}

	md.Cloud.AKS.Name = name
	md.Cloud.AKS.BasicSettings = apiv2.AgentPoolBasics{
		Mode:                string(nodePool.Mode),
		AvailabilityZones:   to.StringSlice(nodePool.AvailabilityZones),
		OrchestratorVersion: to.String(nodePool.OrchestratorVersion),
		VMSize:              to.String(nodePool.VMSize),
		EnableAutoScaling:   to.Bool(nodePool.EnableAutoScaling),
		Count:               to.Int32(nodePool.Count),
		OsDiskSizeGB:        to.Int32(nodePool.OsDiskSizeGB),
	}
	if md.Cloud.AKS.BasicSettings.EnableAutoScaling {
		md.Cloud.AKS.BasicSettings.ScalingConfig.MaxCount = to.Int32(nodePool.MaxCount)
		md.Cloud.AKS.BasicSettings.ScalingConfig.MinCount = to.Int32(nodePool.MinCount)
	}
	md.Cloud.AKS.Configuration = apiv2.AgentPoolConfig{
		OsType:             string(nodePool.OsType),
		OsDiskType:         string(nodePool.OsDiskType),
		VnetSubnetID:       to.String(nodePool.VnetSubnetID),
		PodSubnetID:        to.String(nodePool.VnetSubnetID),
		MaxPods:            to.Int32(nodePool.MaxPods),
		EnableNodePublicIP: to.Bool(nodePool.EnableNodePublicIP),
	}
	md.Cloud.AKS.OptionalSettings = apiv2.AgentPoolOptionalSettings{
		NodeLabels: nodePool.NodeLabels,
		NodeTaints: to.StringSlice(nodePool.NodeTaints),
	}
	if nodePool.UpgradeSettings == nil {
		return md
	}
	md.Cloud.AKS.Configuration.MaxSurgeUpgradeSetting = to.String(nodePool.UpgradeSettings.MaxSurge)
	return md
}

func getAKSNodes(cluster *kubermaticv1.ExternalCluster, nodePoolName string, clusterProvider provider.ExternalClusterProvider) ([]apiv2.ExternalClusterNode, error) {
	var nodesV1 []apiv2.ExternalClusterNode

	nodes, err := clusterProvider.ListNodes(cluster)
	if err != nil {
		return nil, common.KubernetesErrorToHTTPError(err)
	}
	for _, n := range nodes.Items {
		if n.Labels != nil {
			if n.Labels[AKSNodepoolNameLabel] == nodePoolName {
				outNode, err := outputNode(n)
				if err != nil {
					return nil, fmt.Errorf("failed to output node %s: %w", n.Name, err)
				}
				nodesV1 = append(nodesV1, *outNode)
			}
		}
	}

	return nodesV1, err
}

func patchAKSMachineDeployment(ctx context.Context, oldCluster, newCluster *apiv2.ExternalClusterMachineDeployment, secretKeySelector provider.SecretKeySelectorValueFunc, cloud *kubermaticv1.ExternalClusterCloudSpec) (*apiv2.ExternalClusterMachineDeployment, error) {
	cred, err := aks.GetCredentialsForCluster(*cloud, secretKeySelector)
	if err != nil {
		return nil, err
	}

	agentPoolClient, err := getAKSNodePoolClient(cred)
	if err != nil {
		return nil, err
	}

	nodePoolName := newCluster.NodeDeployment.Name
	currentReplicas := oldCluster.NodeDeployment.Spec.Replicas
	desiredReplicas := newCluster.NodeDeployment.Spec.Replicas
	currentVersion := oldCluster.NodeDeployment.Spec.Template.Versions.Kubelet
	desiredVersion := newCluster.NodeDeployment.Spec.Template.Versions.Kubelet
	if desiredReplicas != currentReplicas {
		_, err = resizeAKSNodePool(ctx, *agentPoolClient, cloud, nodePoolName, desiredReplicas)
		if err != nil {
			return nil, err
		}
		newCluster.NodeDeployment.Status.Replicas = desiredReplicas
		return newCluster, nil
	}
	if desiredVersion != currentVersion {
		_, err = upgradeNodePool(ctx, *agentPoolClient, cloud, nodePoolName, desiredVersion)
		if err != nil {
			return nil, err
		}
		newCluster.NodeDeployment.Spec.Replicas = currentReplicas
		return newCluster, nil
	}

	return newCluster, nil
}

func resizeAKSNodePool(ctx context.Context, agentPoolClient containerservice.AgentPoolsClient, cloud *kubermaticv1.ExternalClusterCloudSpec, nodePoolName string, desiredSize int32) (*containerservice.AgentPoolsCreateOrUpdateFuture, error) {
	pool, err := agentPoolClient.Get(ctx, cloud.AKS.ResourceGroup, cloud.AKS.Name, nodePoolName)
	if err != nil {
		return nil, err
	}
	nodePool := containerservice.AgentPool{
		Name: &nodePoolName,
		ManagedClusterAgentPoolProfileProperties: &containerservice.ManagedClusterAgentPoolProfileProperties{
			Count: &desiredSize,
		},
	}
	if pool.ManagedClusterAgentPoolProfileProperties.Mode == containerservice.AgentPoolModeSystem {
		nodePool.ManagedClusterAgentPoolProfileProperties.Mode = containerservice.AgentPoolModeSystem
	}
	update, err := updateAKSNodePool(ctx, agentPoolClient, cloud, nodePoolName, nodePool)
	if err != nil {
		return nil, err
	}

	return update, nil
}

func upgradeNodePool(ctx context.Context, agentPoolClient containerservice.AgentPoolsClient, cloud *kubermaticv1.ExternalClusterCloudSpec, nodePoolName string, desiredVersion string) (*containerservice.AgentPoolsCreateOrUpdateFuture, error) {
	pool, err := agentPoolClient.Get(ctx, cloud.AKS.ResourceGroup, cloud.AKS.Name, nodePoolName)
	if err != nil {
		return nil, err
	}
	nodePool := containerservice.AgentPool{
		Name: &nodePoolName,
		ManagedClusterAgentPoolProfileProperties: &containerservice.ManagedClusterAgentPoolProfileProperties{
			OrchestratorVersion: &desiredVersion,
		},
	}
	if pool.ManagedClusterAgentPoolProfileProperties.Mode == containerservice.AgentPoolModeSystem {
		nodePool.ManagedClusterAgentPoolProfileProperties.Mode = containerservice.AgentPoolModeSystem
	}
	update, err := updateAKSNodePool(ctx, agentPoolClient, cloud, nodePoolName, nodePool)
	if err != nil {
		return nil, err
	}

	return update, nil
}

func getAKSNodePoolClient(cred resources.AKSCredentials) (*containerservice.AgentPoolsClient, error) {
	var err error

	agentPoolClient := containerservice.NewAgentPoolsClient(cred.SubscriptionID)
	agentPoolClient.Authorizer, err = auth.NewClientCredentialsConfig(cred.ClientID, cred.ClientSecret, cred.TenantID).Authorizer()
	if err != nil {
		return nil, fmt.Errorf("failed to create authorizer: %w", err)
	}
	return &agentPoolClient, nil
}

func updateAKSNodePool(ctx context.Context, agentPoolClient containerservice.AgentPoolsClient, cloud *kubermaticv1.ExternalClusterCloudSpec, nodePoolName string, nodePool containerservice.AgentPool) (*containerservice.AgentPoolsCreateOrUpdateFuture, error) {
	result, err := agentPoolClient.CreateOrUpdate(ctx, cloud.AKS.ResourceGroup, cloud.AKS.Name, nodePoolName, nodePool)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func deleteAKSNodeGroup(ctx context.Context, cloud *kubermaticv1.ExternalClusterCloudSpec, nodePoolName string, secretKeySelector provider.SecretKeySelectorValueFunc, credentialsReference *providerconfig.GlobalSecretKeySelector, clusterProvider provider.ExternalClusterProvider) error {
	cred, err := aks.GetCredentialsForCluster(*cloud, secretKeySelector)
	if err != nil {
		return err
	}

	agentPoolClient, err := getAKSNodePoolClient(cred)
	if err != nil {
		return err
	}

	_, err = agentPoolClient.Delete(ctx, cloud.AKS.ResourceGroup, cloud.AKS.Name, nodePoolName)
	if err != nil {
		return err
	}
	return nil
}

func createAKSNodePool(ctx context.Context, cloud *kubermaticv1.ExternalClusterCloudSpec, machineDeployment apiv2.ExternalClusterMachineDeployment, secretKeySelector provider.SecretKeySelectorValueFunc, credentialsReference *providerconfig.GlobalSecretKeySelector) (*apiv2.ExternalClusterMachineDeployment, error) {
	cred, err := aks.GetCredentialsForCluster(*cloud, secretKeySelector)
	if err != nil {
		return nil, err
	}

	agentPoolClient, err := getAKSNodePoolClient(cred)
	if err != nil {
		return nil, err
	}

	nodePool := &containerservice.AgentPool{
		Name: &machineDeployment.Name,
	}

	aksMD := machineDeployment.Cloud.AKS
	if err := checkCreatePoolReqValidity(aksMD); err != nil {
		return nil, err
	}
	basicSettings := aksMD.BasicSettings
	optionalSettings := aksMD.OptionalSettings
	property := containerservice.ManagedClusterAgentPoolProfileProperties{
		VMSize:              to.StringPtr(basicSettings.VMSize),
		Count:               to.Int32Ptr(basicSettings.Count),
		OrchestratorVersion: to.StringPtr(basicSettings.OrchestratorVersion),
		Mode:                (containerservice.AgentPoolMode)(basicSettings.Mode),
		OsDiskSizeGB:        to.Int32Ptr(basicSettings.OsDiskSizeGB),
		AvailabilityZones:   &basicSettings.AvailabilityZones,
		EnableAutoScaling:   to.BoolPtr(basicSettings.EnableAutoScaling),
		NodeLabels:          optionalSettings.NodeLabels,
		NodeTaints:          &optionalSettings.NodeTaints,
	}
	if basicSettings.EnableAutoScaling {
		property.MaxCount = to.Int32Ptr(basicSettings.ScalingConfig.MaxCount)
		property.MinCount = to.Int32Ptr(basicSettings.ScalingConfig.MinCount)
	}
	nodePool.ManagedClusterAgentPoolProfileProperties = &property

	_, err = agentPoolClient.CreateOrUpdate(ctx, cloud.AKS.ResourceGroup, cloud.AKS.Name, *nodePool.Name, *nodePool)
	if err != nil {
		return nil, err
	}

	return &machineDeployment, nil
}

func getAKSClusterDetails(ctx context.Context, apiCluster *apiv2.ExternalCluster, secretKeySelector provider.SecretKeySelectorValueFunc, cloud *kubermaticv1.ExternalClusterCloudSpec) (*apiv2.ExternalCluster, error) {
	cred, err := aks.GetCredentialsForCluster(*cloud, secretKeySelector)
	if err != nil {
		return nil, err
	}
	aksClient, err := aks.GetAKSClusterClient(cred)
	if err != nil {
		return nil, err
	}
	aksCluster, err := aks.GetAKSCluster(ctx, aksClient, cloud)
	if err != nil {
		return nil, err
	}
	clusterSpec := apiv2.AKSClusterSpec{}
	aksClusterProperties := aksCluster.ManagedClusterProperties
	if aksClusterProperties == nil {
		apiCluster.Cloud.AKS.ClusterSpec = &clusterSpec
		return apiCluster, nil
	}
	clusterSpec = apiv2.AKSClusterSpec{
		Location:          to.String(aksCluster.Location),
		DNSPrefix:         to.String(aksCluster.DNSPrefix),
		KubernetesVersion: to.String(aksClusterProperties.KubernetesVersion),
		EnableRBAC:        to.Bool(aksClusterProperties.EnableRBAC),
		NodeResourceGroup: to.String(aksClusterProperties.NodeResourceGroup),
		FqdnSubdomain:     to.String(aksClusterProperties.FqdnSubdomain),
		Fqdn:              to.String(aksClusterProperties.Fqdn),
		PrivateFQDN:       to.String(aksClusterProperties.PrivateFQDN),
	}
	networkProfile := aksClusterProperties.NetworkProfile
	if networkProfile == nil {
		apiCluster.Cloud.AKS.ClusterSpec = &clusterSpec
		return apiCluster, nil
	}
	clusterSpec.NetworkProfile = apiv2.AKSNetworkProfile{
		NetworkPlugin:    string(networkProfile.NetworkPlugin),
		NetworkPolicy:    string(networkProfile.NetworkPolicy),
		NetworkMode:      string(networkProfile.NetworkMode),
		PodCidr:          to.String(networkProfile.PodCidr),
		ServiceCidr:      to.String(networkProfile.ServiceCidr),
		DNSServiceIP:     to.String(networkProfile.DNSServiceIP),
		DockerBridgeCidr: to.String(networkProfile.DockerBridgeCidr),
	}
	if aksCluster.AadProfile == nil {
		apiCluster.Cloud.AKS.ClusterSpec = &clusterSpec
		return apiCluster, nil
	}
	clusterSpec.ManagedAAD = to.Bool(aksCluster.AadProfile.Managed)
	apiCluster.Cloud.AKS.ClusterSpec = &clusterSpec

	return apiCluster, nil
}
