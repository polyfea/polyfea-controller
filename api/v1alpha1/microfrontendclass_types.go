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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MicroFrontendClassSpec defines the desired state of MicroFrontendClass
type MicroFrontendClassSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// BaseUri for which the frontend class will be used
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	BaseUri string `json:"baseUri"`

	// CspHeader that will be used for the frontend class, none if not set
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	CspHeader string `json:"cspHeader,omitempty"`

	// ExtraHeaders that will be used for the frontend class, none if not set
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ExtraHeaders []string `json:"extraHeaders,omitempty"`
}

// MicroFrontendClassStatus defines the observed state of MicroFrontendClass
type MicroFrontendClassStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MicroFrontendClass is the Schema for the microfrontendclasses API
type MicroFrontendClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MicroFrontendClassSpec   `json:"spec,omitempty"`
	Status MicroFrontendClassStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MicroFrontendClassList contains a list of MicroFrontendClass
type MicroFrontendClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MicroFrontendClass `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MicroFrontendClass{}, &MicroFrontendClassList{})
}
