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

package cilium

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	appskubermaticv1 "k8c.io/kubermatic/v2/pkg/apis/apps.kubermatic/v1"
	kubermaticv1 "k8c.io/kubermatic/v2/pkg/apis/kubermatic/v1"
	"k8c.io/kubermatic/v2/pkg/resources"
	"k8c.io/kubermatic/v2/pkg/resources/reconciling"
	"k8c.io/kubermatic/v2/pkg/resources/registry"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const (
	ciliumHelmChartName = "cilium"

	ciliumImageRegistry = "quay.io/cilium/"
)

// ApplicationDefinitionReconciler creates Cilium ApplicationDefinition managed by KKP to be used
// for installing Cilium CNI into KKP usr clusters.
func ApplicationDefinitionReconciler(config *kubermaticv1.KubermaticConfiguration) reconciling.NamedApplicationDefinitionReconcilerFactory {
	return func() (string, reconciling.ApplicationDefinitionReconciler) {
		return kubermaticv1.CNIPluginTypeCilium.String(), func(app *appskubermaticv1.ApplicationDefinition) (*appskubermaticv1.ApplicationDefinition, error) {
			app.Labels = map[string]string{
				appskubermaticv1.ApplicationManagedByLabel: appskubermaticv1.ApplicationManagedByKKPValue,
				appskubermaticv1.ApplicationTypeLabel:      appskubermaticv1.ApplicationTypeCNIValue,
			}

			app.Spec.Description = "Cilium CNI - eBPF-based Networking, Security, and Observability"
			app.Spec.Method = appskubermaticv1.HelmTemplateMethod

			var credentials *appskubermaticv1.HelmCredentials
			if config.Spec.UserCluster.SystemApplications.HelmRegistryConfigFile != nil {
				credentials = &appskubermaticv1.HelmCredentials{
					RegistryConfigFile: config.Spec.UserCluster.SystemApplications.HelmRegistryConfigFile,
				}
			}

			app.Spec.Versions = []appskubermaticv1.ApplicationVersion{
				// NOTE: When introducing a new version, make sure it is:
				//  - introduced in pkg/cni/version.go with the version string exactly matching the Spec.Versions.Version here
				//  - Helm chart is mirrored in Kubermatic OCI registry, use the script cilium-mirror-chart.sh
				{
					Version: "1.13.0",
					Template: appskubermaticv1.ApplicationTemplate{
						Source: appskubermaticv1.ApplicationSource{
							Helm: &appskubermaticv1.HelmSource{
								ChartName:    ciliumHelmChartName,
								ChartVersion: "1.13.0",
								URL:          "oci://" + config.Spec.UserCluster.SystemApplications.HelmRepository,
								Credentials:  credentials,
							},
						},
					},
				},
			}

			// we want to allow overriding the default values, so reconcile them only if nil
			if app.Spec.DefaultValues == nil {
				defaultValues := map[string]any{
					"operator": map[string]any{
						"replicas": 1,
					},
					"hubble": map[string]any{
						"relay": map[string]any{
							"enabled": true,
						},
						"ui": map[string]any{
							"enabled": true,
						},
						// cronJob TLS cert gen method needs to be used for backward compatibility with older KKP
						"tls": map[string]any{
							"auto": map[string]any{
								"method": "cronJob",
							},
						},
					},
				}
				rawValues, err := json.Marshal(defaultValues)
				if err != nil {
					return app, fmt.Errorf("failed to marshall default CNI values: %w", err)
				}
				app.Spec.DefaultValues = &runtime.RawExtension{Raw: rawValues}
			}
			return app, nil
		}
	}
}

