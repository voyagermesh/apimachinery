package v1beta1

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/appscode/kutil/meta"
	"github.com/appscode/voyager/apis/voyager"
)

const (
	IngressKey = "ingress.kubernetes.io"
	EngressKey = "ingress.appscode.com"

	APISchema        = EngressKey + "/" + "api-schema" // APISchema = {APIGroup}/{APIVersion}
	APISchemaEngress = voyager.GroupName + "/v1beta1"
	APISchemaIngress = "extension/v1beta1"

	VoyagerPrefix = "voyager-"

	// LB stats options
	StatsOn          = EngressKey + "/" + "stats"
	StatsPort        = EngressKey + "/" + "stats-port"
	StatsSecret      = EngressKey + "/" + "stats-secret-name"
	StatsServiceName = EngressKey + "/" + "stats-service-name"
	DefaultStatsPort = 56789

	LBTypeHostPort     = "HostPort"
	LBTypeNodePort     = "NodePort"
	LBTypeLoadBalancer = "LoadBalancer" // default
	LBTypeInternal     = "Internal"
	LBType             = EngressKey + "/" + "type"

	// Runs HAProxy on a specific set of a hosts.
	NodeSelector = EngressKey + "/" + "node-selector"

	// Replicas specify # of HAProxy pods run (default 1)
	Replicas = EngressKey + "/" + "replicas"

	// IP to be assigned to cloud load balancer
	LoadBalancerIP = EngressKey + "/" + "load-balancer-ip" // IP or empty

	// BackendWeight is the weight value of a Pod that was
	// addressed by the Endpoint, this weight will be added to server backend.
	// Traffic will be forwarded according to there weight.
	BackendWeight = EngressKey + "/" + "backend-weight"

	// https://github.com/appscode/voyager/issues/103
	// ServiceAnnotations is user provided annotations map that will be
	// applied to the service of that LoadBalancer.
	// ex: "ingress.appscode.com/annotations-service": {"key": "val"}
	ServiceAnnotations = EngressKey + "/" + "annotations-service"

	// PodAnnotations is user provided annotations map that will be
	// applied to the Pods (Deployment/ DaemonSet) of that LoadBalancer.
	// ex: "ingress.appscode.com/annotations-pod": {"key": "val"}
	PodAnnotations = EngressKey + "/" + "annotations-pod"

	// Preserves source IP for LoadBalancer type ingresses. The actual configuration
	// generated depends on the underlying cloud provider.
	//
	//  - gce, gke, azure: Adds annotation service.beta.kubernetes.io/external-traffic: OnlyLocal
	// to services used to expose HAProxy.
	// ref: https://kubernetes.io/docs/tutorials/services/source-ip/#source-ip-for-services-with-typeloadbalancer
	//
	// - aws: Enforces the use of the PROXY protocol over any connection accepted by any of
	// the sockets declared on the same line. Versions 1 and 2 of the PROXY protocol
	// are supported and correctly detected. The PROXY protocol dictates the layer
	// 3/4 addresses of the incoming connection to be used everywhere an address is
	// used, with the only exception of "tcp-request connection" rules which will
	// only see the real connection address. Logs will reflect the addresses
	// indicated in the protocol, unless it is violated, in which case the real
	// address will still be used.  This keyword combined with support from external
	// components can be used as an efficient and reliable alternative to the
	// X-Forwarded-For mechanism which is not always reliable and not even always
	// usable. See also "tcp-request connection expect-proxy" for a finer-grained
	// setting of which client is allowed to use the protocol.
	// ref: https://github.com/kubernetes/kubernetes/blob/release-1.5/pkg/cloudprovider/providers/aws/aws.go#L79
	KeepSourceIP = EngressKey + "/" + "keep-source-ip"

	// Enforces the use of the PROXY protocol over any connection accepted by HAProxy.
	AcceptProxy = EngressKey + "/" + "accept-proxy"

	// Enforces use of the PROXY protocol over any connection established to this server.
	// Possible values are "v1", "v2", "v2-ssl" and "v2-ssl-cn"
	SendProxy = EngressKey + "/" + "send-proxy"

	// Annotations applied to resources offshoot from an ingress
	OriginAPISchema = EngressKey + "/" + "origin-api-schema" // APISchema = {APIGroup}/{APIVersion}
	OriginName      = EngressKey + "/" + "origin-name"

	// https://github.com/appscode/voyager/issues/280
	// Supports all valid timeout option for defaults section of HAProxy
	// https://cbonte.github.io/haproxy-dconv/1.7/configuration.html#4.2-timeout%20check
	// expects a json encoded map
	// ie: "ingress.appscode.com/default-timeout": {"client": "5s"}
	//
	// If the annotation is not set default values used to config defaults section will be:
	//
	// timeout  connect         50s
	// timeout  client          50s
	// timeout  client-fin      50s
	// timeout  server          50s
	// timeout  tunnel          50s
	DefaultsTimeOut = EngressKey + "/" + "default-timeout"

	// https://github.com/appscode/voyager/issues/343
	// Supports all valid options for defaults section of HAProxy config
	// https://cbonte.github.io/haproxy-dconv/1.7/configuration.html#4.2-option%20abortonclose
	// from the list from here
	// expects a json encoded map
	// ie: "ingress.appscode.com/default-option": '{"http-keep-alive": "true", "dontlognull": "true", "clitcpka": "false"}'
	// This will be appended in the defaults section of HAProxy as
	//
	//   option http-keep-alive
	//   option dontlognull
	//   no option clitcpka
	//
	DefaultsOption = EngressKey + "/" + "default-option"

	// Available Options
	//   ssl:
	//    Creates a TLS/SSL socket when connecting to this server in order to cipher/decipher the traffic
	//
	//    if verify not set the following error may occurred
	//    [/etc/haproxy/haproxy.cfg:49] verify is enabled by default but no CA file specified.
	//    If you're running on a LAN where you're certain to trust the server's certificate,
	//    please set an explicit 'verify none' statement on the 'server' line, or use
	//    'ssl-server-verify none' in the global section to disable server-side verifications by default.
	//
	//   verify [none|required]:
	//    Sets HAProxy‘s behavior regarding the certificated presented by the server:
	//   none :
	//    doesn’t verify the certificate of the server
	//
	//   required (default value) :
	//    TLS handshake is aborted if the validation of the certificate presented by the server returns an error.
	//
	//   veryfyhost <hostname>:
	//    Sets a <hostname> to look for in the Subject and SubjectAlternateNames fields provided in the
	//    certificate sent by the server. If <hostname> can’t be found, then the TLS handshake is aborted.
	// ie.
	// ingress.appscode.com/backend-tls: "ssl verify none"
	//
	// If this annotation is not set HAProxy will connect to backend as http,
	// This value should not be set if the backend do not support https resolution.
	BackendTLSOptions = EngressKey + "/backend-tls"

	// StickyIngress configures HAProxy to use sticky connection
	// to the backend servers.
	// Annotations could  be applied to either Ingress or backend Service (since 3.2+).
	// ie: ingress.appscode.com/sticky-session: "true"
	// If applied to Ingress, all the backend connections would be sticky
	// If applied to Service and Ingress do not have this annotation only
	// connection to that backend service will be sticky.
	// Deprecated
	StickySession = EngressKey + "/" + "sticky-session"
	// Specify a method to stick clients to origins across requests.
	// Only supported value is cookie.
	IngressAffinity = IngressKey + "/affinity"
	// When affinity is set to cookie, the name of the cookie to use.
	IngressAffinitySessionCookieName = IngressKey + "/session-cookie-name"
	// When affinity is set to cookie, the hash algorithm used: md5, sha, index.
	IngressAffinitySessionCookieHash = IngressKey + "/session-cookie-hash"

	// Basic Auth: Follows ingress controller standard
	// https://github.com/kubernetes/ingress/tree/master/examples/auth/basic/haproxy
	// HAProxy Ingress read user and password from auth file stored on secrets, one
	// user and password per line.
	// Each line of the auth file should have:
	// user and insecure password separated with a pair of colons: <username>::<plain-text-passwd>; or
	// user and an encrypted password separated with colons: <username>:<encrypted-passwd>
	// Secret name, realm and type are configured with annotations in the ingress
	// Auth can only be applied to HTTP backends.
	// Only supported type is basic
	AuthType = IngressKey + "/auth-type"

	// an optional string with authentication realm
	AuthRealm = IngressKey + "/auth-realm"

	// name of the auth secret
	AuthSecret = IngressKey + "/auth-secret"

	// Name of secret for TLS client certification validation.
	AuthTLSSecret = IngressKey + "/auth-tls-secret"

	// The page that user should be redirected in case of Auth error
	AuthTLSErrorPage = IngressKey + "/auth-tls-error-page"

	// Enables verification of client certificates.
	AuthTLSVerifyClient = IngressKey + "/auth-tls-verify-client"

	// Enables CORS headers in response.
	// Setting this annotations in ingress will add CORS headers to all HTTP
	// frontend. If we need to add cors headers only on specific frontend we can also
	// configure this using FrontendRules for specific frontend.
	// http://blog.nasrulhazim.com/2017/07/haproxy-setting-up-cors/
	CORSEnabled = IngressKey + "/enable-cors"

	// Maximum http request body size. This returns the advertised length of the HTTP request's body in bytes. It
	// will represent the advertised Content-Length header
	// http://cbonte.github.io/haproxy-dconv/1.7/configuration.html#7.3.6-req.body_size
	//
	ProxyBodySize = IngressKey + "/proxy-body-size"

	// Pass TLS connections directly to backend; do not offload.
	SSLPassthrough = IngressKey + "/ssl-passthrough"

	EnableHSTS = IngressKey + "/hsts"
	// This specifies the time (in seconds) the browser should connect to the server using the HTTPS connection.
	// https://blog.stackpath.com/glossary/hsts/
	HSTSMaxAge  = IngressKey + "/hsts-max-age"
	HSTSPreload = IngressKey + "/hsts-preload"
	// If specified, this HSTS rule applies to all of the site's subdomains as well.
	HSTSIncludeSubDomains = IngressKey + "/hsts-include-subdomains"

	WhitelistSourceRange = IngressKey + "/whitelist-source-range"
	MaxConnections       = IngressKey + "/max-connections"

	// https://github.com/appscode/voyager/issues/552
	ForceServicePort = EngressKey + "/force-service-port"
	SSLRedirect      = IngressKey + "/ssl-redirect"
	ForceSSLRedirect = IngressKey + "/force-ssl-redirect"

	// https://github.com/appscode/voyager/issues/525
	ErrorFiles = EngressKey + "/errorfiles"

	// Limit requests per second per IP address
	// http://cbonte.github.io/haproxy-dconv/1.8/configuration.html#7.3.3-sc_conn_rate
	// https://serverfault.com/a/679172/349346
	// https://discourse.haproxy.org/t/solved-how-to-configure-basic-ddos-protection-when-behind-aws-elb-x-forwarded-for/932
	// https://www.haproxy.com/blog/use-a-load-balancer-as-a-first-row-of-defense-against-ddos/
	LimitRPS = IngressKey + "/limit-rps"
	// Limit requests per minute per IP address
	LimitRPM = IngressKey + "/limit-rpm"

	// http://cbonte.github.io/haproxy-dconv/1.8/configuration.html#7.3.3-src_conn_cur
	// https://www.haproxy.com/blog/use-a-load-balancer-as-a-first-row-of-defense-against-ddos/
	LimitConnection = IngressKey + "/limit-connection"

	// https://github.com/appscode/voyager/issues/683
	// https://www.haproxy.com/documentation/aloha/7-0/haproxy/healthchecks/
	CheckHealth     = EngressKey + "/" + "check"
	CheckHealthPort = EngressKey + "/" + "check-port"
)

