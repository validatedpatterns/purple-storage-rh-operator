module github.com/validatedpatterns/purple-storage-rh-operator

go 1.23.0

toolchain go1.23.6

require (
	github.com/Masterminds/semver/v3 v3.3.1
	github.com/go-logr/logr v1.4.2
	github.com/google/uuid v1.6.0
	github.com/manifestival/controller-runtime-client v0.4.0
	github.com/manifestival/manifestival v0.7.2
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onsi/ginkgo/v2 v2.23.0
	github.com/onsi/gomega v1.36.2
	github.com/openshift/api v0.0.0-20240821141839-566edf361acd
	github.com/openshift/client-go v0.0.0-20240528061634-b054aa794d87
	github.com/openshift/library-go v0.0.0-20240711100342-737dc0fa5232
	github.com/openshift/machine-config-operator v0.0.0-00010101000000-000000000000
	github.com/pborman/uuid v1.2.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.80.1
	github.com/prometheus/client_golang v1.21.1
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.9.1
	github.com/stretchr/testify v1.10.0
	go.uber.org/zap v1.27.0
	golang.org/x/mod v0.24.0
	golang.org/x/net v0.37.0
	golang.org/x/sys v0.31.0
	k8s.io/api v0.32.2
	k8s.io/apiextensions-apiserver v0.32.2
	k8s.io/apimachinery v0.32.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/component-helpers v0.32.2
	k8s.io/klog/v2 v2.130.1
	sigs.k8s.io/controller-runtime v0.20.3
	sigs.k8s.io/sig-storage-local-static-provisioner v0.0.0-20250101095114-15a6cc21931d
)

require (
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/ajeddeloh/go-json v0.0.0-20170920214419-6a2fe990e083 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clarketm/json v1.14.1 // indirect
	github.com/coreos/fcct v0.5.0 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/coreos/ign-converter v0.0.0-20200629171308-e40a44f244c5 // indirect
	github.com/coreos/ignition v0.35.0 // indirect
	github.com/coreos/ignition/v2 v2.3.0 // indirect
	github.com/coreos/vcontext v0.0.0-20191017033345-260217907eb5 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/emicklei/go-restful/v3 v3.12.1 // indirect
	github.com/evanphx/json-patch v5.8.1+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.9.11 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-logr/zapr v1.3.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.2.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/pprof v0.0.0-20241210010833-40e02aabc2ad // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/kubernetes-csi/csi-proxy/client v1.0.2 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/miekg/dns v1.1.61 // indirect
	github.com/moby/sys/mountinfo v0.7.2 // indirect
	github.com/moby/sys/userns v0.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/vincent-petithory/dataurl v0.0.0-20160330182126-9a301d65acbb // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.opentelemetry.io/otel v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go4.org v0.0.0-20200104003542-c7e774b10ea0 // indirect
	golang.org/x/oauth2 v0.24.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/time v0.7.0 // indirect
	golang.org/x/tools v0.30.0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.4.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240826202546-f6391c0de4c7 // indirect
	google.golang.org/grpc v1.65.0 // indirect
	google.golang.org/protobuf v1.36.1 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.12.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiserver v0.32.2 // indirect
	k8s.io/component-base v0.32.2 // indirect
	k8s.io/controller-manager v0.28.12 // indirect
	k8s.io/kube-openapi v0.0.0-20241105132330-32ad38e42d3f // indirect
	k8s.io/kubernetes v1.30.2 // indirect
	k8s.io/mount-utils v0.30.2 // indirect
	k8s.io/utils v0.0.0-20241210054802-24370beab758 // indirect
	sigs.k8s.io/json v0.0.0-20241014173422-cfa47c3a1cc8 // indirect
	sigs.k8s.io/kube-storage-version-migrator v0.0.6-0.20230721195810-5c8923c5ff96 // indirect
	sigs.k8s.io/sig-storage-lib-external-provisioner/v6 v6.3.0 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.5.0 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

replace (
	// Usually not needed unless we have newer k8s.io deps and argocd still uses an old gitops-engine
	// github.com/argoproj/gitops-engine => github.com/argoproj/gitops-engine v0.0.0-20240905010810-bd7681ae3f8b
	// github.com/openshift/api => github.com/openshift/api v0.0.0-20240124164020-e2ce40831f2e
	// github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20240115204758-e6bf7d631d5e
	github.com/openshift/machine-config-operator => github.com/openshift/machine-config-operator v0.0.1-0.20200913004441-7eba765c69c9
	// Caused by Argo importing 'k8s.io/api'
	// Override all the v0.0.0 entries
	k8s.io/api => k8s.io/api v0.32.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.32.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.32.2
	k8s.io/apiserver => k8s.io/apiserver v0.32.2
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.32.2
	k8s.io/client-go => k8s.io/client-go v0.32.2
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.32.2
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.32.2
	k8s.io/code-generator => k8s.io/code-generator v0.32.2
	k8s.io/component-base => k8s.io/component-base v0.32.2
	k8s.io/component-helpers => k8s.io/component-helpers v0.32.2
	k8s.io/controller-manager => k8s.io/controller-manager v0.32.2
	k8s.io/cri-api => k8s.io/cri-api v0.32.2
	k8s.io/cri-client => k8s.io/cri-client v0.32.2
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.32.2
	k8s.io/dynamic-resource-allocation => k8s.io/dynamic-resource-allocation v0.32.2
	k8s.io/endpointslice => k8s.io/endpointslice v0.32.2
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.32.2
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.32.2
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.32.2
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.32.2
	k8s.io/kubectl => k8s.io/kubectl v0.32.2
	k8s.io/kubelet => k8s.io/kubelet v0.32.2
	k8s.io/kubernetes => k8s.io/kubernetes v1.32.2
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.32.2
	k8s.io/metrics => k8s.io/metrics v0.32.2
	k8s.io/mount-utils => k8s.io/mount-utils v0.32.2
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.32.2
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.32.2
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.20.3
	sigs.k8s.io/sig-storage-local-static-provisioner => sigs.k8s.io/sig-storage-local-static-provisioner v0.0.0-20250101095114-15a6cc21931d
)
