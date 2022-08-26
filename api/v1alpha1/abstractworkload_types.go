/*
Copyright 2022.

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
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AbstractWorkloadSpec defines the desired state of AbstractWorkload
type AbstractWorkloadSpec struct {
	// +required
	Replicas int `json:"replicas"`
	// +required
	ContainerImage string `json:"containerImage"`
	// +required
	WorkloadType WorkloadType `json:"workloadType"`
}

// AbstractWorkloadStatus defines the observed state of AbstractWorkload
type AbstractWorkloadStatus struct {
	Workload CrossNamespaceObjectReference `json:"workload"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AbstractWorkload is the Schema for the abstractworkloads API
// +kubebuilder:resource:shortName="aw"
// +kubebuilder:printcolumn:name="Replicas",type=string,JSONPath=`.spec.replicas`
// +kubebuilder:printcolumn:name="WorkloadType",type=string,JSONPath=`.spec.workloadType`
// +kubebuilder:printcolumn:name="WorkloadType",type=string,JSONPath=`.status.workload.kind`
type AbstractWorkload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AbstractWorkloadSpec   `json:"spec,omitempty"`
	Status AbstractWorkloadStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AbstractWorkloadList contains a list of AbstractWorkload
type AbstractWorkloadList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AbstractWorkload `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AbstractWorkload{}, &AbstractWorkloadList{})
}

type WorkloadType int64

const (
	Stateless WorkloadType = iota
	Stateful
)

func (w WorkloadType) String() string {
	switch w {
	case Stateless:
		return "stateless"
	case Stateful:
		return "stateful"
	}
	return "unknown"
}

type CrossNamespaceObjectReference struct {
	// API version of the referent.
	// +optional
	APIVersion string `json:"apiVersion,omitempty"`

	// Kind of the referent.
	// +required
	Kind string `json:"kind"`

	// Name of the referent.
	// +required
	Name string `json:"name"`

	// Namespace of the referent, defaults to the namespace of the Kubernetes resource object that contains the reference.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

func (s *CrossNamespaceObjectReference) String() string {
	if s.Namespace != "" {
		return fmt.Sprintf("%s/%s/%s", s.Kind, s.Namespace, s.Name)
	}
	return fmt.Sprintf("%s/%s", s.Kind, s.Name)
}