const (
	ACMEUserEmail        = "ACME_EMAIL"
	ACMEUserPrivatekey   = "ACME_USER_PRIVATE_KEY"
	ACMERegistrationData = "ACME_REGISTRATION_DATA"
	ACMEServerURL        = "ACME_SERVER_URL"
)

type ProxyProtocolVersion string

const (
	proxyProtocolV1      ProxyProtocolVersion = "v1"
	proxyProtocolV2      ProxyProtocolVersion = "v2"
	proxyProtocolV2SSL   ProxyProtocolVersion = "v2-ssl"
	proxyProtocolV2SSLCN ProxyProtocolVersion = "v2-ssl-cn"
)

func ProxyProtocolCommand(version string) string {
	switch ProxyProtocolVersion(version) {
	case proxyProtocolV1:
		return "send-proxy"
	case proxyProtocolV2, proxyProtocolV2SSL, proxyProtocolV2SSLCN:
		return "send-proxy-" + version
	default:
		return ""
	}
}

func (r Ingress) OffshootName() string {
	return VoyagerPrefix + r.Name
}

func (r Ingress) OffshootLabels() map[string]string {
	lbl := map[string]string{
		"origin":      "voyager",
		"origin-name": r.Name,
	}

	gv := strings.SplitN(r.APISchema(), "/", 2)
	if len(gv) == 2 {
		lbl["origin-api-group"] = gv[0]
	}
	return lbl
}

