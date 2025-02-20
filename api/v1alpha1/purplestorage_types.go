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

type MachineConfig struct {
	Labels map[string]string `json:"labels,omitempty"`
}
type IBMSpectrumCluster struct {
	// Boolean to create the CNSA cluster object
	Create bool `json:"create,omitempty"`
}

// PurpleStorageSpec defines the desired state of PurpleStorage
type PurpleStorageSpec struct {
	// Version of IBMs installation manifests found at https://github.com/IBM/ibm-spectrum-scale-container-native
	IbmCnsaVersion string        `json:"ibm_cnsa_version,omitempty"`
	MachineConfig  MachineConfig `json:"machineconfig,omitempty"`
	// PullSecret is the secret that contains the credentials to pull the images from the Container Registry
	PullSecret string             `json:"pull_secret,omitempty"`
	Cluster    IBMSpectrumCluster `json:"ibm_cnsa_cluster,omitempty"`
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
