// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	metav1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

// ContainerRuntimeConfigSpecApplyConfiguration represents a declarative configuration of the ContainerRuntimeConfigSpec type for use
// with apply.
type ContainerRuntimeConfigSpecApplyConfiguration struct {
	MachineConfigPoolSelector *metav1.LabelSelectorApplyConfiguration          `json:"machineConfigPoolSelector,omitempty"`
	ContainerRuntimeConfig    *ContainerRuntimeConfigurationApplyConfiguration `json:"containerRuntimeConfig,omitempty"`
}

// ContainerRuntimeConfigSpecApplyConfiguration constructs a declarative configuration of the ContainerRuntimeConfigSpec type for use with
// apply.
func ContainerRuntimeConfigSpec() *ContainerRuntimeConfigSpecApplyConfiguration {
	return &ContainerRuntimeConfigSpecApplyConfiguration{}
}

// WithMachineConfigPoolSelector sets the MachineConfigPoolSelector field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MachineConfigPoolSelector field is set to the value of the last call.
func (b *ContainerRuntimeConfigSpecApplyConfiguration) WithMachineConfigPoolSelector(value *metav1.LabelSelectorApplyConfiguration) *ContainerRuntimeConfigSpecApplyConfiguration {
	b.MachineConfigPoolSelector = value
	return b
}

// WithContainerRuntimeConfig sets the ContainerRuntimeConfig field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ContainerRuntimeConfig field is set to the value of the last call.
func (b *ContainerRuntimeConfigSpecApplyConfiguration) WithContainerRuntimeConfig(value *ContainerRuntimeConfigurationApplyConfiguration) *ContainerRuntimeConfigSpecApplyConfiguration {
	b.ContainerRuntimeConfig = value
	return b
}
