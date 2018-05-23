package v1beta1

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/appscode/kutil"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation"
)

// +k8s:openapi-gen=false
type indices struct {
	RuleIndex int
	PathIndex int
}

// +k8s:openapi-gen=false
type Paths map[string]indices

// +k8s:openapi-gen=false
type address struct {
	Protocol       string // tcp, http
	Address        string // IPv4, IPv6
	PodPort        int
	NodePort       int
	FirstRuleIndex int
	Hosts          map[string]Paths
}

func (a address) String() string {
	for h := range a.Hosts {
		return fmt.Sprintf("%s:%d", h, a.PodPort)
	}
	return fmt.Sprintf("%s:%d", a.Address, a.PodPort)
}

func (r *Ingress) Migrate() {
	for ti, tls := range r.Spec.TLS {
		if tls.SecretName != "" {
			r.Spec.TLS[ti].Ref = &LocalTypedReference{
				APIVersion: "v1",
				Kind:       "Secret",
				Name:       tls.SecretName,
			}
			r.Spec.TLS[ti].SecretName = ""
		}
	}
}

func (r Ingress) IsValid(cloudProvider string) error {
	for key, fn := range get {
		if _, err := fn(r.Annotations); err != nil && err != kutil.ErrNotFound {
			return errors.Errorf("can not parse annotation %s. Reason: %s", key, err)
		}
	}

	timeouts, _ := get[DefaultsTimeOut](r.Annotations)
	if err := checkMapKeys(timeouts.(map[string]string), sets.StringKeySet(timeoutDefaults)); err != nil {
		return errors.Errorf("invalid value for annotation %s. Reason: %s", DefaultsTimeOut, err)
	}

	for ri, rule := range r.Spec.FrontendRules {
		if _, err := checkRequiredPort(rule.Port); err != nil {
			return errors.Errorf("spec.frontendRules[%d].port %s is invalid. Reason: %s", ri, rule.Port, err)
		}
	}
	for ti, tls := range r.Spec.TLS {
		if tls.SecretName != "" {
			return errors.Errorf("spec.tls[%d].secretName must be migrated to spec.tls[%d].ref", ti, ti)
		} else if tls.Ref == nil {
			return errors.Errorf("spec.tls[%d] specifies no secret name and secret ref", ti)
		} else {
			if tls.Ref.Kind != "" && !(strings.EqualFold(tls.Ref.Kind, "Secret") || strings.EqualFold(tls.Ref.Kind, "Certificate")) {
				return errors.Errorf("spec.tls[%d].ref.kind %s is unsupported", ti, tls.Ref.Kind)
			}
			if tls.Ref.Name == "" {
				return errors.Errorf("spec.tls[%d] specifies no secret name and secret ref name", ti)
			}
		}
	}

	addrs := make(map[string]*address)
	nodePorts := make(map[int]int)
	usesHTTPRule := false
	sslPassthrough := r.SSLPassthrough()

	for ri, rule := range r.Spec.Rules {
		if rule.HTTP != nil && rule.TCP == nil {
			usesHTTPRule = true
			var err error
			var podPort, nodePort int
			podPort, err = checkOptionalPort(rule.HTTP.Port)
			if err != nil {
				return errors.Errorf("spec.rule[%d].http.port %s is invalid. Reason: %s", ri, rule.HTTP.Port, err)
			}
			if podPort == 0 { // detect port
				if r.UseTLSForRule(rule) {
					podPort = 443
				} else {
					podPort = 80
				}
			}
			if nodePort, err = checkOptionalPort(rule.HTTP.NodePort); err != nil {
				return errors.Errorf("spec.rule[%d].http.nodePort %s is invalid. Reason: %s", ri, rule.HTTP.NodePort, err)
			} else if nodePort > 0 {
				if r.LBType() == LBTypeHostPort || r.LBType() == LBTypeInternal {
					return errors.Errorf("spec.rule[%d].http.nodePort %s may not be specified when `LBType` is %s", ri, rule.HTTP.NodePort, r.LBType())
				}
			}
			bindAddress, err := checkOptionalAddress(rule.HTTP.Address)
			if err != nil {
				return errors.Errorf("spec.rule[%d].http.address %s is invalid. Reason: %s", ri, rule.HTTP.Address, err)
			} else if err = checkExclusiveWildcard(bindAddress, podPort, addrs); err != nil {
				return errors.Errorf("spec.rule[%d].http.address %s is invalid. Reason: %s", ri, rule.HTTP.Address, err)
			}

			var a *address
			var addrKey = fmt.Sprintf("%s:%d", bindAddress, podPort)

			if ea, found := addrs[addrKey]; found {
				if ea.Protocol == "tcp" {
					return errors.Errorf("spec.rule[%d].http is reusing port %d, also used in spec.rule[%d]", ri, ea.PodPort, ea.FirstRuleIndex)
				}
				if nodePort > 0 {
					if ea.NodePort > 0 && nodePort != ea.NodePort {
						return errors.Errorf("spec.rule[%d].http.nodePort %d does not match with nodePort %d", ri, nodePort, ea.NodePort)
					} else {
						nodePorts[nodePort] = ri
					}
				}

				// check for conflicting TLS
				if r.UseTLSForRule(rule) != r.UseTLSForRule(r.Spec.Rules[ea.FirstRuleIndex]) {
					return errors.Errorf("spec.rule[%d].http has conflicting TLS spec with spec.rule[%d].http", ri, ea.FirstRuleIndex)
				}

				// check for conflicting ALPN
				if rule.ParseALPNOptions() != r.Spec.Rules[ea.FirstRuleIndex].ParseALPNOptions() {
					return errors.Errorf("spec.rule[%d].HTTP has conflicting ALPN spec with spec.rule[%d].HTTP", ri, ea.FirstRuleIndex)
				}

				a = ea // paths will be merged into the original one
			} else {
				a = &address{
					Protocol:       "http",
					Address:        bindAddress,
					PodPort:        podPort,
					NodePort:       nodePort,
					FirstRuleIndex: ri,
					Hosts:          map[string]Paths{},
				}
				if nodePort > 0 {
					if ei, found := nodePorts[nodePort]; found {
						return errors.Errorf("spec.rule[%d].http is reusing nodePort %d for addr %s, also used in spec.rule[%d]", ri, nodePort, a, ei)
					} else {
						nodePorts[nodePort] = ri
					}
				}
				addrs[addrKey] = a
			}

			for pi, path := range rule.HTTP.Paths {
				if _, found := a.Hosts[rule.GetHost()]; !found {
					a.Hosts[rule.GetHost()] = Paths{}
				}
				if ei, found := a.Hosts[rule.GetHost()][path.Path]; found {
					return errors.Errorf("spec.rule[%d].http.paths[%d] is reusing path %s for addr %s, also used in spec.rule[%d].http.paths[%d]", ri, pi, path.Path, a, ei.RuleIndex, ei.PathIndex)
				}
				a.Hosts[rule.GetHost()][path.Path] = indices{RuleIndex: ri, PathIndex: pi}

				if !checkBackendServiceName(path.Backend.ServiceName) {
					return errors.Errorf("spec.rule[%d].http.paths[%d] has invalid serviceName for addr %s and path %s", ri, pi, a, path.Path)
				}
				if errs := validation.IsDNS1123Subdomain(path.Backend.ServiceName); len(errs) > 0 {
					return errors.Errorf("spec.rule[%d].http.paths[%d] is using invalid serviceName for addr %s. Reason: %s", ri, pi, a, strings.Join(errs, ","))
				}
				if _, err := checkRequiredPort(path.Backend.ServicePort); err != nil {
					return errors.Errorf("spec.rule[%d].http.paths[%d] is using invalid servicePort %s for addr %s and path %s. Reason: %s", ri, pi, path.Backend.ServicePort, a, path.Path, err)
				}
				for hi, hdr := range path.Backend.HeaderRules {
					if len(strings.Fields(hdr)) == 1 {
						return errors.Errorf("spec.rule[%d].http.paths[%d].backend.headerRules[%d] is invalid for addr %s and path %s", ri, pi, hi, a, path.Path)
					}
				}
			}
		} else if rule.TCP != nil && rule.HTTP == nil {
			var err error
			var podPort, nodePort int

			if podPort, err = checkRequiredPort(rule.TCP.Port); err != nil {
				return errors.Errorf("spec.rule[%d].tcp.port %s is invalid. Reason: %s", ri, rule.TCP.Port, err)
			}
			if nodePort, err = checkOptionalPort(rule.TCP.NodePort); err != nil {
				return errors.Errorf("spec.rule[%d].tcp.nodePort %s is invalid. Reason: %s", ri, rule.TCP.NodePort, err)
			} else if nodePort > 0 {
				if r.LBType() == LBTypeHostPort || r.LBType() == LBTypeInternal {
					return errors.Errorf("spec.rule[%d].tcp.nodePort %s may not be specified when `LBType` is %s", ri, rule.TCP.NodePort, r.LBType())
				}
			}
			bindAddress, err := checkOptionalAddress(rule.TCP.Address)
			if err != nil {
				return errors.Errorf("spec.rule[%d].tcp.address %s is invalid. Reason: %s", ri, rule.TCP.Address, err)
			} else if err = checkExclusiveWildcard(bindAddress, podPort, addrs); err != nil {
				return errors.Errorf("spec.rule[%d].tcp.address %s is invalid. Reason: %s", ri, rule.TCP.Address, err)
			}

			// should not use TLS in passthrough mode
			useTLS := r.UseTLSForRule(rule)
			if sslPassthrough && useTLS {
				return errors.Errorf("TLS defined for spec.rule[%d].tcp in SSLPassthrough mode", ri)
			}

			var a *address
			var addrKey = fmt.Sprintf("%s:%d", bindAddress, podPort)

			if ea, found := addrs[addrKey]; found {
				if ea.Protocol != "tcp" {
					return errors.Errorf("spec.rule[%d].tcp is reusing port %d, also used in spec.rule[%d].http", ri, ea.PodPort, ea.FirstRuleIndex)
				}
				if nodePort != ea.NodePort {
					return errors.Errorf("spec.rule[%d].tcp.nodePort %d does not match with nodePort %d", ri, nodePort, ea.NodePort)
				} else {
					nodePorts[nodePort] = ri
				}

				// for empty-host and wildcard-host, there can not be more than one rules under same address-binder
				if rule.Host == "" || rule.Host == "*" { // current host is wildcard-host but previously found one/more rules
					return errors.Errorf("multiple rules with one/more wildcard/empty host found for address %s", addrKey)
				} else { // current host is not wildcard-host but previously found rules with with wildcard-host
					if _, found := ea.Hosts[""]; found {
						return errors.Errorf("multiple rules with one/more wildcard/empty host found for address %s", addrKey)
					}
					if _, found := ea.Hosts["*"]; found {
						return errors.Errorf("multiple rules with one/more wildcard/empty host found for address %s", addrKey)
					}
				}

				// check for conflicting TLS
				if useTLS != r.UseTLSForRule(r.Spec.Rules[ea.FirstRuleIndex]) {
					return errors.Errorf("spec.rule[%d].TCP has conflicting TLS spec with spec.rule[%d].TCP", ri, ea.FirstRuleIndex)
				}

				// check for conflicting ALPN
				if rule.ParseALPNOptions() != r.Spec.Rules[ea.FirstRuleIndex].ParseALPNOptions() {
					return errors.Errorf("spec.rule[%d].TCP has conflicting ALPN spec with spec.rule[%d].TCP", ri, ea.FirstRuleIndex)
				}

				a = ea
			} else {
				a = &address{
					Protocol:       "tcp",
					Address:        bindAddress,
					PodPort:        podPort,
					FirstRuleIndex: ri,
					Hosts:          map[string]Paths{},
				}
				if nodePort > 0 {
					if ei, found := nodePorts[nodePort]; found {
						return errors.Errorf("spec.rule[%d].tcp is reusing nodePort %d for addr %s, also used in spec.rule[%d]", ri, nodePort, a, ei)
					} else {
						a.NodePort = nodePort
						nodePorts[nodePort] = ri
					}
				}
				addrs[addrKey] = a
			}

			if _, found := a.Hosts[rule.GetHost()]; !found {
				a.Hosts[rule.GetHost()] = Paths{
					"": indices{RuleIndex: ri}, // for tcp no paths, just store indices in empty-path
				}
			} else { // same host under same address-binder
				ei := a.Hosts[rule.GetHost()][""]
				return errors.Errorf("spec.rule[%d].tcp is reusing host %s for addr %s, also used in spec.rule[%d]", ri, rule.Host, a, ei.RuleIndex)
			}

			if !checkBackendServiceName(rule.TCP.Backend.ServiceName) {
				return errors.Errorf("spec.rule[%d].tcp has invalid serviceName for addr %s", ri, a)
			}
			if errs := validation.IsDNS1123Subdomain(rule.TCP.Backend.ServiceName); len(errs) > 0 {
				return errors.Errorf("spec.rule[%d].tcp is using invalid serviceName for addr %s. Reason: %s", ri, a, strings.Join(errs, ","))
			}
			if _, err := checkRequiredPort(rule.TCP.Backend.ServicePort); err != nil {
				return errors.Errorf("spec.rule[%d].tcp is using invalid servicePort %s for addr %s. Reason: %s", ri, rule.TCP.Backend.ServicePort, a, err)
			}
		} else if rule.TCP == nil && rule.HTTP == nil {
			return errors.Errorf("spec.rule[%d] is missing both HTTP and TCP specification", ri)
		} else {
			return errors.Errorf("spec.rule[%d] can specify either HTTP or TCP", ri)
		}
	}

	// If Ingress does not use any HTTP rule but defined a default backend, we need to open port 80
	if !usesHTTPRule && r.Spec.Backend != nil {
		if !checkBackendServiceName(r.Spec.Backend.ServiceName) {
			return errors.Errorf("invalid serviceName for default backend")
		}
		addrs["*:80"] = &address{Protocol: "http", Address: "*", PodPort: 80}
	}
	// ref: https://github.com/appscode/voyager/issues/188
	if cloudProvider == "aws" && r.LBType() == LBTypeLoadBalancer {
		if ans, ok := r.ServiceAnnotations(cloudProvider); ok {
			if v, usesAWSCertManager := ans["service.beta.kubernetes.io/aws-load-balancer-ssl-cert"]; usesAWSCertManager && v != "" {
				var tp80, sp443 bool
				for ip, target := range addrs {
					svcPort := strings.Split(ip, ":")[1]
					if target.PodPort == 80 {
						tp80 = true
					}
					if svcPort == "443" {
						sp443 = true
					}
				}
				if !tp80 || sp443 {
					return errors.Errorf("failed to open port 443 on service for AWS cert manager for Ingress %s/%s", r.Namespace, r.Name)
				}
			}
		}
	}
	if !r.SupportsLBType(cloudProvider) {
		return errors.Errorf("ingress %s/%s uses unsupported LBType %s for cloud provider %s", r.Namespace, r.Name, r.LBType(), cloudProvider)
	}

	if (r.LBType() == LBTypeNodePort || r.LBType() == LBTypeHostPort || r.LBType() == LBTypeInternal) && len(r.Spec.LoadBalancerSourceRanges) > 0 {
		return errors.Errorf("ingress %s/%s of type %s can't use `spec.LoadBalancerSourceRanges`", r.Namespace, r.Name, r.LBType())
	}

	// validate external auth
	for ri, rule := range r.Spec.FrontendRules {
		if rule.Auth != nil && rule.Auth.OAuth != nil {
			oauthHosts := make(map[string]int)
			for ii, oauth := range rule.Auth.OAuth {
				// check multiple oauth for same host under same port
				if jj, found := oauthHosts[oauth.Host]; found {
					return errors.Errorf("spec.frontendRules[%d].oauth[%d] is reusing host %s for port %s, also used in spec.frontendRules[%d].oauth[%d]", ri, ii, rule.Auth.OAuth[ii].Host, rule.Port, ri, jj)
				} else {
					oauthHosts[oauth.Host] = ii
				}

				// check auth backend exists
				authBackendFound := false
				for _, addr := range addrs {
					if addr.PodPort == int(rule.Port.IntVal) {
						for _, path := range addr.Hosts[oauth.Host] {
							if r.Spec.Rules[path.RuleIndex].HTTP.Paths[path.PathIndex].Backend.Name == oauth.AuthBackend {
								authBackendFound = true
							}
						}
					}
				}
				if !authBackendFound {
					return errors.Errorf("specified auth backend not found for spec.frontendRules[%d].oauth[%d]", ri, ii)
				}
			}
		}
	}

	return nil
}

