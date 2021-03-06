//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright AppsCode Inc. and Contributors

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

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1

import (
	common "k8s.io/kube-openapi/pkg/common"
	spec "k8s.io/kube-openapi/pkg/validation/spec"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"kmodules.xyz/monitoring-agent-api/api/v1.AgentSpec":              schema_kmodulesxyz_monitoring_agent_api_api_v1_AgentSpec(ref),
		"kmodules.xyz/monitoring-agent-api/api/v1.PrometheusExporterSpec": schema_kmodulesxyz_monitoring_agent_api_api_v1_PrometheusExporterSpec(ref),
		"kmodules.xyz/monitoring-agent-api/api/v1.PrometheusSpec":         schema_kmodulesxyz_monitoring_agent_api_api_v1_PrometheusSpec(ref),
		"kmodules.xyz/monitoring-agent-api/api/v1.ServiceMonitorSpec":     schema_kmodulesxyz_monitoring_agent_api_api_v1_ServiceMonitorSpec(ref),
	}
}

func schema_kmodulesxyz_monitoring_agent_api_api_v1_AgentSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"agent": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"prometheus": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("kmodules.xyz/monitoring-agent-api/api/v1.PrometheusSpec"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"kmodules.xyz/monitoring-agent-api/api/v1.PrometheusSpec"},
	}
}

func schema_kmodulesxyz_monitoring_agent_api_api_v1_PrometheusExporterSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"port": {
						SchemaProps: spec.SchemaProps{
							Description: "Port number for the exporter side car.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"args": {
						SchemaProps: spec.SchemaProps{
							Description: "Arguments to the entrypoint. The docker image's CMD is used if this is not provided. Variable references $(VAR_NAME) are expanded using the container's environment. If a variable cannot be resolved, the reference in the input string will be unchanged. The $(VAR_NAME) syntax can be escaped with a double $$, ie: $$(VAR_NAME). Escaped references will never be expanded, regardless of whether the variable exists or not. Cannot be updated. More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: "",
										Type:    []string{"string"},
										Format:  "",
									},
								},
							},
						},
					},
					"env": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-patch-merge-key": "name",
								"x-kubernetes-patch-strategy":  "merge",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "List of environment variables to set in the container. Cannot be updated.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("k8s.io/api/core/v1.EnvVar"),
									},
								},
							},
						},
					},
					"resources": {
						SchemaProps: spec.SchemaProps{
							Description: "Compute Resources required by exporter container. Cannot be updated. More info: https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/",
							Default:     map[string]interface{}{},
							Ref:         ref("k8s.io/api/core/v1.ResourceRequirements"),
						},
					},
					"securityContext": {
						SchemaProps: spec.SchemaProps{
							Description: "Security options the pod should run with. More info: https://kubernetes.io/docs/concepts/policy/security-context/ More info: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/",
							Ref:         ref("k8s.io/api/core/v1.SecurityContext"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"k8s.io/api/core/v1.EnvVar", "k8s.io/api/core/v1.ResourceRequirements", "k8s.io/api/core/v1.SecurityContext"},
	}
}

func schema_kmodulesxyz_monitoring_agent_api_api_v1_PrometheusSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"exporter": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("kmodules.xyz/monitoring-agent-api/api/v1.PrometheusExporterSpec"),
						},
					},
					"serviceMonitor": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("kmodules.xyz/monitoring-agent-api/api/v1.ServiceMonitorSpec"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"kmodules.xyz/monitoring-agent-api/api/v1.PrometheusExporterSpec", "kmodules.xyz/monitoring-agent-api/api/v1.ServiceMonitorSpec"},
	}
}

func schema_kmodulesxyz_monitoring_agent_api_api_v1_ServiceMonitorSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"labels": {
						SchemaProps: spec.SchemaProps{
							Description: "Labels are key value pairs that is used to select Prometheus instance via ServiceMonitor labels.",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: "",
										Type:    []string{"string"},
										Format:  "",
									},
								},
							},
						},
					},
					"interval": {
						SchemaProps: spec.SchemaProps{
							Description: "Interval at which metrics should be scraped",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
	}
}
