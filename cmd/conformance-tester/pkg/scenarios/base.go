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

package scenarios

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"go.uber.org/zap"

	clusterv1alpha1 "github.com/kubermatic/machine-controller/pkg/apis/cluster/v1alpha1"
	"github.com/kubermatic/machine-controller/pkg/cloudprovider/util"
	providerconfig "github.com/kubermatic/machine-controller/pkg/providerconfig/types"
	"k8c.io/kubermatic/v2/cmd/conformance-tester/pkg/types"
	apiv1 "k8c.io/kubermatic/v2/pkg/api/v1"
	kubermaticv1 "k8c.io/kubermatic/v2/pkg/apis/kubermatic/v1"
	"k8c.io/kubermatic/v2/pkg/semver"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
)

type Scenario interface {
	// these are all satisfied by the baseScenario

	CloudProvider() providerconfig.CloudProvider
	OperatingSystem() providerconfig.OperatingSystem
	ContainerRuntime() string
	DualstackEnabled() bool
	Version() semver.Semver
	Datacenter() *kubermaticv1.Datacenter
	Name() string
	Log(log *zap.SugaredLogger) *zap.SugaredLogger
	NamedLog(log *zap.SugaredLogger) *zap.SugaredLogger
	IsValid(opts *types.Options, log *zap.SugaredLogger) bool

	// these are implemented per provider
	OperatingSystemSpec() (*apiv1.OperatingSystemSpec, error)

	Cluster(secrets types.Secrets) *kubermaticv1.ClusterSpec
	MachineDeployments(ctx context.Context, num int, secrets types.Secrets, cluster *kubermaticv1.Cluster) ([]clusterv1alpha1.MachineDeployment, error)

	SetDualstackEnabled(bool)
}

type baseScenario struct {
	cloudProvider    providerconfig.CloudProvider
	operatingSystem  providerconfig.OperatingSystem
	version          semver.Semver
	containerRuntime string
	dualstackEnabled bool
	datacenter       *kubermaticv1.Datacenter
}

func (s *baseScenario) CloudProvider() providerconfig.CloudProvider {
	return s.cloudProvider
}

func (s *baseScenario) OperatingSystem() providerconfig.OperatingSystem {
	return s.operatingSystem
}

func (s *baseScenario) Version() semver.Semver {
	return s.version
}

func (s *baseScenario) ContainerRuntime() string {
	return s.containerRuntime
}

func (s *baseScenario) DualstackEnabled() bool {
	return s.dualstackEnabled
}

func (s *baseScenario) Datacenter() *kubermaticv1.Datacenter {
	return s.datacenter
}

func (s *baseScenario) Log(log *zap.SugaredLogger) *zap.SugaredLogger {
	return log.With(
		"provider", s.cloudProvider,
		"os", s.operatingSystem,
		"version", s.version.String(),
		"cri", s.containerRuntime,
	)
}

func (s *baseScenario) NamedLog(log *zap.SugaredLogger) *zap.SugaredLogger {
	return log.With("scenario", s.Name())
}

func (s *baseScenario) Name() string {
	return fmt.Sprintf("%s-%s-%s-%s", s.cloudProvider, s.operatingSystem, s.containerRuntime, s.version.String())
}

func (s *baseScenario) IsValid(opts *types.Options, log *zap.SugaredLogger) bool {
	return true
}

func (s *baseScenario) createMachineDeployment(replicas int, providerSpec interface{}) (clusterv1alpha1.MachineDeployment, error) {
	osSpec, err := s.OperatingSystemSpec()
	if err != nil {
		return clusterv1alpha1.MachineDeployment{}, fmt.Errorf("failed to build OS spec: %w", err)
	}

	return createMachineDeployment(replicas, &s.version, s.operatingSystem, osSpec, s.cloudProvider, providerSpec, s.dualstackEnabled)
}

