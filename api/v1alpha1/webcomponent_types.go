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

// WebComponentSpec defines the desired state of WebComponent
type WebComponentSpec struct {
	// Reference to a microfrontend from which the webcomponent would be served.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	MicroFrontend *string `json:"microFrontend"`

	// The HTML element tag name to be used when the matcher is matched.
	// +kubebuilder:example="my-menu-item"
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Element *string `json:"element"`

	// This is a list of key-value pairs that allows you to assign specific attributes to the element. The name field is used as the attribute name, while the value field can be any valid JSON type.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Attributes []Attribute `json:"attributes,omitempty"`

	// DisplayRules defines the conditions under which the web component should be loaded. If not specified, the web component will always be loaded.
	// There is an or opperation between the elements of the DisplayRules list. If any of the DisplayRules is matched, the web component will be loaded.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	DisplayRules []DisplayRules `json:"displayRules,omitempty"`

	// Priority defines the priority of the webcomponent. Used for ordering the webcomponent within the shell. The higher the number, the higher the priority. The default priority is 0.
	// +kubebuilder:default=0
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Priority *int32 `json:"priority,omitempty"`

	// Styles defines the styles that should be applied to the webcomponent.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Style *string `json:"style,omitempty"`
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

// DisplayRules defines the conditions under which the web component should be loaded.
// There is an and opperation between AllOf, AnyOf and NoneOf lists.
type DisplayRules struct {
	// If all of the matchers in this list are matched, the web component will be loaded.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	AllOf []Matcher `json:"allOf,omitempty"`

	// If any of the matchers in this list are matched, the web component will be loaded.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	AnyOf []Matcher `json:"anyOf,omitempty"`

	// If none of the matchers in this list are matched, the web component will be loaded.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	NoneOf []Matcher `json:"noneOf,omitempty"`
}

// Matcher defines the conditions under which the web component should be loaded.
// +kubebuilder:validation:MaxProperties=1
type Matcher struct {
	// This is a list of context names in which this element is intended to be shown.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ContextName string `json:"context-names,omitempty"`

	// The list of paths in which this element is intended to be shown.
	// +kubebuilder:example="/my-menu-item"
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Path string `json:"paths,omitempty"`

	// The list of roles for which this element is intended to be shown.
	// +kubebuilder:example="admin"
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Role string `json:"roles,omitempty"`
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
