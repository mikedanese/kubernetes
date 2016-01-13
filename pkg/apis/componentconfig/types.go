/*
Copyright 2015 The Kubernetes Authors All rights reserved.

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

package componentconfig

import "k8s.io/kubernetes/pkg/api/unversioned"

type KubeProxyConfiguration struct {
	unversioned.TypeMeta

	// bindAddress is the IP address for the proxy server to serve on (set to 0.0.0.0 for all interfaces)
	BindAddress string `json:"bindAddress"`
	// cleanupIPTables
	CleanupIPTables bool `json:"cleanupIPTables"`
	// healthzBindAddress is the IP address for the health check server to serve on, defaulting to 127.0.0.1 (set to 0.0.0.0 for all interfaces)
	HealthzBindAddress string `json:"healthzBindAddress"`
	// healthzPort is the port to bind the health check server. Use 0 to disable.
	HealthzPort int `json:"healthzPort"`
	// hostnameOverride, if non-empty, will be used as the identity instead of the actual hostname.
	HostnameOverride string `json:"hostnameOverride"`
	// iptablesSyncPeriodSeconds is the period that iptables rules are refreshed (e.g. '5s', '1m', '2h22m').  Must be greater than 0.
	IPTablesSyncePeriodSeconds int `json:"iptablesSyncPeriodSeconds"`
	// kubeAPIBurst is the burst to use while talking with kubernetes apiserver
	KubeAPIBurst int `json:"kubeAPIBurst"`
	// kubeAPIQPS is the max QPS to use while talking with kubernetes apiserver
	KubeAPIQPS int `json:"kubeAPIQPS"`
	// kubeconfigPath is the path to the kubeconfig file with authorization information (the master location is set by the master flag).
	KubeconfigPath string `json:"kubeconfigPath"`
	// masqueradeAll tells kube-proxy to SNAT everything if using the pure iptables proxy mode.
	MasqueradeAll bool `json:"masqueradeAll"`
	// master is the address of the Kubernetes API server (overrides any value in kubeconfig)
	Master string `json:"master"`
	// oomScoreAdj is the oom-score-adj value for kube-proxy process. Values must be within the range [-1000, 1000]
	OOMScoreAdj *int `json:"oomScoreAdj"`
	// mode specifies which proxy mode to use.
	Mode ProxyMode `json:"mode"`
	// portRange is the range of host ports (beginPort-endPort, inclusive) that may be consumed in order to proxy service traffic. If unspecified (0-0) then ports will be randomly chosen.
	PortRange string `json:"portRange"`
	// resourceContainer is the bsolute name of the resource-only container to create and run the Kube-proxy in (Default: /kube-proxy).
	ResourceContainer string `json:"resourceContainer"`
	// udpTimeoutMilliseconds is how long an idle UDP connection will be kept open (e.g. '250ms', '2s').  Must be greater than 0. Only applicable for proxyMode=userspace.
	UDPTimeoutMilliseconds int `json:"udpTimeoutMilliseconds"`
}

// Currently two modes of proxying are available: 'userspace' (older, stable) or 'iptables' (experimental). If blank, look at the Node object on the Kubernetes API and respect the 'net.experimental.kubernetes.io/proxy-mode' annotation if provided.  Otherwise use the best-available proxy (currently userspace, but may change in future versions).  If the iptables proxy is selected, regardless of how, but the system's kernel or iptables versions are insufficient, this always falls back to the userspace proxy.
type ProxyMode string

const (
	ProxyModeUserspace ProxyMode = "userspace"
	ProxyModeIPTables  ProxyMode = "iptables"
)

type KubeAPIServerConfiguration struct {
	// insecurePort is the port on which to serve unsecured, unauthenticated
	// access. Default 8080. It is assumed that firewall rules are set up
	// such that this port is not reachable from outside of the cluster and
	// that port 443 on the cluster's public address is proxied to this port.
	// This is performed by nginx in the default setup.
	InsecurePort int `json:"insecurePort"`
	// insecureBindAddress is the IP address on which to serve the
	// insecurePort (set to 0.0.0.0 for all interfaces).
	InsecureBindAddress ip `json:"insecureBindAddress"`
	// bindAddress is the IP address on which to listen for the securePort.
	// The associated interface(s) must be reachable by the rest of the
	// cluster, and by CLI/web clients.
	BindAddress ip `json:"bindAddress"`
	// advertiseAddress is the IP address on which to advertise the apiserver
	// to members of the cluster. This address must be reachable by the rest
	// of the cluster. If blank, the --bind-address will be used.
	AdvertiseAddress ip `json:"advertiseAddress"`
	// securePort is the port on which to serve HTTPS with authentication and
	// authorization.
	SecurePort int `json:"securePort"`
	// tLSCertFile is file containing x509 Certificate for HTTPS.  (CA cert,
	// if any, concatenated after server cert). If HTTPS serving is enabled,
	// and tlsCertFile and tlsPrivateKeyFile are not provided, a self-signed
	// certificate and key are generated for the public address and saved to
	// /var/run/kubernetes.
	TLSCertFile string `json:"tLSCertFile"`
	// tLSPrivateKeyFile is file containing x509 private key matching tlsCertFile.
	TLSPrivateKeyFile string `json:"tLSPrivateKeyFile"`
	// certDirectory is the directory where the TLS certs are located (by
	// default /var/run/kubernetes). If --tls-cert-file and tlsPrivateKeyFile
	// are provided, this flag will be ignored.
	CertDirectory string `json:"certDirectory"`
	// storageVersions is the versions to store resources with. Different
	// groups may be stored in different versions. Specified in the format
	// "group1/version1,group2/version2...". This flag expects a complete
	// list of storage versions of ALL groups registered in the server. It
	// defaults to a list of preferred versions of all registered groups,
	// which is derived from the KUBE_API_VERSIONS environment variable.
	StorageVersions string `json:"storageVersions"`
	// cloudProvider is the provider for cloud services.  Empty string for
	// no provider.
	CloudProvider string `json:"cloudProvider"`
	// cloudConfigFile is the path to the cloud provider configuration file.
	// Empty string for no configuration file.
	CloudConfigFile string `json:"cloudConfigFile"`
	// eventTTL is Amount of time to retain events. Default 1 hour.
	EventTTL duration `json:"eventTTL"`
	// basicAuthFile is If set, the file that will be used to admit requests
	// to the secure port of the API server via http basic authentication.
	BasicAuthFile string `json:"basicAuthFile"`
	// clientCAFile is If set, any request presenting a client certificate
	// signed by one of the authorities in the client-ca-file is authenticated
	// with an identity corresponding to the CommonName of the client certificate.
	ClientCAFile string `json:"clientCAFile"`
	// tokenAuthFile is If set, the file that will be used to secure the
	// secure port of the API server via token authentication.
	TokenAuthFile string `json:"tokenAuthFile"`
	// oIDCIssuerURL is The URL of the OpenID issuer, only HTTPS scheme will
	// be accepted. If set, it will be used to verify the OIDC JSON Web
	// Token (JWT)
	OIDCIssuerURL string `json:"oIDCIssuerURL"`
	// oIDCClientID is The client ID for the OpenID Connect client, must be
	// set if oidc-issuer-url is set
	OIDCClientID string `json:"oIDCClientID"`
	// oIDCCAFile is If set, the OpenID server's certificate will be verified
	// by one of the authorities in the oidc-ca-file, otherwise the host's
	// root CA set will be used
	OIDCCAFile string `json:"oIDCCAFile"`
	// oIDCUsernameClaim is The OpenID claim to use as the user name. Note
	// that claims other than the default ('sub') is not guaranteed to be
	// unique and immutable. This flag is experimental, please see the
	// authentication documentation for further details.
	OIDCUsernameClaim string `json:"oIDCUsernameClaim"`
	// serviceAccountKeyFile is File containing PEM-encoded x509 RSA private
	// or public key, used to verify ServiceAccount tokens. If unspecified,
	// tlsPrivateKeyFile is used.
	ServiceAccountKeyFile string `json:"serviceAccountKeyFile"`
	// serviceAccountLookup is If true, validate ServiceAccount tokens exist
	// in etcd as part of authentication.
	ServiceAccountLookup bool `json:"serviceAccountLookup"`
	// keystoneURL is If passed, activates the keystone authentication plugin
	KeystoneURL string `json:"keystoneURL"`
	// authorizationMode is Ordered list of plug-ins to do authorization on
	// secure port. Comma-delimited list of AuthorizationModeChoices
	AuthorizationMode string `json:"authorizationMode"`
	// authorizationPolicyFile is File with authorization policy in csv format,
	// used with authorizationMode=ABAC, on the secure port.
	AuthorizationPolicyFile string `json:"authorizationPolicyFile"`
	// admissionControl is Ordered list of plug-ins to do admission control
	// of resources into cluster. Comma-delimited list of: AdmissionControls
	AdmissionControl string `json:"admissionControl"`
	// admissionControlConfigFile is File with admission control configuration.
	AdmissionControlConfigFile string `json:"admissionControlConfigFile"`
	// etcdServerList is List of etcd servers to watch (http://ip:port),
	// comma separated. Mutually exclusive with -etcd-config
	EtcdServerList stringslice `json:"etcdServerList"`
	// etcdServersOverrides is Per-resource etcd servers overrides, comma
	// separated. The individual override format: group/resource#servers,
	// where servers are http://ip:port, semicolon separated.
	EtcdServersOverrides stringslice `json:"etcdServersOverrides"`
	// etcdPathPrefix is The prefix for all resource paths in etcd.
	EtcdPathPrefix string `json:"etcdPathPrefix"`
	// corsAllowedOriginList is List of allowed origins for CORS, comma
	// separated.  An allowed origin can be a regular expression to support
	// subdomain matching.  If this list is empty CORS will not be enabled.
	CorsAllowedOriginList stringslice `json:"corsAllowedOriginList"`
	// allowPrivileged is If true, allow privileged containers.
	AllowPrivileged bool `json:"allowPrivileged"`
	// serviceClusterIPRange is A CIDR notation IP range from which to assign
	// service cluster IPs. This must not overlap with any IP ranges assigned
	// to nodes for pods.
	ServiceClusterIPRange ipnet `json:"serviceClusterIPRange"`
	// serviceNodePort is a port range to reserve for services with NodePort
	// visibility.  Example: '30000-32767'.  Inclusive at both ends of the range.
	ServiceNodePortRange map[string]string `json:"serviceNodePortRange`
	// masterServiceNamespace is The namespace from which the kubernetes master
	// services should be injected into pods
	MasterServiceNamespace string `json:"masterServiceNamespace"`
	// masterCount is The number of apiservers running in the cluster
	MasterCount int `json:"masterCount"`
	// runtimeConfig is a set of key=value pairs that describe runtime configuration
	// that may be passed to apiserver. apis/<groupVersion> key can be used to
	// turn on/off specific api versions. apis/<groupVersion>/<resource> can
	// be used to turn on/off specific resources. api/all and api/legacy are
	// special keys to control all and legacy api versions respectively.
	RuntimeConfig map[string]string `json:"runtimeConfig"`
	// enableProfiling is Enable profiling via web interface host:port/debug/pprof/
	EnableProfiling bool `json:"enableProfiling"`
	// enableWatchCache is Enable watch caching in the apiserver
	EnableWatchCache bool `json:"enableWatchCache"`
	// externalHost is The hostname to use when generating externalized URLs
	// for this master (e.g. Swagger API Docs.)
	ExternalHost string `json:"externalHost"`
	// maxRequestsInFlight is The maximum number of requests in flight at a
	// given time.  When the server exceeds this, it rejects requests.
	// Zero for no limit.
	MaxRequestsInFlight int `json:"maxRequestsInFlight"`
	// minRequestTimeout is An optional field indicating the minimum number
	// of seconds a handler must keep a request open before timing it out.
	// Currently only honored by the watch request handler, which picks a
	// randomized value above this number as the connection timeout, to
	// spread out load.
	MinRequestTimeout int `json:"minRequestTimeout"`
	// longRunningRequestRE is A regular expression matching long running
	// requests which should be excluded from maximum inflight request handling.
	LongRunningRequestRE string `json:"longRunningRequestRE"`
	// sshUser is If non-empty, use secure SSH proxy to the nodes, using this
	// user name
	SSHUser string `json:"sshUser"`
	// sshKeyfile is If non-empty, use secure SSH proxy to the nodes, using
	// this user keyfile
	SSHKeyfile string `json:"sshKeyfile"`
	// maxConnectionBytesPerSec is If non-zero, throttle each user connection
	// to this number of bytes/sec.  Currently only applies to long-running
	// requests
	MaxConnectionBytesPerSec int64 `json:"maxConnectionBytesPerSec"`

	// !!!!! Uh oh: 	fs.BoolVar(&s.KubeletConfig.EnableHttps, "kubelet-https", s.KubeletConfig.EnableHttps, "Use https for kubelet connections"), []
	// !!!!! Uh oh: 	fs.DurationVar(&s.KubeletConfig.HTTPTimeout, "kubelet-timeout", s.KubeletConfig.HTTPTimeout, "Timeout for kubelet operations"), []
	// !!!!! Uh oh: 	fs.StringVar(&s.KubeletConfig.CertFile, "kubelet-client-certificate", s.KubeletConfig.CertFile, "Path to a client cert file for TLS."), []
	// !!!!! Uh oh: 	fs.StringVar(&s.KubeletConfig.KeyFile, "kubelet-client-key", s.KubeletConfig.KeyFile, "Path to a client key file for TLS."), []
	// !!!!! Uh oh: 	fs.StringVar(&s.KubeletConfig.CAFile, "kubelet-certificate-authority", s.KubeletConfig.CAFile, "Path to a cert. file for the certificate authority."), []

	// kubernetesServiceNodePort is If non-zero, the Kubernetes master service
	// (which apiserver creates/maintains) will be of type NodePort, using
	// this as the value of the port. If zero, the Kubernetes master service
	// will be of type ClusterIP.
	KubernetesServiceNodePort int `json:"kubernetesServiceNodePort"`
	// !!!!! Uh oh: 	fs.BoolVar(&validation.RepairMalformedUpdates, "repair-malformed-updates", true, "If true, server will do its best to fix the update request to pass the validation, e.g., setting empty UID in update request to its existing value. This flag can be turned off after we fix all the clients that send malformed updates."), []
}
