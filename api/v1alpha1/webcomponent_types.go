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
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WebComponentSpec defines the desired state of WebComponent
type WebComponentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Reference to a microfrontend from which the webcomponent would be served.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	MicroFrontend string `json:"microFrontend,omitempty"`

	// The HTML element tag name to be used when the matcher is matched.
	// +kubebuilder:example="my-menu-item"
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Element string `json:"element"`

	// This is a list of key-value pairs that allows you to assign specific attributes to the element. The name field is used as the attribute name, while the value field can be any valid JSON type.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Attributes []Attribute `json:"attributes,omitempty"`

	// This is a list of matchers that allows you to define the conditions under which the web component should be loaded.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Matchers []Matcher `json:"matchers,omitempty"`
}

// Attribute defines a key-value pair that allows you to assign specific attributes to the element. The name field is used as the attribute name, while the value field can be any valid JSON type.
type Attribute struct {

	// The name of the attribute.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Name string `json:"name"`

	// The value of the attribute.
	// +kubebuilder:validation:XPreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Value runtime.RawExtension `json:"value"`
}

// Matcher defines the conditions under which the web component should be loaded.
type Matcher struct {
	// This is a list of context names in which this element is intended to be shown.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ContextNames []string `json:"context-names,omitempty"`

	// The list of paths in which this element is intended to be shown.
	// +kubebuilder:example="/my-menu-item"
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Paths []string `json:"paths,omitempty"`

	// The list of roles for which this element is intended to be shown.
	// +kubebuilder:example="admin"
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Roles []string `json:"roles,omitempty"`
}

// WebComponentStatus defines the observed state of WebComponent
type WebComponentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// WebComponent is the Schema for the webcomponents API
type WebComponent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebComponentSpec   `json:"spec,omitempty"`
	Status WebComponentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// WebComponentList contains a list of WebComponent
type WebComponentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WebComponent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WebComponent{}, &WebComponentList{})
}
