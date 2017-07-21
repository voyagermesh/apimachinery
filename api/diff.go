package api

import (
	"errors"
	"net"
	"reflect"
	"sort"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ingressClassAnnotationKey   = "kubernetes.io/ingress.class"
	ingressClassAnnotationValue = "voyager"
)

// if ingressClass == "voyager", then only handle ingress that has voyager annotation
// if ingressClass == "", then handle no annotaion or voyager annotation
func (r Ingress) ShouldHandleIngress(ingressClass string) bool {
	// https://github.com/appscode/voyager/blob/master/api/conversion_v1beta1.go#L44
	if r.Annotations[APISchema] == APISchemaEngress {
		// Resource Type is Extended Ingress So we should always Handle this
		return true
	}
	kubeAnnotation, _ := r.Annotations[ingressClassAnnotationKey]
	return kubeAnnotation == ingressClass || kubeAnnotation == ingressClassAnnotationValue
}

func (r Ingress) HasChanged(o Ingress) (bool, error) {
	if r.Name != o.Name ||
		r.Namespace != o.Namespace ||
		r.APISchema() != o.APISchema() {
		return false, errors.New("Not the same Ingress.")
	}

	if !reflect.DeepEqual(r.Spec, o.Spec) {
		return true, nil
	}

	ra := map[string]string{}
	for k, v := range r.Annotations {
		if strings.HasPrefix(k, EngressKey+"/") {
			ra[k] = v
		}
	}
	oa := map[string]string{}
	for k, v := range o.Annotations {
		if strings.HasPrefix(k, EngressKey+"/") {
			oa[k] = v
		}
	}
	return !reflect.DeepEqual(ra, oa), nil
}

func (r Ingress) Ports() []int {
	usesTLS := func(h string) bool {
		for _, tls := range r.Spec.TLS {
			for _, host := range tls.Hosts {
				if host == h {
					return true
				}
			}
		}
		return false
	}

	usesHTTPRule := false
	ports := map[int]string{}
	for _, rule := range r.Spec.Rules {
		if rule.HTTP != nil {
			usesHTTPRule = true
			if usesTLS(rule.Host) {
				ports[443] = "https"
			} else {
				ports[80] = "http"
			}
		}
		for _, port := range rule.TCP {
			p := port.Port.IntValue()
			if p > 0 {
				ports[p] = "tcp"
			}
		}
	}
	if !usesHTTPRule && r.Spec.Backend != nil {
		ports[80] = "http"
	}

	result := make([]int, len(ports))
	i := 0
	for k := range ports {
		result[i] = k
		i++
	}
	return result
}

func (r Ingress) IsPortChanged(o Ingress) bool {
	rPorts := r.Ports()
	oPorts := o.Ports()

	if len(rPorts) != len(oPorts) {
		return true
	}

	sort.Ints(rPorts)
	sort.Ints(oPorts)
	for i := range rPorts {
		if rPorts[i] != oPorts[i] {
			return true
		}
	}
	return false
}

func (r Ingress) IsSecretChanged(o Ingress) bool {
	oldSecretLists := map[string]bool{}
	for _, rs := range r.Spec.TLS {
		oldSecretLists[rs.SecretName] = true
	}
	for _, rs := range r.Spec.Rules {
		for _, tcp := range rs.TCP {
			oldSecretLists[tcp.SecretName] = true
		}
	}

	for _, rs := range o.Spec.Rules {
		for _, port := range rs.TCP {
			if _, ok := oldSecretLists[port.SecretName]; !ok {
				return true
			}
		}
	}
	for _, rs := range o.Spec.TLS {
		if _, ok := oldSecretLists[rs.SecretName]; !ok {
			return true
		}
	}
	return false
}

func (r Ingress) IsLoadBalancerSourceRangeChanged(o Ingress) bool {
	oldipset := make(map[string]bool)
	for _, oldrange := range r.Spec.LoadBalancerSourceRanges {
		k, ok := ipnet(oldrange)
		if ok {
			oldipset[k] = true
		}
	}
	newipset := make(map[string]bool)
	for _, newrange := range o.Spec.LoadBalancerSourceRanges {
		k, ok := ipnet(newrange)
		if ok {
			newipset[k] = true
			if _, found := oldipset[k]; !found {
				return true
			}
		}
	}
	return len(newipset) != len(oldipset)
}

func ipnet(spec string) (string, bool) {
	spec = strings.TrimSpace(spec)
	_, ipnet, err := net.ParseCIDR(spec)
	if err != nil {
		return "", false
	}
	return ipnet.String(), true
}

func (r Ingress) IsStatsChanged(o Ingress) bool {
	return isMapKeyChanged(r.Annotations, o.Annotations, StatsOn) ||
		isMapKeyChanged(r.Annotations, o.Annotations, StatsPort) ||
		isMapKeyChanged(r.Annotations, o.Annotations, StatsServiceName) ||
		isMapKeyChanged(r.Annotations, o.Annotations, StatsSecret)
}

func (r Ingress) IsStatsSecretChanged(o Ingress) bool {
	return isMapKeyChanged(r.Annotations, o.Annotations, StatsSecret)
}

func (r Ingress) IsKeepSourceChanged(o Ingress, cloudProvider string) bool {
	return cloudProvider == "aws" &&
		o.LBType() == LBTypeLoadBalancer &&
		isMapKeyChanged(r.Annotations, o.Annotations, KeepSourceIP)
}

func isMapKeyChanged(oldMap map[string]string, newMap map[string]string, key string) bool {
	oldValue, oldOk := oldMap[key]
	newValue, newOk := newMap[key]
	return oldOk != newOk || oldValue != newValue
}

func (r Ingress) BackendServices() map[string]metav1.ObjectMeta {
	services := map[string]metav1.ObjectMeta{}

	record := func(svcName string) {
		parts := strings.SplitN(svcName, ".", 2)
		if len(parts) == 2 {
			services[svcName] = metav1.ObjectMeta{
				Name:      parts[0],
				Namespace: parts[1],
			}
		} else {
			services[svcName+"."+r.Namespace] = metav1.ObjectMeta{
				Name:      parts[0],
				Namespace: r.Namespace,
			}
		}
	}

	if r.Spec.Backend != nil {
		record(r.Spec.Backend.ServiceName)
	}
	for _, rules := range r.Spec.Rules {
		if rules.HTTP != nil {
			for _, svc := range rules.HTTP.Paths {
				record(svc.Backend.ServiceName)
			}
		}
		for _, svc := range rules.TCP {
			record(svc.Backend.ServiceName)
		}
	}

	return services
}

func (r Ingress) HasBackendService(name, namespace string) bool {
	svcFQN := name + "." + namespace

	fqn := func(svcName string) string {
		if strings.ContainsRune(svcName, '.') {
			return svcName
		}
		return svcName + "." + r.Namespace
	}

	if r.Spec.Backend != nil {
		if fqn(r.Spec.Backend.ServiceName) == svcFQN {
			return true
		}
	}
	for _, rules := range r.Spec.Rules {
		if rules.HTTP != nil {
			for _, svc := range rules.HTTP.Paths {
				if fqn(svc.Backend.ServiceName) == svcFQN {
					return true
				}
			}
		}
		for _, svc := range rules.TCP {
			if fqn(svc.Backend.ServiceName) == svcFQN {
				return true
			}
		}
	}
	return false
}
