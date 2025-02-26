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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PurpleStorageSpec defines the desired state of PurpleStorage
type PurpleStorageSpec struct {
	// Version of IBMs installation manifests found at https://github.com/IBM/ibm-spectrum-scale-container-native
	// +operator-sdk:csv:customresourcedefinitions:type=spec,order=1
	IbmCnsaVersion string `json:"ibm_cnsa_version,omitempty"`
	// MachineConfig labelling for the installation of kernel-devel package
	// +operator-sdk:csv:customresourcedefinitions:type=spec,order=2
	MachineConfig MachineConfigLabels `json:"mcoconfig,omitempty"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec,order=3
	Cluster IBMSpectrumCluster `json:"ibm_cnsa_cluster,omitempty"`
}

type IBMSpectrumCluster struct {
	// Boolean to create the CNSA cluster object
	// +operator-sdk:csv:customresourcedefinitions:type=spec,order=4,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	// +kubebuilder:default:=true
	Create bool `json:"create,omitempty"`
	// Nodes with this label will be part of the cluster, must have at least 3 nodes with this
	// +operator-sdk:csv:customresourcedefinitions:type=spec,order=5,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldDependency:ibm_cnsa_cluster.create:true"}
	Daemon_nodeSelector map[string]string `json:"daemon_nodeSelector,omitempty"`
}

type MachineConfigLabels struct {
	// Labels to be used for the machineconfigpool
	// +operator-sdk:csv:customresourcedefinitions:type=spec,order=6
	McoLabels map[string]string `json:"mco_labels,omitempty"`
}

// PurpleStorageStatus defines the observed state of PurpleStorage
type PurpleStorageStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
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
