package api

import (
	schema "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/runtime"
	versionedwatch "k8s.io/kubernetes/pkg/watch/versioned"
)

// SchemeGroupVersion is group version used to register these objects
var V1beta1SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1beta1"}

var (
	V1beta1SchemeBuilder = runtime.NewSchemeBuilder(v1addKnownTypes, addConversionFuncs)
	V1betaAddToScheme    = V1beta1SchemeBuilder.AddToScheme
)

// Adds the list of known types to api.Scheme.
func v1addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(V1beta1SchemeGroupVersion,
		&Ingress{},
		&IngressList{},

		&Alert{},
		&AlertList{},

		&Certificate{},
		&CertificateList{},

		&Backup{},
		&BackupList{},

		&v1.ListOptions{},
	)
	versionedwatch.AddToGroupVersion(scheme, V1beta1SchemeGroupVersion)
	return nil
}
