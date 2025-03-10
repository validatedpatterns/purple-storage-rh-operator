package common

import (
	"os"
)

const (
	defaultDiskMakerImageVersion = "quay.io/hybridcloudpatterns/purple-storage-rh-operator-diskmaker"
	defaultKubeProxyImage        = "quay.io/openshift/origin-kube-rbac-proxy:latest"

	// OwnerNamespaceLabel references the owning object's namespace
	OwnerNamespaceLabel = "purple.purplestorage.com/owner-namespace"
	// OwnerNameLabel references the owning object
	OwnerNameLabel = "purple.purplestorage.com/owner-name"

	// DiskMakerImageEnv is used by the operator to read the DISKMAKER_IMAGE from the environment
	DiskMakerImageEnv = "DISKMAKER_IMAGE"
	// KubeRBACProxyImageEnv is used by the operator to read the KUBE_RBAC_PROXY_IMAGE from the environment
	KubeRBACProxyImageEnv = "KUBE_RBAC_PROXY_IMAGE"

	// DiscoveryNodeLabelKey is the label key on the discovery result CR used to identify the node it belongs to.
	// the value is the node's name
	DiscoveryNodeLabel = "discovery-result-node"

	DiskMakerDiscoveryDaemonSetTemplate = "templates/diskmaker-discovery-daemonset.yaml"
	MetricsServiceTemplate              = "templates/localmetrics/service.yaml"
	MetricsServiceMonitorTemplate       = "templates/localmetrics/service-monitor.yaml"
	PrometheusRuleTemplate              = "templates/localmetrics/prometheus-rule.yaml"

	// DiscoveryServiceName is the name of the service created for the diskmaker discovery daemon
	DiscoveryServiceName = "local-storage-discovery-metrics"

	// DiscoveryMetricsServingCert is the name of secret created for discovery service to store TLS config
	DiscoveryMetricsServingCert = "discovery-metric-serving-cert"
)

// GetDiskMakerImage returns the image to be used for diskmaker daemonset
func GetDiskMakerImage() string {
	if diskMakerImageFromEnv := os.Getenv(DiskMakerImageEnv); diskMakerImageFromEnv != "" {
		return diskMakerImageFromEnv
	}
	return defaultDiskMakerImageVersion
}

// GetKubeRBACProxyImage returns the image to be used for Kube RBAC Proxy sidecar container
func GetKubeRBACProxyImage() string {
	if kubeRBACProxyImageFromEnv := os.Getenv(KubeRBACProxyImageEnv); kubeRBACProxyImageFromEnv != "" {
		return kubeRBACProxyImageFromEnv
	}
	return defaultKubeProxyImage
}
