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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

// IngressRuleValueApplyConfiguration represents an declarative configuration of the IngressRuleValue type for use
// with apply.
type IngressRuleValueApplyConfiguration struct {
	HTTP *HTTPIngressRuleValueApplyConfiguration `json:"http,omitempty"`
	TCP  *TCPIngressRuleValueApplyConfiguration  `json:"tcp,omitempty"`
}

// IngressRuleValueApplyConfiguration constructs an declarative configuration of the IngressRuleValue type for use with
// apply.
func IngressRuleValue() *IngressRuleValueApplyConfiguration {
	return &IngressRuleValueApplyConfiguration{}
}

// WithHTTP sets the HTTP field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the HTTP field is set to the value of the last call.
func (b *IngressRuleValueApplyConfiguration) WithHTTP(value *HTTPIngressRuleValueApplyConfiguration) *IngressRuleValueApplyConfiguration {
	b.HTTP = value
	return b
}

// WithTCP sets the TCP field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the TCP field is set to the value of the last call.
func (b *IngressRuleValueApplyConfiguration) WithTCP(value *TCPIngressRuleValueApplyConfiguration) *IngressRuleValueApplyConfiguration {
	b.TCP = value
	return b
}
