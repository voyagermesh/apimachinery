package v1beta1

import (
	"crypto/x509"
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/appscode/voyager/apis/voyager"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	core_util "kmodules.xyz/client-go/core/v1"
)

const (
	ingressClassAnnotationKey   = "kubernetes.io/ingress.class"
	ingressClassAnnotationValue = "voyager"
)

// if ingressClass == "voyager", then only handle ingress that has voyager annotation
// if ingressClass == "", then handle no annotation or voyager annotation
func (r Ingress) ShouldHandleIngress(ingressClass string) bool {
	// https://github.com/appscode/voyager/blob/master/api/conversion_v1beta1.go#L44
	if r.APISchema() == APISchemaEngress {
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
		return false, errors.New("not the same Ingress")
	}

	if o.DeletionTimestamp != nil && core_util.HasFinalizer(o.ObjectMeta, voyager.GroupName) {
		return true, nil
	}

	specEqual := cmp.Equal(r.Spec, o.Spec, cmp.Comparer(func(x, y resource.Quantity) bool {
		return x.Cmp(y) == 0
	}))
	if !specEqual {
		return true, nil
	}

	ra := map[string]string{}
	for k, v := range r.Annotations {
		if strings.HasPrefix(k, EngressKey+"/") || k == ingressClassAnnotationKey {
			ra[k] = v
		}
	}
	oa := map[string]string{}
	for k, v := range o.Annotations {
		if strings.HasPrefix(k, EngressKey+"/") || k == ingressClassAnnotationKey {
			oa[k] = v
		}
	}
	return !reflect.DeepEqual(ra, oa), nil
}

func (r Ingress) UseTLSForRule(rule IngressRule) bool {
	if (rule.HTTP != nil && rule.HTTP.NoTLS) || (rule.TCP != nil && rule.TCP.NoTLS) {
		return false
	}
	_, ok := r.FindTLSSecret(rule.Host)
	return ok
}

func (r Ingress) FindTLSSecret(h string) (*LocalTypedReference, bool) {
	if h == "" {
		return nil, false
	}
	for _, tls := range r.Spec.TLS {
		for _, host := range tls.Hosts {
			if host == h || (h != "" && host == "*."+h) {
				return tls.Ref, true
			}
		}
	}
	return nil, false
}

