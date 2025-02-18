package controller

import (
	machineconfigv1 "github.com/openshift/api/machineconfiguration/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func NewMachineConfig(label string) *machineconfigv1.MachineConfig {
	return &machineconfigv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: machineconfigv1.SchemeGroupVersion.String(),
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "00-worker-ibm-spectrum-scale-kernel-devel",
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": "worker",
			},
		},
		Spec: *NewMachineConfigSpec(),
	}
}

func NewMachineConfigSpec() *machineconfigv1.MachineConfigSpec {
	return &machineconfigv1.MachineConfigSpec{
		// config is a Ignition Config object.
		// +optional

		// extensions contains a list of additional features that can be enabled on host
		// +listType=atomic
		// +optional
		Extensions: []string{"kernel-devel"},
	}
}