func (r Ingress) StatsLabels() map[string]string {
	lbl := r.OffshootLabels()
	lbl["feature"] = "stats"
	return lbl
}

func (r Ingress) APISchema() string {
	if v := meta.GetString(r.Annotations, APISchema); v != "" {
		return v
	}
	return APISchemaEngress
}

func (r Ingress) Sticky() bool {
	// Specify a method to stick clients to origins across requests.
	// Like nginx HAProxy only supports the value cookie.
	if len(meta.GetString(r.Annotations, IngressAffinity)) > 0 {
		return true
	}
	v, _ := meta.GetBool(r.Annotations, StickySession)
	return v
}

func (r Ingress) StickySessionCookieName() string {
	// When affinity is set to cookie, the name of the cookie to use.
	cookieName := meta.GetString(r.Annotations, IngressAffinitySessionCookieName)
	if len(cookieName) > 0 {
		return cookieName
	}
	return "SERVERID"
}

func (r Ingress) StickySessionCookieHashType() string {
	return meta.GetString(r.Annotations, IngressAffinitySessionCookieHash)
}

func (r Ingress) EnableCORS() bool {
	v, _ := meta.GetBool(r.Annotations, CORSEnabled)
	return v
}

func (r Ingress) ForceServicePort() bool {
	v, _ := meta.GetBool(r.Annotations, ForceServicePort)
	return v
}

