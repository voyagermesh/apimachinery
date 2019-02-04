package v1beta1

import (
	core "k8s.io/api/core/v1"
)

// LocalTypedReference contains enough information to let you inspect or modify the referred object.
type LocalTypedReference struct {
	// Kind of the referent.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	Kind string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`
	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,3,opt,name=name"`
	// API version of the referent.
	// +optional
	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,5,opt,name=apiVersion"`
}

// Represents the source of a volume to mount.
// Only one of its members may be specified.
type VolumeSource struct {
	Name string `json:"name,omitempty"`
	// Secret represents a secret that should populate this volume.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#secret
	// +optional
	Secret *core.SecretVolumeSource `json:"secret,omitempty"`
	// ConfigMap represents a configMap that should populate this volume
	// +optional
	ConfigMap *core.ConfigMapVolumeSource `json:"configMap,omitempty"`
	// Items for all in one resources secrets, configmaps, and downward API
	// +optional
	Projected *core.ProjectedVolumeSource `json:"projected,omitempty"`
	// Path within the container at which the volume should be mounted.  Must
	// not contain ':'.
	MountPath string `json:"mountPath"`
}
