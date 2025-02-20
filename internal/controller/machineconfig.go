package controller

import (
	"encoding/json"

	machineconfigv1 "github.com/openshift/api/machineconfiguration/v1"
	ctrlcommon "github.com/openshift/machine-config-operator/pkg/controller/common"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// apiVersion: machineconfiguration.openshift.io/v1
// kind: MachineConfig
// metadata:
//   labels:
//     machineconfiguration.openshift.io/role: worker
//   name: 00-worker-ibm-spectrum-scale-kernel-devel
// spec:
//   config:
//     ignition:
//       version: 3.2.0
//   extensions:
//   - kernel-devel

func NewMachineConfig(labels map[string]string) *machineconfigv1.MachineConfig {
	return &machineconfigv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: machineconfigv1.SchemeGroupVersion.String(),
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "00-worker-ibm-spectrum-scale-kernel-devel",
			Labels: labels,
		},
		Spec: *NewMachineConfigSpec(),
	}
}

func NewMachineConfigSpec() *machineconfigv1.MachineConfigSpec {
	tmpIgnCfg := ctrlcommon.NewIgnConfig()
	rawTmpIgnCfg, _ := json.Marshal(tmpIgnCfg)

	return &machineconfigv1.MachineConfigSpec{
		// config is a Ignition Config object.
		// +optional
		Config: runtime.RawExtension{
			Raw: rawTmpIgnCfg,
		},
		// extensions contains a list of additional features that can be enabled on host
		// +listType=atomic
		// +optional
		Extensions: []string{"kernel-devel"},
	}
}