func (r Ingress) SupportsLBType(cloudProvider string) bool {
	switch r.LBType() {
	case LBTypeLoadBalancer:
		return cloudProvider == "aws" ||
			cloudProvider == "gce" ||
			cloudProvider == "gke" ||
			cloudProvider == "azure" ||
			cloudProvider == "acs" ||
			cloudProvider == "openstack" ||
			cloudProvider == "minikube" ||
			cloudProvider == "metallb"
	case LBTypeNodePort:
		return cloudProvider != "acs"
	case LBTypeHostPort:
		// TODO: https://github.com/appscode/voyager/issues/374
		return cloudProvider != "acs" && cloudProvider != "azure" && cloudProvider != "gce" && cloudProvider != "gke"
	case LBTypeInternal:
		return true
	default:
		return false
	}
}

func checkBackendServiceName(name string) bool {
	if strings.Contains(name, ".") {
		idx := strings.Index(name, ".")
		return name[:idx] != "" && name[idx+1:] != ""
	} else {
		return name != ""
	}
}

func checkRequiredPort(port intstr.IntOrString) (int, error) {
	if port.Type == intstr.Int {
		if port.IntVal <= 0 {
			return 0, errors.Errorf("port %s must a +ve integer", port)
		}
		return int(port.IntVal), nil
	} else if port.Type == intstr.String {
		return strconv.Atoi(port.StrVal)
	}
	return 0, errors.Errorf("invalid data type %v for port %s", port.Type, port)
}

