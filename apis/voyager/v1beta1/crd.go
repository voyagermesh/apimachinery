/*
Copyright The Voyager Authors.

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
package v1beta1

import (
	"github.com/go-openapi/spec"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/kube-openapi/pkg/common"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
)

var (
	EnableStatusSubresource bool
)

func (r Ingress) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourceIngresses,
		Singular:      ResourceIngress,
		Kind:          ResourceKindIngress,
		ShortNames:    []string{"ing"},
		Categories:    []string{"networking", "appscode", "all"},
		ResourceScope: string(apiextensions.NamespaceScoped),
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    SchemeGroupVersion.Version,
				Served:  true,
				Storage: true,
			},
		},
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "voyager"},
		},
		SpecDefinitionName:      "github.com/appscode/voyager/apis/voyager/v1beta1.Ingress",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: EnableStatusSubresource,
		AdditionalPrinterColumns: []apiextensions.CustomResourceColumnDefinition{
			{
				Name:     "HOSTS",
				Type:     "string",
				JSONPath: ".spec.rules[0].host",
			},
			// Address? Port ?
			{
				Name:     "LOAD_BALANCER_IP",
				Type:     "string",
				JSONPath: ".status.loadBalancer.ingress",
			},
			{
				Name:     "AGE",
				Type:     "date",
				JSONPath: ".metadata.creationTimestamp",
			},
		},
	}, setNameSchema)
}

func (c Certificate) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourceCertificates,
		Singular:      ResourceCertificate,
		Kind:          ResourceKindCertificate,
		ShortNames:    []string{"cert"},
		Categories:    []string{"networking", "appscode", "all"},
		ResourceScope: string(apiextensions.NamespaceScoped),
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    SchemeGroupVersion.Version,
				Served:  true,
				Storage: true,
			},
		},
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "voyager"},
		},
		SpecDefinitionName:      "github.com/appscode/voyager/apis/voyager/v1beta1.Certificate",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: EnableStatusSubresource,
		AdditionalPrinterColumns: []apiextensions.CustomResourceColumnDefinition{
			{
				Name:     "Domains",
				Type:     "string",
				JSONPath: ".spec.domains[0]",
			},
			{
				Name:     "Age",
				Type:     "date",
				JSONPath: ".metadata.creationTimestamp",
			},
		},
	})
}

func setNameSchema(openapiSpec map[string]common.OpenAPIDefinition) {
	// ref: https://github.com/kubedb/project/issues/166
	// https://github.com/kubernetes/apimachinery/blob/94ebb086c69b9fec4ddbfb6a1433d28ecca9292b/pkg/util/validation/validation.go#L153
	var maxLength int64 = 63
	openapiSpec["k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"].Schema.SchemaProps.Properties["name"] = spec.Schema{
		SchemaProps: spec.SchemaProps{
			Description: "Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names",
			Type:        []string{"string"},
			Format:      "",
			Pattern:     `^[a-z]([-a-z0-9]*[a-z0-9])?$`,
			MaxLength:   &maxLength,
		},
	}
}