func (r Ingress) IsPortChanged(o Ingress, cloudProvider string) bool {
	rpm, err := r.PortMappings(cloudProvider)
	if err != nil {
		return false
	}
	opm, err := o.PortMappings(cloudProvider)
	if err != nil {
		return false
	}
	return !reflect.DeepEqual(rpm, opm)
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
		(isMapKeyChanged(r.Annotations, o.Annotations, KeepSourceIP) || isMapKeyChanged(r.Annotations, o.Annotations, HealthCheckNodeport))
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
	for _, rule := range r.Spec.Rules {
		if rule.HTTP != nil {
			for _, svc := range rule.HTTP.Paths {
				record(svc.Backend.ServiceName)
			}
		} else if rule.TCP != nil {
			record(rule.TCP.Backend.ServiceName)
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
	for _, rule := range r.Spec.Rules {
		if rule.HTTP != nil {
			for _, svc := range rule.HTTP.Paths {
				if fqn(svc.Backend.ServiceName) == svcFQN {
					return true
				}
			}
		} else if rule.TCP != nil {
			if fqn(rule.TCP.Backend.ServiceName) == svcFQN {
				return true
			}
		}
	}
	return false
}

func (c Certificate) MatchesDomains(crt *x509.Certificate) bool {
	crtDomains := sets.NewString()
	if crt != nil {
		crtDomains.Insert(crt.Subject.CommonName)
		crtDomains.Insert(crt.DNSNames...)
	}
	return crtDomains.Equal(sets.NewString(c.Spec.Domains...))
}

func (c Certificate) ShouldRenew(crt *x509.Certificate) bool {
	bufferDays := c.Spec.RenewalBufferDays
	if bufferDays <= 0 {
		bufferDays = 15 // default 15 days
	}
	return !crt.NotAfter.After(time.Now().Add(time.Duration(bufferDays) * 24 * time.Hour))
}

func (c Certificate) IsRateLimited() bool {
	for _, cond := range c.Status.Conditions {
		if cond.Type == CertificateRateLimited {
			return time.Now().Add(-60 * time.Minute).Before(cond.LastUpdateTime.Time)
		}
	}
	return false
}

func (c Certificate) SecretName() string {
	if c.Spec.Storage.Vault != nil {
		if c.Spec.Storage.Vault.Name != "" {
			return c.Spec.Storage.Vault.Name
		}
		return "tls-" + c.Name
	}
	if c.Spec.Storage.Secret != nil && c.Spec.Storage.Secret.Name != "" {
		return c.Spec.Storage.Secret.Name
	}
	return "tls-" + c.Name
}

func (r Ingress) UsesAuthSecret(namespace, name string) bool {
	if r.Namespace != namespace {
		return false
	}

	if r.AuthSecretName() == name {
		return true
	}
	return false
}

// unify wildcard-host into empty-host, this will ease template rendering
// paths under empty-host and wildcard-host will be merged
// but TLS will be matched separately
func (r IngressRule) GetHost() string {
	host := strings.TrimSpace(r.Host)
	if host == `` || host == `*` {
		return ``
	}
	return host
}

func (r IngressRule) ParseALPNOptions() string {
	if r.HTTP != nil {
		if len(r.HTTP.ALPN) > 0 {
			return "alpn " + strings.Join(r.HTTP.ALPN, ",")
		} else if r.HTTP.Proto == "" {
			return "alpn http/1.1" // maintain backward compatibility
		}
	}
	if r.TCP != nil && len(r.TCP.ALPN) > 0 {
		return "alpn " + strings.Join(r.TCP.ALPN, ",")
	}
	return ""
}

func (r IngressBackend) ParseALPNOptions() string {
	if len(r.ALPN) > 0 {
		return "alpn " + strings.Join(r.ALPN, ",")
	}
	return ""
}

// check if both alpn and proto specified in the same rule/backend
func (r Ingress) ProtoWithALPN() error {
	// check default backend
	if r.Spec.Backend != nil {
		if r.Spec.Backend.Proto != "" && len(r.Spec.Backend.ALPN) > 0 {
			return fmt.Errorf("both proto and ALPN specified for spec.backend")
		}
	}
	for ri, rule := range r.Spec.Rules {
		if rule.HTTP != nil {
			if rule.HTTP.Proto != "" && len(rule.HTTP.ALPN) > 0 {
				return fmt.Errorf("both proto and ALPN specified for spec.rules[%d].http", ri)
			}
			for pi, path := range rule.HTTP.Paths {
				if path.Backend.Proto != "" && len(path.Backend.ALPN) > 0 {
					return fmt.Errorf("both proto and ALPN specified for spec.rules[%d].http.paths[%d]", ri, pi)
				}
			}
		}
		if rule.TCP != nil {
			if rule.TCP.Proto != "" && len(rule.TCP.ALPN) > 0 {
				return fmt.Errorf("both proto and ALPN specified for spec.rules[%d].tcp", ri)
			}
			if rule.TCP.Backend.Proto != "" && len(rule.TCP.Backend.ALPN) > 0 {
				return fmt.Errorf("both proto and ALPN specified for spec.rules[%d].tcp.backend", ri)
			}
		}
	}
	return nil
}

// check any alpn/proto contains "h2"
func (r Ingress) UseHTX() bool {
	// check default backend
	if r.Spec.Backend != nil {
		if containsH2(r.Spec.Backend.Proto, r.Spec.Backend.ALPN) {
			return true
		}
	}
	for _, rule := range r.Spec.Rules {
		if rule.HTTP != nil {
			if containsH2(rule.HTTP.Proto, rule.HTTP.ALPN) {
				return true
			}
			for _, path := range rule.HTTP.Paths {
				if containsH2(path.Backend.Proto, path.Backend.ALPN) {
					return true
				}
			}
		}
		if rule.TCP != nil {
			if containsH2(rule.TCP.Proto, rule.TCP.ALPN) {
				return true
			}
			if containsH2(rule.TCP.Backend.Proto, rule.TCP.Backend.ALPN) {
				return true
			}
		}
	}
	return false
}

func containsH2(proto string, alpn []string) bool {
	if proto == "h2" {
		return true
	}
	for _, option := range alpn {
		if option == "h2" {
			return true
		}
	}
	return false
}