// GetAppInstallOverrideValues returns Helm values to be enforced on the cluster's ApplicationInstallation
// of the Cilium CNI managed by KKP.
func GetAppInstallOverrideValues(cluster *kubermaticv1.Cluster, overwriteRegistry string) map[string]any {
	values := map[string]any{
		"cni": map[string]any{
			// we run Cilium as non-exclusive CNI to allow for Multus use-cases
			"exclusive": false,
		},
	}

	valuesOperator := map[string]any{
		"securityContext": map[string]any{
			"seccompProfile": map[string]any{
				"type": "RuntimeDefault",
			},
		},
	}
	values["operator"] = valuesOperator

	if cluster.Spec.ClusterNetwork.ProxyMode == resources.EBPFProxyMode {
		values["kubeProxyReplacement"] = "strict"
		values["k8sServiceHost"] = cluster.Status.Address.ExternalName
		values["k8sServicePort"] = cluster.Status.Address.Port

		nodePortRange := cluster.Spec.ComponentsOverride.Apiserver.NodePortRange
		if nodePortRange != "" && nodePortRange != resources.DefaultNodePortRange {
			values["nodePort"] = map[string]any{
				"range": strings.ReplaceAll(nodePortRange, "-", ","),
			}
		}
	} else {
		values["kubeProxyReplacement"] = "disabled"
	}

	ipamOperator := map[string]any{
		"clusterPoolIPv4PodCIDR":  cluster.Spec.ClusterNetwork.Pods.GetIPv4CIDR(),
		"clusterPoolIPv4MaskSize": fmt.Sprintf("%d", *cluster.Spec.ClusterNetwork.NodeCIDRMaskSizeIPv4),
	}
	if cluster.IsDualStack() {
		values["ipv6"] = map[string]any{"enabled": true}
		ipamOperator["clusterPoolIPv6PodCIDR"] = cluster.Spec.ClusterNetwork.Pods.GetIPv6CIDR()
		ipamOperator["clusterPoolIPv6MaskSize"] = fmt.Sprintf("%d", *cluster.Spec.ClusterNetwork.NodeCIDRMaskSizeIPv6)
	}
	values["ipam"] = map[string]any{"operator": ipamOperator}

	// Override images if registry override is configured
	if overwriteRegistry != "" {
		values["image"] = map[string]any{
			"repository": registry.Must(registry.RewriteImage(ciliumImageRegistry+"cilium", overwriteRegistry)),
			"useDigest":  false,
		}
		valuesOperator["image"] = map[string]any{
			"repository": registry.Must(registry.RewriteImage(ciliumImageRegistry+"operator", overwriteRegistry)),
			"useDigest":  false,
		}
		values["hubble"] = map[string]any{
			"relay": map[string]any{
				"image": map[string]any{
					"repository": registry.Must(registry.RewriteImage(ciliumImageRegistry+"hubble-relay", overwriteRegistry)),
					"useDigest":  false,
				},
			},
			"ui": map[string]any{
				"backend": map[string]any{
					"image": map[string]any{
						"repository": registry.Must(registry.RewriteImage(ciliumImageRegistry+"hubble-ui-backend", overwriteRegistry)),
						"useDigest":  false,
					},
				},
				"frontend": map[string]any{
					"image": map[string]any{
						"repository": registry.Must(registry.RewriteImage(ciliumImageRegistry+"hubble-ui", overwriteRegistry)),
						"useDigest":  false,
					},
				},
			},
		}
		values["certgen"] = map[string]any{
			"image": map[string]any{
				"repository": registry.Must(registry.RewriteImage(ciliumImageRegistry+"certgen", overwriteRegistry)),
				"useDigest":  false,
			},
		}
	}

	return values
}

// ValidateValuesUpdate validates the update operation on provided Cilium Helm values.
func ValidateValuesUpdate(newValues, oldValues map[string]any, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// Validate immutability of specific top-level value subtrees, managed solely by KKP
	allErrs = append(allErrs, validateImmutableValues(newValues, oldValues, fieldPath, []string{
		"cni",
		"ipam",
		"ipv6",
	})...)

	// Validate that mandatory top-level values are present
	allErrs = append(allErrs, validateMandatoryValues(newValues, fieldPath, []string{
		"kubeProxyReplacement", // can be changed if switching KKP proxy mode, but must be present
	})...)

	// Validate config for strict kubeProxyReplacement (ebpf "proxy mode")
	if newValues["kubeProxyReplacement"] == "strict" {
		// Validate mandatory values for kube-proxy-replacement
		allErrs = append(allErrs, validateMandatoryValues(newValues, fieldPath, []string{
			"k8sServiceHost",
			"k8sServicePort",
		})...)
	}

	// Validate immutability of operator.securityContext
	operPath := fieldPath.Child("operator")
	newOperator, ok := newValues["operator"].(map[string]any)
	if !ok {
		allErrs = append(allErrs, field.Invalid(operPath, newValues["operator"], "value missing or incorrect"))
	}
	oldOperator, ok := oldValues["operator"].(map[string]any)
	if !ok {
		allErrs = append(allErrs, field.Invalid(operPath, oldValues["operator"], "value missing or incorrect"))
	}
	allErrs = append(allErrs, validateImmutableValues(newOperator, oldOperator, operPath, []string{
		"securityContext",
	})...)

	return allErrs
}

func validateImmutableValues(newValues, oldValues map[string]any, fieldPath *field.Path, immutableValues []string) field.ErrorList {
	allErrs := field.ErrorList{}
	for _, v := range immutableValues {
		if !reflect.DeepEqual(newValues[v], oldValues[v]) {
			allErrs = append(allErrs, field.Invalid(fieldPath.Child(v), newValues[v], "value is immutable"))
		}
	}
	return allErrs
}

func validateMandatoryValues(newValues map[string]any, fieldPath *field.Path, mandatoryValues []string) field.ErrorList {
	allErrs := field.ErrorList{}
	for _, v := range mandatoryValues {
		if _, ok := newValues[v]; !ok {
			allErrs = append(allErrs, field.NotFound(fieldPath.Child(v), newValues[v]))
		}
	}
	return allErrs
}