func (r Ingress) EnableHSTS() bool {
	v, err := meta.GetBool(r.Annotations, EnableHSTS)
	if err != nil {
		return true // enable hsts by default
	}
	return v
}

func (r Ingress) HSTSMaxAge() int {
	v := meta.GetString(r.Annotations, HSTSMaxAge)
	ageInSec, err := strconv.Atoi(v)
	if err == nil {
		return ageInSec
	}
	d, err := time.ParseDuration(v)
	if err == nil {
		return int(d.Seconds())
	}
	// default 6 months
	return 15768000
}

func (r Ingress) HSTSPreload() bool {
	v, _ := meta.GetBool(r.Annotations, HSTSPreload)
	return v
}

func (r Ingress) HSTSIncludeSubDomains() bool {
	v, _ := meta.GetBool(r.Annotations, HSTSIncludeSubDomains)
	return v
}

func (r Ingress) WhitelistSourceRange() string {
	return meta.GetString(r.Annotations, WhitelistSourceRange)
}

func (r Ingress) MaxConnections() int {
	v, _ := meta.GetInt(r.Annotations, MaxConnections)
	return v
}

func (r Ingress) SSLRedirect() bool {
	v, err := meta.GetBool(r.Annotations, SSLRedirect)
	if err == nil && !v {
		return false
	}
	return true
}

