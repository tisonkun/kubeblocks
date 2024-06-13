/*
Copyright (C) 2022-2024 ApeCloud Co., Ltd

This file is part of KubeBlocks project

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ClusterDefinitionSpec defines the desired state of ClusterDefinition
type ClusterDefinitionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ClusterDefinition. Edit clusterdefinition_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// ClusterDefinitionStatus defines the observed state of ClusterDefinition
type ClusterDefinitionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +genclient:nonNamespaced
// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:resource:categories={kubeblocks},scope=Cluster,shortName=cd
// +kubebuilder:printcolumn:name="Topologies",type="string",JSONPath=".status.topologies",description="topologies"
// +kubebuilder:printcolumn:name="ServiceRefs",type="string",JSONPath=".status.serviceRefs",description="service references"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.phase",description="status phase"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

// ClusterDefinition is the Schema for the clusterdefinitions API
type ClusterDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterDefinitionSpec   `json:"spec,omitempty"`
	Status ClusterDefinitionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterDefinitionList contains a list of ClusterDefinition
type ClusterDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterDefinition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterDefinition{}, &ClusterDefinitionList{})
}