func checkOptionalPort(port intstr.IntOrString) (int, error) {
	if port.Type == intstr.Int {
		if port.IntVal < 0 {
			return 0, errors.Errorf("port %s can't be -ve integer", port)
		}
		return int(port.IntVal), nil
	} else if port.Type == intstr.String {
		if port.StrVal == "" {
			return 0, nil
		}
		return strconv.Atoi(port.StrVal)
	}
	return 0, errors.Errorf("invalid data type %v for port %s", port.Type, port)
}

func checkOptionalAddress(address string) (string, error) {
	if address != "" && net.ParseIP(address) == nil {
		return "", errors.Errorf("could not parse IPv4 or IPv6 address")
	} else if address == "" {
		return "*", nil
	}

	return address, nil
}

func checkExclusiveWildcard(address string, port int, defined map[string]*address) error {
	var wildcard = fmt.Sprintf("*:%d", port)

	if address == "*" {
		// If a wildcard already exists for the port, we've passed validation for this port before.
		if _, ok := defined[wildcard]; ok {
			return nil
		}

		// Check defined addresses for existing bind against the specified port and a non-wildcard IP.
		for i := range defined {
			if defined[i].PodPort == port && defined[i].Address != "*" {
				return errors.Errorf("cannot use wildcard address for port %d, bind already exists for address %s", port, defined[i].Address)
			}
		}
	} else if _, ok := defined[wildcard]; ok {
		return errors.Errorf("cannot use address %s for port %d, one or more rules use a wildcard bind address", address, port)
	}

	return nil
}

