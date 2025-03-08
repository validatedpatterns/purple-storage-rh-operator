/*
Copyright 2025.

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

package v1alpha1

import (
	operatorv1 "github.com/openshift/api/operator/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PurpleStorageSpec defines the desired state of PurpleStorage
type PurpleStorageSpec struct {
	// MachineConfig labeling for the installation of kernel-devel package
	// +operator-sdk:csv:customresourcedefinitions:type=spec,order=1
	MachineConfig MachineConfig `json:"mco_config,omitempty"`

	// Version of IBMs installation manifests found at https://github.com/IBM/ibm-spectrum-scale-container-native
	// +operator-sdk:csv:customresourcedefinitions:type=spec,order=2
	IbmCnsaVersion string `json:"ibm_cnsa_version,omitempty"`

	// +operator-sdk:csv:customresourcedefinitions:type=spec,order=3
	Cluster IBMSpectrumCluster `json:"ibm_cnsa_cluster,omitempty"`

	// Inherited from LVSet to provide control over node selector and device filtering capabilities
	// +operator-sdk:csv:customresourcedefinitions:type=spec,order=4
	NodeSpec NodeSpec `json:"node_spec,omitempty"`
}

type NodeSpec struct {
	// Nodes on which the automatic detection policies must run.
	// +optional
	Selector *corev1.NodeSelector `json:"selector,omitempty"`

	// If specified, a list of tolerations to pass to the discovery daemons.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// DeviceInclusionSpec is the filtration rule for including a device in the device discovery
	// +optional
	DeviceInclusionSpec *DeviceInclusionSpec `json:"deviceInclusionSpec,omitempty"`
}

type MachineConfig struct {
	// Boolean to create the MachinConfig objects
	// +operator-sdk:csv:customresourcedefinitions:type=spec,order=4,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	// +kubebuilder:default:=true
	Create bool `json:"create,omitempty"`
	// Labels to be used for the machineconfigpool
	// +operator-sdk:csv:customresourcedefinitions:type=spec,order=5,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldDependency:mco_config.create:true"}
	Labels map[string]string `json:"labels,omitempty"`
}

type IBMSpectrumCluster struct {
	// Boolean to create the CNSA cluster object
	// +operator-sdk:csv:customresourcedefinitions:type=spec,order=6,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	// +kubebuilder:default:=true
	Create bool `json:"create,omitempty"`
	// Nodes with this label will be part of the cluster, must have at least 3 nodes with this
	// +operator-sdk:csv:customresourcedefinitions:type=spec,order=7,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldDependency:ibm_cnsa_cluster.create:true"}
	Daemon_nodeSelector map[string]string `json:"daemon_nodeSelector,omitempty"`
}

// PurpleStorageStatus defines the observed state of PurpleStorage
type PurpleStorageStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Conditions is a list of conditions and their status.
	Conditions []operatorv1.OperatorCondition `json:"conditions,omitempty"`
	// TotalProvisionedDeviceCount is the count of the total devices over which the PVs has been provisioned
	TotalProvisionedDeviceCount *int32 `json:"totalProvisionedDeviceCount,omitempty"`
	// observedGeneration is the last generation change the operator has dealt with
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PurpleStorage is the Schema for the purplestorages API
type PurpleStorage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PurpleStorageSpec   `json:"spec,omitempty"`
	Status PurpleStorageStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PurpleStorageList contains a list of PurpleStorage
type PurpleStorageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PurpleStorage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PurpleStorage{}, &PurpleStorageList{})
}

// DeviceMechanicalProperty holds the device's mechanical spec. It can be rotational or nonRotational
type DeviceMechanicalProperty string

// The mechanical properties of the devices
const (
	// Rotational refers to magnetic disks
	Rotational DeviceMechanicalProperty = "Rotational"
	// NonRotational refers to ssds
	NonRotational DeviceMechanicalProperty = "NonRotational"
)

// DeviceType is the types that will be supported by the LSO.
type DeviceType string

const (
	// RawDisk represents a device-type of block disk
	RawDisk DeviceType = "disk"
	// Partition represents a device-type of partition
	Partition DeviceType = "part"
	// Loop type device
	Loop DeviceType = "loop"
	// Multipath device type
	MultiPath DeviceType = "mpath"
)

// DeviceInclusionSpec holds the inclusion filter spec
type DeviceInclusionSpec struct {
	// Devices is the list of devices that should be used for automatic detection.
	// This would be one of the types supported by the local-storage operator. Currently,
	// the supported types are: disk, part. If the list is empty only `disk` types will be selected
	// +optional
	DeviceTypes []DeviceType `json:"deviceTypes,omitempty"`
	// DeviceMechanicalProperty denotes whether Rotational or NonRotational disks should be used.
	// by default, it selects both
	// +optional
	DeviceMechanicalProperties []DeviceMechanicalProperty `json:"deviceMechanicalProperties,omitempty"`
	// MinSize is the minimum size of the device which needs to be included. Defaults to `1Gi` if empty
	// +optional
	MinSize *resource.Quantity `json:"minSize,omitempty"`
	// MaxSize is the maximum size of the device which needs to be included
	// +optional
	MaxSize *resource.Quantity `json:"maxSize,omitempty"`
	// Models is a list of device models. If not empty, the device's model as outputted by lsblk needs
	// to contain at least one of these strings.
	// +optional
	Models []string `json:"models,omitempty"`
	// Vendors is a list of device vendors. If not empty, the device's model as outputted by lsblk needs
	// to contain at least one of these strings.
	// +optional
	Vendors []string `json:"vendors,omitempty"`
}
