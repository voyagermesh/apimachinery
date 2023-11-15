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

// AuthOptionApplyConfiguration represents an declarative configuration of the AuthOption type for use
// with apply.
type AuthOptionApplyConfiguration struct {
	Basic *BasicAuthApplyConfiguration `json:"basic,omitempty"`
	TLS   *TLSAuthApplyConfiguration   `json:"tls,omitempty"`
	OAuth []OAuthApplyConfiguration    `json:"oauth,omitempty"`
}

// AuthOptionApplyConfiguration constructs an declarative configuration of the AuthOption type for use with
// apply.
func AuthOption() *AuthOptionApplyConfiguration {
	return &AuthOptionApplyConfiguration{}
}

// WithBasic sets the Basic field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Basic field is set to the value of the last call.
func (b *AuthOptionApplyConfiguration) WithBasic(value *BasicAuthApplyConfiguration) *AuthOptionApplyConfiguration {
	b.Basic = value
	return b
}

// WithTLS sets the TLS field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the TLS field is set to the value of the last call.
func (b *AuthOptionApplyConfiguration) WithTLS(value *TLSAuthApplyConfiguration) *AuthOptionApplyConfiguration {
	b.TLS = value
	return b
}

// WithOAuth adds the given value to the OAuth field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the OAuth field.
func (b *AuthOptionApplyConfiguration) WithOAuth(values ...*OAuthApplyConfiguration) *AuthOptionApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithOAuth")
		}
		b.OAuth = append(b.OAuth, *values[i])
	}
	return b
}
