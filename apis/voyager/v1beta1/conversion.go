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

package v1beta1

import (
	"unsafe"

	v1 "voyagermesh.dev/apimachinery/apis/voyager/v1"

	corev1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/util/intstr"
	kbconv "sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this CronJob to the Hub version (v1).
func (src *Ingress) ConvertTo(dstRaw kbconv.Hub) error {
	return Convert_v1beta1_Ingress_To_v1_Ingress(src, dstRaw.(*v1.Ingress), nil)
}

// ConvertFrom converts from the Hub version (v1) to this version.
func (dst *Ingress) ConvertFrom(srcRaw kbconv.Hub) error {
	return Convert_v1_Ingress_To_v1beta1_Ingress(srcRaw.(*v1.Ingress), dst, nil)
}

func Convert_intstr_IntOrString_To_int32(in *intstr.IntOrString, out *int32, s conversion.Scope) error {
	*out = int32(in.IntValue())
	return nil
}

func Convert_int32_To_intstr_IntOrString(in *int32, out *intstr.IntOrString, s conversion.Scope) error {
	*out = intstr.FromInt(int(*in))
	return nil
}

func Convert_v1beta1_IngressBackend_To_v1_IngressBackend(in *IngressBackend, out *v1.IngressBackend, s conversion.Scope) error {
	out.Name = in.Name
	out.HostNames = *(*[]string)(unsafe.Pointer(&in.HostNames))
	out.Service = networking.IngressServiceBackend{
		Name: in.ServiceName,
	}
	if portNo := in.ServicePort.IntValue(); portNo > 0 {
		out.Service.Port.Number = int32(portNo)
	} else {
		out.Service.Port.Name = in.ServicePort.StrVal
	}
	out.BackendRules = *(*[]string)(unsafe.Pointer(&in.BackendRules))
	out.ALPN = *(*[]string)(unsafe.Pointer(&in.ALPN))
	out.Proto = in.Proto
	out.LoadBalanceOn = in.LoadBalanceOn
	return nil
}

func Convert_v1_IngressBackend_To_v1beta1_IngressBackend(in *v1.IngressBackend, out *IngressBackend, s conversion.Scope) error {
	out.Name = in.Name
	out.HostNames = *(*[]string)(unsafe.Pointer(&in.HostNames))
	out.ServiceName = in.Service.Name
	if in.Service.Port.Name != "" {
		out.ServicePort = intstr.FromString(in.Service.Port.Name)
	} else {
		out.ServicePort = intstr.FromInt(int(in.Service.Port.Number))
	}
	out.BackendRules = *(*[]string)(unsafe.Pointer(&in.BackendRules))
	out.ALPN = *(*[]string)(unsafe.Pointer(&in.ALPN))
	out.Proto = in.Proto
	out.LoadBalanceOn = in.LoadBalanceOn
	return nil
}

func Convert_v1beta1_IngressSpec_To_v1_IngressSpec(in *IngressSpec, out *v1.IngressSpec, s conversion.Scope) error {
	if err := Convert_v1beta1_HTTPIngressBackend_To_v1_HTTPIngressBackend(in.Backend, out.DefaultBackend, s); err != nil {
		return err
	}
	out.TLS = *(*[]v1.IngressTLS)(unsafe.Pointer(&in.TLS))
	out.ConfigVolumes = *(*[]v1.VolumeSource)(unsafe.Pointer(&in.ConfigVolumes))
	if in.FrontendRules != nil {
		in, out := &in.FrontendRules, &out.FrontendRules
		*out = make([]v1.FrontendRule, len(*in))
		for i := range *in {
			if err := Convert_v1beta1_FrontendRule_To_v1_FrontendRule(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.FrontendRules = nil
	}
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]v1.IngressRule, len(*in))
		for i := range *in {
			if err := Convert_v1beta1_IngressRule_To_v1_IngressRule(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Rules = nil
	}
	out.LoadBalancerSourceRanges = *(*[]string)(unsafe.Pointer(&in.LoadBalancerSourceRanges))
	out.Resources = in.Resources
	out.NodeSelector = *(*map[string]string)(unsafe.Pointer(&in.NodeSelector))
	out.Affinity = (*corev1.Affinity)(unsafe.Pointer(in.Affinity))
	out.SchedulerName = in.SchedulerName
	out.Tolerations = *(*[]corev1.Toleration)(unsafe.Pointer(&in.Tolerations))
	out.ImagePullSecrets = *(*[]corev1.LocalObjectReference)(unsafe.Pointer(&in.ImagePullSecrets))
	out.PriorityClassName = in.PriorityClassName
	out.Priority = (*int32)(unsafe.Pointer(in.Priority))
	out.SecurityContext = (*corev1.PodSecurityContext)(unsafe.Pointer(in.SecurityContext))
	out.ExternalIPs = *(*[]string)(unsafe.Pointer(&in.ExternalIPs))
	out.LivenessProbe = (*corev1.Probe)(unsafe.Pointer(in.LivenessProbe))
	out.ReadinessProbe = (*corev1.Probe)(unsafe.Pointer(in.ReadinessProbe))
	out.TerminationGracePeriodSeconds = (*int64)(unsafe.Pointer(in.TerminationGracePeriodSeconds))
	return nil
}

func Convert_v1_IngressSpec_To_v1beta1_IngressSpec(in *v1.IngressSpec, out *IngressSpec, s conversion.Scope) error {
	if err := Convert_v1_HTTPIngressBackend_To_v1beta1_HTTPIngressBackend(in.DefaultBackend, out.Backend, s); err != nil {
		return err
	}
	out.TLS = *(*[]IngressTLS)(unsafe.Pointer(&in.TLS))
	out.ConfigVolumes = *(*[]VolumeSource)(unsafe.Pointer(&in.ConfigVolumes))
	if in.FrontendRules != nil {
		in, out := &in.FrontendRules, &out.FrontendRules
		*out = make([]FrontendRule, len(*in))
		for i := range *in {
			if err := Convert_v1_FrontendRule_To_v1beta1_FrontendRule(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.FrontendRules = nil
	}
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]IngressRule, len(*in))
		for i := range *in {
			if err := Convert_v1_IngressRule_To_v1beta1_IngressRule(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Rules = nil
	}
	out.LoadBalancerSourceRanges = *(*[]string)(unsafe.Pointer(&in.LoadBalancerSourceRanges))
	out.Resources = in.Resources
	out.NodeSelector = *(*map[string]string)(unsafe.Pointer(&in.NodeSelector))
	out.Affinity = (*corev1.Affinity)(unsafe.Pointer(in.Affinity))
	out.SchedulerName = in.SchedulerName
	out.Tolerations = *(*[]corev1.Toleration)(unsafe.Pointer(&in.Tolerations))
	out.ImagePullSecrets = *(*[]corev1.LocalObjectReference)(unsafe.Pointer(&in.ImagePullSecrets))
	out.PriorityClassName = in.PriorityClassName
	out.Priority = (*int32)(unsafe.Pointer(in.Priority))
	out.SecurityContext = (*corev1.PodSecurityContext)(unsafe.Pointer(in.SecurityContext))
	out.ExternalIPs = *(*[]string)(unsafe.Pointer(&in.ExternalIPs))
	out.LivenessProbe = (*corev1.Probe)(unsafe.Pointer(in.LivenessProbe))
	out.ReadinessProbe = (*corev1.Probe)(unsafe.Pointer(in.ReadinessProbe))
	out.TerminationGracePeriodSeconds = (*int64)(unsafe.Pointer(in.TerminationGracePeriodSeconds))
	return nil
}