func (r Ingress) ForceSSLRedirect() bool {
	v, _ := meta.GetBool(r.Annotations, ForceSSLRedirect)
	return v
}

func (r Ingress) ProxyBodySize() string {
	return meta.GetString(r.Annotations, ProxyBodySize)
}

func (r Ingress) SSLPassthrough() bool {
	v, _ := meta.GetBool(r.Annotations, SSLPassthrough)
	return v
}

func (r Ingress) Stats() bool {
	v, _ := meta.GetBool(r.Annotations, StatsOn)
	return v
}

func (r Ingress) StatsSecretName() string {
	return meta.GetString(r.Annotations, StatsSecret)
}

func (r Ingress) StatsPort() int {
	if v, _ := meta.GetInt(r.Annotations, StatsPort); v > 0 {
		return v
	}
	return DefaultStatsPort
}

func (r Ingress) StatsServiceName() string {
	if v := meta.GetString(r.Annotations, StatsServiceName); v != "" {
		return v
	}
	return VoyagerPrefix + r.Name + "-stats"
}

func (r Ingress) LBType() string {
	if v := meta.GetString(r.Annotations, LBType); v != "" {
		return v
	}
	return LBTypeLoadBalancer
}

func (r Ingress) Replicas() int32 {
	if v, _ := meta.GetInt(r.Annotations, Replicas); v > 0 {
		return int32(v)
	}
	return 1
}

func (r Ingress) NodeSelector() map[string]string {
	if v, _ := meta.GetMap(r.Annotations, NodeSelector); len(v) > 0 {
		return v
	}
	return ParseDaemonNodeSelector(meta.GetString(r.Annotations, EngressKey+"/"+"daemon.nodeSelector"))
}

func (r Ingress) LoadBalancerIP() net.IP {
	if v := meta.GetString(r.Annotations, LoadBalancerIP); v != "" {
		return net.ParseIP(v)
	}
	return nil
}

func (r Ingress) ServiceAnnotations(provider string) (map[string]string, bool) {
	ans, err := meta.GetMap(r.Annotations, ServiceAnnotations)
	if err == nil {
		filteredMap := make(map[string]string)
		for k, v := range ans {
			if !strings.HasPrefix(strings.TrimSpace(k), EngressKey+"/") {
				filteredMap[k] = v
			}
		}
		if r.LBType() == LBTypeLoadBalancer && r.KeepSourceIP() {
			switch provider {
			case "aws":
				// ref: https://github.com/kubernetes/kubernetes/blob/release-1.5/pkg/cloudprovider/providers/aws/aws.go#L79
				filteredMap["service.beta.kubernetes.io/aws-load-balancer-proxy-protocol"] = "*"
			}
		}
		return filteredMap, true
	}
	return ans, false
}

func (r Ingress) PodsAnnotations() (map[string]string, bool) {
	ans, err := meta.GetMap(r.Annotations, PodAnnotations)
	if err == nil {
		filteredMap := make(map[string]string)
		for k, v := range ans {
			if !strings.HasPrefix(strings.TrimSpace(k), EngressKey+"/") {
				filteredMap[k] = v
			}
		}
		return filteredMap, true
	}
	return ans, false
}

func (r Ingress) KeepSourceIP() bool {
	v, _ := meta.GetBool(r.Annotations, KeepSourceIP)
	return v
}

func (r Ingress) AcceptProxy() bool {
	v, _ := meta.GetBool(r.Annotations, AcceptProxy)
	return v
}

var timeoutDefaults = map[string]string{
	// Maximum time to wait for a connection attempt to a server to succeed.
	"connect": "50s",

	// Maximum inactivity time on the client side.
	// Applies when the client is expected to acknowledge or send data.
	"client": "50s",

	// Inactivity timeout on the client side for half-closed connections.
	// Applies when the client is expected to acknowledge or send data
	// while one direction is already shut down.
	"client-fin": "50s",

	// Maximum inactivity time on the server side.
	"server": "50s",

	// Timeout to use with WebSocket and CONNECT
	"tunnel": "50s",
}