func (c Certificate) IsValid(cloudProvider string) error {
	if len(c.Spec.Domains) == 0 {
		return errors.Errorf("domain list is empty")
	}

	if c.Spec.ChallengeProvider.HTTP == nil && c.Spec.ChallengeProvider.DNS == nil {
		return errors.Errorf("certificate has no valid challenge provider")
	}

	if c.Spec.ChallengeProvider.HTTP != nil && c.Spec.ChallengeProvider.DNS != nil {
		return errors.Errorf("invalid provider specification, used both http and dns provider")
	}

	if c.Spec.ChallengeProvider.HTTP != nil {
		if len(c.Spec.ChallengeProvider.HTTP.Ingress.Name) == 0 ||
			(c.Spec.ChallengeProvider.HTTP.Ingress.APIVersion != SchemeGroupVersion.String() && c.Spec.ChallengeProvider.HTTP.Ingress.APIVersion != "extensions/v1beta1") {
			return errors.Errorf("invalid ingress reference")
		}
	}

	if c.Spec.ChallengeProvider.DNS != nil {
		if len(c.Spec.ChallengeProvider.DNS.Provider) == 0 {
			return errors.Errorf("no dns provider name specified")
		}
		if c.Spec.ChallengeProvider.DNS.CredentialSecretName == "" {
			useCredentialFromEnv := (cloudProvider == "aws" && sets.NewString("aws", "route53").Has(c.Spec.ChallengeProvider.DNS.Provider) && c.Spec.ChallengeProvider.DNS.CredentialSecretName == "") ||
				(sets.NewString("gce", "gke").Has(cloudProvider) && sets.NewString("googlecloud", "gcloud", "gce", "gke").Has(c.Spec.ChallengeProvider.DNS.Provider) && c.Spec.ChallengeProvider.DNS.CredentialSecretName == "")
			if !useCredentialFromEnv {
				return errors.Errorf("missing dns challenge provider credential")
			}
		}
	}

	if len(c.Spec.ACMEUserSecretName) == 0 {
		return errors.Errorf("no user secret name specified")
	}

	if c.Spec.Storage.Secret != nil && c.Spec.Storage.Vault != nil {
		return errors.Errorf("invalid storage specification, used both storage")
	}

	return nil
}

func checkMapKeys(m map[string]string, keys sets.String) error {
	diff := sets.StringKeySet(m).Difference(keys)
	if diff.Len() != 0 {
		return errors.Errorf("invalid keys: %v", diff.List())
	}
	return nil
}
