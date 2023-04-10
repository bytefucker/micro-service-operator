/*
Copyright 2023.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ServicesGroupSpec defines the desired state of ServicesGroup
type ServicesGroupSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ServicesGroup. Edit servicesgroup_types.go to remove/update
	Services []Service `json:"services,omitempty"`
}

type Service struct {
	Name          string `json:"name,omitempty"`  //服务名称
	Image         string `json:"image,omitempty"` //服务镜像
	ContainerPort int32  `json:"containerPort,omitempty"`
	Rank          int    `json:"rank,omitempty"`      //启动顺序
	Replicas      *int32 `json:"replicas,omitempty" ` //分片数量
}

// ServicesGroupStatus defines the observed state of ServicesGroup
type ServicesGroupStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ServicesGroup is the Schema for the servicesgroups API
type ServicesGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServicesGroupSpec   `json:"spec,omitempty"`
	Status ServicesGroupStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ServicesGroupList contains a list of ServicesGroup
type ServicesGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServicesGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServicesGroup{}, &ServicesGroupList{})
}
