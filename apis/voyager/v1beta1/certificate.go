package v1beta1

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindCertificate     = "Certificate"
	ResourceSingularCertificate = "certificate"
	ResourcePluralCertificate   = "certificates"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Certificate is a collection of domains for which a SSL certificate is
// issued from Let's Encrypt.
type Certificate struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CertificateSpec   `json:"spec,omitempty"`
	Status            CertificateStatus `json:"status,omitempty"`
}

type CertificateSpec struct {
	// Tries to obtain a single certificate using all domains passed into Domains.
	// The first domain in domains is used for the CommonName field of the certificate, all other
	// domains are added using the Subject Alternate Names extension.
	Domains []string `json:"domains,omitempty"`

	// ChallengeProvider details to verify domains
	ChallengeProvider ChallengeProvider `json:"challengeProvider"`

	// Secret contains ACMEUser information. Secret must contain a key `email`
	// If empty tries to find an Secret via domains
	// if not found create an ACMEUser and stores as a secret.
	// Secrets key to be expected:
	//  ACME_EMAIL -> required, if not provided it will through error.
	//  ACME_SERVER_URL -> custom server url to generate certificates, default is lets encrypt.
	//  ACME_USER_DATA -> user data, if not found one will be created for the provided email,
	//    and stored in the key.
	ACMEUserSecretName string `json:"acmeUserSecretName"`

	// Storage backend to store the certificates currently, kubernetes secret and vault.
	Storage CertificateStorage `json:"storage,omitempty"`

	// Indicates that the certificate is paused.
	// +optional
	Paused bool `json:"paused,omitempty"`
}

type ChallengeProvider struct {
	HTTP *HTTPChallengeProvider `json:"http,omitempty"`
	DNS  *DNSChallengeProvider  `json:"dns,omitempty"`
}

type HTTPChallengeProvider struct {
	Ingress LocalTypedReference `json:"ingress,omitempty"`
}

type DNSChallengeProvider struct {
	// DNS Provider from the list https://github.com/appscode/voyager/blob/master/docs/tasks/certificate/providers.md
	Provider             string `json:"provider,omitempty"`
	CredentialSecretName string `json:"credentialSecretName,omitempty"`
}

type CertificateStorage struct {
	Secret *core.LocalObjectReference `json:"secret,omitempty"`
	Vault  *VaultStore                `json:"vault,omitempty"`
}

type VaultStore struct {
	Name   string `json:"name,omitempty"`
	Prefix string `json:"prefix,omitempty"`
}

type CertificateStatus struct {
	CreationTime          *metav1.Time           `json:"creationTime,omitempty"`
	Conditions            []CertificateCondition `json:"conditions,omitempty"`
	LastIssuedCertificate *CertificateDetails    `json:"lastIssuedCertificate,omitempty"`
}

type ACMECertificateDetails struct {
	Domain        string `json:"domain"`
	CertURL       string `json:"certUrl"`
	CertStableURL string `json:"certStableUrl"`
	AccountRef    string `json:"accountRef,omitempty"`
}

type CertificateDetails struct {
	SerialNumber  string      `json:"serialNumber,omitempty"`
	NotBefore     metav1.Time `json:"notBefore,omitempty"`
	NotAfter      metav1.Time `json:"notAfter,omitempty"`
	CertURL       string      `json:"certURL"`
	CertStableURL string      `json:"certStableURL"`
	AccountRef    string      `json:"accountRef,omitempty"`
}

type RequestConditionType string

// These are the possible conditions for a certificate create request.
const (
	CertificateIssued      RequestConditionType = "Issued"
	CertificateFailed      RequestConditionType = "Failed"
	CertificateRateLimited RequestConditionType = "RateLimited"
)

type CertificateCondition struct {
	// request approval state, currently Approved or Denied.
	Type RequestConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=RequestConditionType"`
	// brief reason for the request state
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,2,opt,name=reason"`
	// human readable message with details about the request state
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,3,opt,name=message"`
	// timestamp for the last update to this condition
	// +optional
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty" protobuf:"bytes,4,opt,name=lastUpdateTime"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CertificateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Certificate `json:"items,omitempty"`
}
