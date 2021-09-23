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

package v1

import (
	"testing"

	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIsValid(t *testing.T) {
	for k, expected := range dataTables {
		err := k.IsValid("kind")
		actual := err == nil
		if expected != actual {
			t.Errorf("Failed Tests: %s, Expected: %v, Actual: %v, Reason %v", k.Name, expected, actual, err)
		}
	}
}

var dataTables = map[*Ingress]bool{
	{
		ObjectMeta: metav1.ObjectMeta{Name: "Empty backend service name"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
			},
		},
	}: false,
	{
		ObjectMeta: metav1.ObjectMeta{Name: "Named backend service port"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "service-with-named-port",
									Port: networking.ServiceBackendPort{Name: "http"},
								},
							},
						},
					},
				},
			},
		},
	}: true,
	{
		ObjectMeta: metav1.ObjectMeta{Name: "Invalid backend service name"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: ".default",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
			},
		},
	}: false,
	{
		ObjectMeta: metav1.ObjectMeta{Name: "spec.rule[0] can specify either HTTP or TCP"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{},
						TCP:  &TCPIngressRuleValue{},
					},
				},
			},
		},
	}: false,
	{
		ObjectMeta: metav1.ObjectMeta{Name: "No DefaultBackend Service For TCP"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
						},
					},
				},
			},
		},
	}: false,
	{
		ObjectMeta: metav1.ObjectMeta{Name: "No Listen Port for TCP"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
			},
		},
	}: false,
	{
		ObjectMeta: metav1.ObjectMeta{Name: "TCP and HTTP in Same Port specified"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 80,
							Paths: []HTTPIngressPath{
								{
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 3444},
										},
									},
								},
							},
						},
					},
				},
				{
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 80,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
			},
		},
	}: false,
	{
		ObjectMeta: metav1.ObjectMeta{Name: "TCP and HTTP in Same Port not specified"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Paths: []HTTPIngressPath{
								{
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 3444},
										},
									},
								},
							},
						},
					},
				},
				{
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 80,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
			},
		},
	}: false,
	{
		ObjectMeta: metav1.ObjectMeta{Name: "HTTP with host and path"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					Host: "this-is-host-one",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 80,
							Paths: []HTTPIngressPath{
								{
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 3444},
										},
									},
								},
							},
						},
					},
				},
				{
					Host: "this-is-host-one",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 80,
							Paths: []HTTPIngressPath{
								{
									Path: "/should-be-true",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 3444},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}: true,
	{
		ObjectMeta: metav1.ObjectMeta{Name: "HTTP with hosts"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					Host: "this-is-host-one",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 80,
							Paths: []HTTPIngressPath{
								{
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 3444},
										},
									},
								},
							},
						},
					},
				},
				{
					Host: "this-is-host-one",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 8091,
							Paths: []HTTPIngressPath{
								{
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 3444},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}: true,
	{
		ObjectMeta: metav1.ObjectMeta{Name: "TCP multi-host in same port"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					Host: "voyager.appscode.test",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
				{
					Host: "voyager.appscode.com",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo2",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
			},
		},
	}: true, // this should work after adding tcp-sni support
	{
		ObjectMeta: metav1.ObjectMeta{Name: "TCP multi-host with TLS"},
		Spec: IngressSpec{
			TLS: []IngressTLS{
				{
					SecretName: "voyager-cert",
					Hosts: []string{
						"voyager.appscode.test",
						"voyager.appscode.com",
					},
				},
			},
			Rules: []IngressRule{
				{
					Host: "voyager.appscode.test",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
				{
					Host: "voyager.appscode.com",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo2",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
			},
		},
	}: true, // useTLS for TCP multi-host
	{
		ObjectMeta: metav1.ObjectMeta{Name: "TCP multi-host with conflicting TLS"},
		Spec: IngressSpec{
			TLS: []IngressTLS{
				{
					SecretName: "voyager-cert",
					Hosts: []string{
						"voyager.appscode.test",
					},
				},
			},
			Rules: []IngressRule{
				{
					Host: "voyager.appscode.test",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
				{
					Host: "voyager.appscode.com",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo2",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
			},
		},
	}: false, // conflicting useTLS
	{
		ObjectMeta: metav1.ObjectMeta{Name: "TCP multi-host ALPN conflict"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					Host: "voyager.appscode.test",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							ALPN: []string{"a", "b"},
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
				{
					Host: "voyager.appscode.com",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							ALPN: []string{"x", "y"},
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo2",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
			},
		},
	}: false, // conflicting ALPN
	{
		ObjectMeta: metav1.ObjectMeta{Name: "TCP multi-host ALPN conflict"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					Host: "voyager.appscode.test",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
				{
					Host: "voyager.appscode.com",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							ALPN: []string{"x", "y"},
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo2",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
			},
		},
	}: false, // conflicting ALPN with NoALPN
	{
		ObjectMeta: metav1.ObjectMeta{Name: "TCP multi-host with NoTLS"},
		Spec: IngressSpec{
			TLS: []IngressTLS{
				{
					SecretName: "voyager-cert",
					Hosts: []string{
						"voyager.appscode.test",
					},
				},
			},
			Rules: []IngressRule{
				{
					Host: "voyager.appscode.test",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							NoTLS: true,
							Port:  3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
				{
					Host: "voyager.appscode.com",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo2",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
			},
		},
	}: true, // useTLS will be false for both rules
	{
		ObjectMeta: metav1.ObjectMeta{Name: "TCP multi-host with reusing host"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					Host: "voyager.appscode.test",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
				{
					Host: "voyager.appscode.test",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo2",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
			},
		},
	}: false, // some host under same address-binder
	{
		ObjectMeta: metav1.ObjectMeta{Name: "TCP multi-host with wildcard host"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					Host: "*",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
				{
					Host: "voyager.appscode.test",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo2",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
			},
		},
	}: false, // wildcard host with multiple rules under same address-binder
	{
		ObjectMeta: metav1.ObjectMeta{Name: "TCP multi-host with wildcard host (2)"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					Host: "voyager.appscode.test",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
				{
					Host: "*",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo2",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
			},
		},
	}: false, // wildcard host with multiple rules under same address-binder
	{
		ObjectMeta: metav1.ObjectMeta{Name: "TCP multi-host with empty host"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
				{
					Host: "voyager.appscode.test",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo2",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
			},
		},
	}: false, // empty host with multiple rules under same address-binder
	{
		ObjectMeta: metav1.ObjectMeta{Name: "TCP with different port"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
				{
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3435,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo2",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
			},
		},
	}: true,
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "TLS in passthrough mode",
			Annotations: map[string]string{
				SSLPassthrough: "true",
			},
		},
		Spec: IngressSpec{
			TLS: []IngressTLS{
				{
					SecretName: "voyager-cert",
					Hosts: []string{
						"voyager.appscode.test",
					},
				},
			},
			Rules: []IngressRule{
				{
					Host: "voyager.appscode.test",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							Port: 3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
			},
		},
	}: false, // TLS defined in passthrough mode
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "NoTLS in passthrough mode",
			Annotations: map[string]string{
				SSLPassthrough: "true",
			},
		},
		Spec: IngressSpec{
			TLS: []IngressTLS{
				{
					SecretName: "voyager-cert",
					Hosts: []string{
						"voyager.appscode.test",
					},
				},
			},
			Rules: []IngressRule{
				{
					Host: "voyager.appscode.test",
					IngressRuleValue: IngressRuleValue{
						TCP: &TCPIngressRuleValue{
							NoTLS: true,
							Port:  3434,
							Backend: IngressBackend{
								Service: networking.IngressServiceBackend{
									Name: "foo",
									Port: networking.ServiceBackendPort{Number: 3444},
								},
							},
						},
					},
				},
			},
		},
	}: true, // NoTLS in passthrough mode
	{
		ObjectMeta: metav1.ObjectMeta{Name: "Multi rule"},
		Spec: IngressSpec{
			DefaultBackend: &IngressBackend{
				Service: networking.IngressServiceBackend{
					Name: "foo",
					Port: networking.ServiceBackendPort{Number: 80},
				}},
			Rules: []IngressRule{
				{
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Paths: []HTTPIngressPath{
								{
									Path: "/test-dns",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 80},
										}},
								},
								{
									Path: "/test-no-dns",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 80},
										}},
								},
								{
									Path: "/test-no-backend-redirect",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 80},
										}},
								},
								{
									Path: "/test-no-backend-rule-redirect",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 80},
										},
										BackendRules: []string{
											"http-request redirect location https://google.com code 302",
										},
									},
								},
							},
						},
					},
				},
				{
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Paths: []HTTPIngressPath{
								{
									Path: "/redirect-rule",
									Backend: IngressBackend{
										BackendRules: []string{
											"http-request redirect location https://github.com/appscode/discuss/issues code 301",
										},
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 80},
										},
									},
								},
							},
						},
					},
				},
				{
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Paths: []HTTPIngressPath{
								{
									Path: "/redirect",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 80},
										},
									},
								},
							},
						},
					},
				},
				{
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Paths: []HTTPIngressPath{
								{
									Path: "/back-end",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 8989},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}: true,
	{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "https://github.com/voyagermesh/voyager/issues/420",
			Namespace: "default",
			Annotations: map[string]string{
				"ingress.appscode.com/type": "HostPort",
			},
		},
		Spec: IngressSpec{
			TLS: []IngressTLS{
				{
					SecretName: "voyager-cert",
					Hosts: []string{
						"minicluster.example.com",
					},
				},
			},
			Rules: []IngressRule{
				{
					Host: "minicluster.example.com",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Paths: []HTTPIngressPath{
								{
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "cluster-nginx",
											Port: networking.ServiceBackendPort{Number: 80},
										},
									},
								},
							},
						},
					},
				},
				{
					Host: "minicluster.example.com",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							NoTLS: true,
							Paths: []HTTPIngressPath{
								{
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "cluster-nginx",
											Port: networking.ServiceBackendPort{Number: 80},
										},
									},
								},
							},
						},
					},
				},
				{
					Host: "domain1.com",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Paths: []HTTPIngressPath{
								{
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "cluster-nginx",
											Port: networking.ServiceBackendPort{Number: 80},
										},
									},
								},
							},
						},
					},
				},
				{
					Host: "domain2.com",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Paths: []HTTPIngressPath{
								{
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "cluster-nginx",
											Port: networking.ServiceBackendPort{Number: 80},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}: true,
	{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "https://github.com/voyagermesh/voyager/pull/768",
			Namespace: "default",
			Annotations: map[string]string{
				"ingress.appscode.com/type": "Internal",
			},
		},
		Spec: IngressSpec{
			TLS: []IngressTLS{
				{
					SecretName: "voyager-cert",
					Hosts: []string{
						"*.example.org.dmmy.me",
					},
				},
			},
			Rules: []IngressRule{
				{
					Host: "*.example.org.dmmy.me",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Paths: []HTTPIngressPath{
								{
									Path: "/resources",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "distro-static",
											Port: networking.ServiceBackendPort{Number: 80},
										},
									},
								},
								{
									Path: "/admin_resources",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "admin-resources",
											Port: networking.ServiceBackendPort{Number: 80},
										},
									},
								},
							},
						},
					},
				},
				{
					Host: "*.example.org.dmmy.me",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							NoTLS: true,
							Paths: []HTTPIngressPath{
								{
									Path: "/resources",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "distro-static",
											Port: networking.ServiceBackendPort{Number: 80},
										},
									},
								},
								{
									Path: "/admin_resources",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "admin-resources",
											Port: networking.ServiceBackendPort{Number: 80},
										},
									},
								},
							},
						},
					},
				},
				{
					Host: "example.org.dmmy.me",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Paths: []HTTPIngressPath{
								{
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "verify-tls-cert-dmmy-me",
											Port: networking.ServiceBackendPort{Number: 80},
										},
									},
								},
							},
						},
					},
				},
			},
			DefaultBackend: &IngressBackend{
				Service: networking.IngressServiceBackend{
					Name: "distro-biz",
					Port: networking.ServiceBackendPort{Number: 80},
				},
			},
		},
	}: true,
	{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "acme-http-challenge-path",
			Namespace: "default",
			Annotations: map[string]string{
				"ingress.appscode.com/type": LBTypeLoadBalancer,
			},
		},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							NoTLS: true,
							Paths: []HTTPIngressPath{
								{
									Path: "/.well-known/acme-challenge/",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "voyager-operator.kube-system",
											Port: networking.ServiceBackendPort{Number: 56791},
										},
									},
								},
							},
						},
					},
				},
				{
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							NoTLS:    true,
							NodePort: 32666,
							Paths: []HTTPIngressPath{
								{
									Path: "/",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "web",
											Port: networking.ServiceBackendPort{Number: 80},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}: true,
	{
		ObjectMeta: metav1.ObjectMeta{Name: "Merging empty-host with wildcard-host"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 80,
							Paths: []HTTPIngressPath{
								{
									Path: "/foo",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 3444},
										},
									},
								},
							},
						},
					},
				},
				{
					Host: "*",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 80,
							Paths: []HTTPIngressPath{
								{
									Path: "/not-foo",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 3444},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}: true, // merging "*" host with empty-host
	{
		ObjectMeta: metav1.ObjectMeta{Name: "Conflict merging empty-host with wildcard-host"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 80,
							Paths: []HTTPIngressPath{
								{
									Path: "/foo",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 3444},
										},
									},
								},
							},
						},
					},
				},
				{
					Host: "*",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 80,
							Paths: []HTTPIngressPath{
								{
									Path: "/foo",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 3444},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}: false, // conflicting paths while merging "*" host with empty-host
	{
		ObjectMeta: metav1.ObjectMeta{Name: "Merging empty-host with wildcard-host"},
		Spec: IngressSpec{
			TLS: []IngressTLS{
				{
					SecretName: "voyager-cert",
					Hosts: []string{
						"*",
					},
				},
			},
			Rules: []IngressRule{
				{
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 80,
							Paths: []HTTPIngressPath{
								{
									Path: "/foo",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 3444},
										},
									},
								},
							},
						},
					},
				},
				{
					Host: "*",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 80,
							Paths: []HTTPIngressPath{
								{
									Path: "/not-foo",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 3444},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}: false, // conflicting TLS merging "*" host with empty-host
	{
		ObjectMeta: metav1.ObjectMeta{Name: "HTTP ALPN conflict"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					Host: "voyager.appscode.test",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 3434,
							ALPN: []string{"a", "b", "c"},
							Paths: []HTTPIngressPath{
								{
									Path: "/path-1",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 3444},
										},
									},
								},
							},
						},
					},
				},
				{
					Host: "voyager.appscode.com",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 3434,
							ALPN: []string{"x", "y"},
							Paths: []HTTPIngressPath{
								{
									Path: "/path-2",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 3444},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}: false, // conflicting ALPN
	{
		ObjectMeta: metav1.ObjectMeta{Name: "HTTP conflicting ALPN with NoALPN"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					Host: "voyager.appscode.test",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 3434,
							ALPN: []string{"a", "b", "c"},
							Paths: []HTTPIngressPath{
								{
									Path: "/path-1",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 3444},
										},
									},
								},
							},
						},
					},
				},
				{
					Host: "voyager.appscode.com",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 3434,
							Paths: []HTTPIngressPath{
								{
									Path: "/path-2",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 3444},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}: false, // conflicting ALPN with NoALPN
	{
		ObjectMeta: metav1.ObjectMeta{Name: "HTTP conflicting Proto"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					Host: "voyager.appscode.test",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port:  3434,
							Proto: "http/1.1",
							Paths: []HTTPIngressPath{
								{
									Path: "/path-1",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 3444},
										},
									},
								},
							},
						},
					},
				},
				{
					Host: "voyager.appscode.com",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port:  3434,
							Proto: "h2",
							Paths: []HTTPIngressPath{
								{
									Path: "/path-2",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 3444},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}: false, // conflicting Proto
	{
		ObjectMeta: metav1.ObjectMeta{Name: "Default backend Proto with ALPN"},
		Spec: IngressSpec{
			DefaultBackend: &IngressBackend{
				Service: networking.IngressServiceBackend{
					Name: "foo",
					Port: networking.ServiceBackendPort{Number: 3444},
				},
				ALPN:  []string{"a", "b", "c"},
				Proto: "h2",
			},
		},
	}: false, // Proto with ALPN
	{
		ObjectMeta: metav1.ObjectMeta{Name: "HTTP Proto with ALPN"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					Host: "voyager.appscode.test",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port:  3434,
							ALPN:  []string{"a", "b", "c"},
							Proto: "h2",
						},
					},
				},
			},
		},
	}: false, // Proto with ALPN
	{
		ObjectMeta: metav1.ObjectMeta{Name: "HTTP backend Proto with ALPN"},
		Spec: IngressSpec{
			Rules: []IngressRule{
				{
					Host: "voyager.appscode.test",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Paths: []HTTPIngressPath{
								{
									Path: "/path-1",
									Backend: IngressBackend{
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 3444},
										},
										ALPN:  []string{"a", "b", "c"},
										Proto: "h2",
									},
								},
							},
						},
					},
				},
			},
		},
	}: false, // Proto with ALPN
	{
		ObjectMeta: metav1.ObjectMeta{Name: "Multiple Oauth under same port"},
		Spec: IngressSpec{
			FrontendRules: []FrontendRule{
				{
					Port: 80,
					Auth: &AuthOption{
						OAuth: []OAuth{
							{
								Host:        "team01.appscode.test",
								AuthBackend: "auth-be-1",
							},
							{
								Host:        "team02.appscode.test",
								AuthBackend: "auth-be-2",
							},
						},
					},
				},
			},
			Rules: []IngressRule{
				{
					Host: "team01.appscode.test",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 80,
							Paths: []HTTPIngressPath{
								{
									Path: "/path-1",
									Backend: IngressBackend{
										Name: "auth-be-1",
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 8080},
										},
									},
								},
							},
						},
					},
				},
				{
					Host: "team02.appscode.test",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 80,
							Paths: []HTTPIngressPath{
								{
									Path: "/path-2",
									Backend: IngressBackend{
										Name: "auth-be-2",
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 8080},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}: true,
	{
		ObjectMeta: metav1.ObjectMeta{Name: "Multiple Oauth under same port"},
		Spec: IngressSpec{
			FrontendRules: []FrontendRule{
				{
					Port: 80,
					Auth: &AuthOption{
						OAuth: []OAuth{
							{
								Host:        "auth.appscode.test",
								AuthBackend: "auth-be-1",
							},
							{
								Host:        "auth.appscode.test",
								AuthBackend: "auth-be-2",
							},
						},
					},
				},
			},
			Rules: []IngressRule{
				{
					Host: "auth.appscode.test",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 80,
							Paths: []HTTPIngressPath{
								{
									Path: "/path-1",
									Backend: IngressBackend{
										Name: "auth-be-1",
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 8080},
										},
									},
								},
								{
									Path: "/path-2",
									Backend: IngressBackend{
										Name: "auth-be-2",
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 8080},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}: false, // multiple oauth for same host under same port
	{
		ObjectMeta: metav1.ObjectMeta{Name: "Auth backend not found"},
		Spec: IngressSpec{
			FrontendRules: []FrontendRule{
				{
					Port: 80,
					Auth: &AuthOption{
						OAuth: []OAuth{
							{
								Host:        "auth.appscode.test",
								AuthBackend: "auth-be",
							},
						},
					},
				},
			},
			Rules: []IngressRule{
				{
					Host: "auth.appscode.test",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 8080,
							Paths: []HTTPIngressPath{
								{
									Path: "/path-1",
									Backend: IngressBackend{
										Name: "auth-be-1",
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 8080},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}: false, // auth backend not found: port not matched
	{
		ObjectMeta: metav1.ObjectMeta{Name: "Auth backend not found"},
		Spec: IngressSpec{
			FrontendRules: []FrontendRule{
				{
					Port: 80,
					Auth: &AuthOption{
						OAuth: []OAuth{
							{
								Host:        "auth.appscode.test",
								AuthBackend: "auth-be",
							},
						},
					},
				},
			},
			Rules: []IngressRule{
				{
					Host: "voyager.appscode.test",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 80,
							Paths: []HTTPIngressPath{
								{
									Path: "/path-1",
									Backend: IngressBackend{
										Name: "auth-be-1",
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 8080},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}: false, // auth backend not found: host not matched
	{
		ObjectMeta: metav1.ObjectMeta{Name: "Auth backend not found"},
		Spec: IngressSpec{
			FrontendRules: []FrontendRule{
				{
					Port: 80,
					Auth: &AuthOption{
						OAuth: []OAuth{
							{
								Host:        "auth.appscode.test",
								AuthBackend: "auth-be",
							},
						},
					},
				},
			},
			Rules: []IngressRule{
				{
					Host: "auth.appscode.test",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 80,
							Paths: []HTTPIngressPath{
								{
									Path: "/path-1",
									Backend: IngressBackend{
										Name: "auth-be-2",
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 8080},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}: false, // auth backend not found: backend name not matched
	{
		ObjectMeta: metav1.ObjectMeta{Name: "FrontendRule without auth"},
		Spec: IngressSpec{
			FrontendRules: []FrontendRule{
				{
					Port: 80,
					Rules: []string{
						"acl acl_fake path_beg /fake",
					},
				},
			},
			Rules: []IngressRule{
				{
					Host: "auth.appscode.test",
					IngressRuleValue: IngressRuleValue{
						HTTP: &HTTPIngressRuleValue{
							Port: 80,
							Paths: []HTTPIngressPath{
								{
									Path: "/path-1",
									Backend: IngressBackend{
										Name: "auth-be-2",
										Service: networking.IngressServiceBackend{
											Name: "foo",
											Port: networking.ServiceBackendPort{Number: 8080},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}: true, // frontend-rule without auth
}