func (r Ingress) Timeouts() map[string]string {
	ans, _ := meta.GetMap(r.Annotations, DefaultsTimeOut)
	if ans == nil {
		ans = make(map[string]string)
	}

	// If the timeouts specified in `defaultTimeoutValues` are not set specifically set
	// we need to set default timeout values.
	// An unspecified timeout results in an infinite timeout, which
	// is not recommended. Such a usage is accepted and works but reports a warning
	// during startup because it may results in accumulation of expired sessions in
	// the system if the system's timeouts are not configured either.
	for k, v := range timeoutDefaults {
		if _, ok := ans[k]; !ok {
			ans[k] = v
		}
	}

	return ans
}

func (r Ingress) HAProxyOptions() map[string]bool {
	ans, _ := meta.GetMap(r.Annotations, DefaultsOption)
	if ans == nil {
		ans = make(map[string]string)
	}

	ret := make(map[string]bool)
	for k := range ans {
		val, err := meta.GetBool(ans, k)
		if err != nil {
			continue
		}
		ret[k] = val
	}

	if len(ret) == 0 {
		ret["http-server-close"] = true
		ret["dontlognull"] = true
	}

	return ret
}

func (r Ingress) BasicAuthEnabled() bool {
	if r.Annotations == nil {
		return false
	}

	// Check auth type is basic; other auth mode is not supported
	if val := meta.GetString(r.Annotations, AuthType); val != "basic" {
		return false
	}

	// Check secret name is not empty
	if val := meta.GetString(r.Annotations, AuthSecret); len(val) == 0 {
		return false
	}

	return true
}

func (r Ingress) AuthRealm() string {
	return meta.GetString(r.Annotations, AuthRealm)
}

func (r Ingress) AuthSecretName() string {
	return meta.GetString(r.Annotations, AuthSecret)
}

func (r Ingress) AuthTLSSecret() string {
	return meta.GetString(r.Annotations, AuthTLSSecret)
}

func (r Ingress) AuthTLSVerifyClient() TLSAuthVerifyOption {
	str := meta.GetString(r.Annotations, AuthTLSVerifyClient)
	if str == string(TLSAuthVerifyOptional) {
		return TLSAuthVerifyOptional
	}
	return TLSAuthVerifyRequired
}

func (r Ingress) AuthTLSErrorPage() string {
	return meta.GetString(r.Annotations, AuthTLSErrorPage)
}

func (r Ingress) ErrorFilesConfigMapName() string {
	return meta.GetString(r.Annotations, ErrorFiles)
}

func (r Ingress) LimitRPS() int {
	value, _ := meta.GetInt(r.Annotations, LimitRPS)
	return value
}

func (r Ingress) LimitRPM() int {
	value, _ := meta.GetInt(r.Annotations, LimitRPM)
	return value
}

func (r Ingress) LimitConnections() int {
	value, _ := meta.GetInt(r.Annotations, LimitConnection)
	return value
}

// ref: https://github.com/kubernetes/kubernetes/blob/078238a461a0872a8eacb887fbb3d0085714604c/staging/src/k8s.io/apiserver/pkg/apis/example/v1/types.go#L134
// Deprecated, for newer ones use '{"k1":"v1", "k2", "v2"}' form
// This expects the form k1=v1,k2=v2
func ParseDaemonNodeSelector(labels string) map[string]string {
	selectorMap := make(map[string]string)
	for _, label := range strings.Split(labels, ",") {
		label = strings.TrimSpace(label)
		if len(label) > 0 && strings.Contains(label, "=") {
			data := strings.SplitN(label, "=", 2)
			if len(data) >= 2 {
				if len(data[0]) > 0 && len(data[1]) > 0 {
					selectorMap[data[0]] = data[1]
				}
			}
		}
	}
	return selectorMap
}