func (s *baseScenario) OperatingSystemSpec() (*apiv1.OperatingSystemSpec, error) {
	switch s.operatingSystem {
	case providerconfig.OperatingSystemUbuntu:
		return &apiv1.OperatingSystemSpec{
			Ubuntu: &apiv1.UbuntuSpec{},
		}, nil

	case providerconfig.OperatingSystemCentOS:
		return &apiv1.OperatingSystemSpec{
			CentOS: &apiv1.CentOSSpec{},
		}, nil

	case providerconfig.OperatingSystemSLES:
		return &apiv1.OperatingSystemSpec{
			SLES: &apiv1.SLESSpec{},
		}, nil

	case providerconfig.OperatingSystemRHEL:
		return &apiv1.OperatingSystemSpec{
			RHEL: &apiv1.RHELSpec{},
		}, nil

	case providerconfig.OperatingSystemFlatcar:
		return &apiv1.OperatingSystemSpec{
			Flatcar: &apiv1.FlatcarSpec{},
		}, nil

	case providerconfig.OperatingSystemRockyLinux:
		return &apiv1.OperatingSystemSpec{
			RockyLinux: &apiv1.RockyLinuxSpec{},
		}, nil

	case providerconfig.OperatingSystemAmazonLinux2:
		return &apiv1.OperatingSystemSpec{
			AmazonLinux: &apiv1.AmazonLinuxSpec{},
		}, nil

	default:
		return nil, errors.New("cannot create OS spec based on the scenario: unknown OS")
	}
}

func (s *baseScenario) SetDualstackEnabled(enabled bool) {
	s.dualstackEnabled = enabled
}

func createMachineDeployment(replicas int, version *semver.Semver, os providerconfig.OperatingSystem, osSpec interface{}, provider providerconfig.CloudProvider, providerSpec interface{}, dualstack bool) (clusterv1alpha1.MachineDeployment, error) {
	replicas32 := int32(replicas)

	var encodedOSSpec json.RawMessage
	var err error

	switch osSpec := osSpec.(type) {
	case *apiv1.OperatingSystemSpec:
		switch os {
		case providerconfig.OperatingSystemRHEL:
			encodedOSSpec, err = json.Marshal(osSpec.RHEL)
		case providerconfig.OperatingSystemSLES:
			encodedOSSpec, err = json.Marshal(osSpec.SLES)
		case providerconfig.OperatingSystemUbuntu:
			encodedOSSpec, err = json.Marshal(osSpec.Ubuntu)
		case providerconfig.OperatingSystemFlatcar:
			encodedOSSpec, err = json.Marshal(osSpec.Flatcar)
		case providerconfig.OperatingSystemCentOS:
			encodedOSSpec, err = json.Marshal(osSpec.CentOS)
		case providerconfig.OperatingSystemRockyLinux:
			encodedOSSpec, err = json.Marshal(osSpec.RockyLinux)
		case providerconfig.OperatingSystemAmazonLinux2:
			encodedOSSpec, err = json.Marshal(osSpec.AmazonLinux)
		}
	default:
		encodedOSSpec, err = json.Marshal(osSpec)
	}

	if err != nil {
		return clusterv1alpha1.MachineDeployment{}, err
	}

	encodedCloudProviderSpec, err := json.Marshal(providerSpec)
	if err != nil {
		return clusterv1alpha1.MachineDeployment{}, err
	}

	cfg := providerconfig.Config{
		CloudProvider: provider,
		CloudProviderSpec: runtime.RawExtension{
			Raw: encodedCloudProviderSpec,
		},
		OperatingSystem: os,
		OperatingSystemSpec: runtime.RawExtension{
			Raw: encodedOSSpec,
		},
	}

	if dualstack {
		cfg.Network = &providerconfig.NetworkConfig{
			IPFamily: util.DualStack,
		}
	}

	encodedConfig, err := json.Marshal(cfg)
	if err != nil {
		return clusterv1alpha1.MachineDeployment{}, err
	}

	machineLabels := map[string]string{
		"machine": "md-" + utilrand.String(5),
	}

	return clusterv1alpha1.MachineDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "e2e-" + utilrand.String(5),
			Namespace: metav1.NamespaceSystem,
		},
		Spec: clusterv1alpha1.MachineDeploymentSpec{
			Replicas: &replicas32,
			Selector: metav1.LabelSelector{
				MatchLabels: machineLabels,
			},
			Template: clusterv1alpha1.MachineTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: machineLabels,
				},
				Spec: clusterv1alpha1.MachineSpec{
					Versions: clusterv1alpha1.MachineVersionInfo{
						Kubelet: version.String(),
					},
					ProviderSpec: clusterv1alpha1.ProviderSpec{
						Value: &runtime.RawExtension{
							Raw: encodedConfig,
						},
					},
				},
			},
		},
	}, nil
}
