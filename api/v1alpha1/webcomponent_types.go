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
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// WebComponentSpec defines the desired state of WebComponent
type WebComponentSpec struct {
	// Reference to a microfrontend from which the webcomponent would be served.
	// +optional
	// +kubebuilder:validation:MaxLength=253
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	MicroFrontend *string `json:"microFrontend,omitempty"`

	// The HTML element tag name to be used when the matcher is matched.
	// +kubebuilder:example="my-menu-item"
	// +kubebuilder:validation:MaxLength=256
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Element *string `json:"element"`

	// This is a list of key-value pairs that allows you to assign specific attributes to the element. The name field is used as the attribute name, while the value field can be any valid JSON type.
	// +kubebuilder:validation:MaxItems=64
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Attributes []Attribute `json:"attributes,omitempty"`

	// DisplayRules defines the conditions under which the web component should be loaded.
	// There is an or opperation between the elements of the DisplayRules list. If any of the DisplayRules is matched, the web component will be loaded.
	// +kubebuilder:validation:MaxItems=32
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	DisplayRules []DisplayRules `json:"displayRules"`

	// Priority defines the priority of the webcomponent. Used for ordering the webcomponent within the shell. The higher the number, the higher the priority. The default priority is 0.
	// +kubebuilder:default=0
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Priority *int32 `json:"priority,omitempty"`

	// Styles defines the styles that should be applied to the webcomponent.
	// +kubebuilder:validation:MaxItems=64
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Style []Style `json:"style,omitempty"`
}

// Attribute defines a key-value pair that allows you to assign specific attributes to the element. The name field is used as the attribute name, while the value field can be any valid JSON type.
type Attribute struct {
	// The name of the attribute.
	// +kubebuilder:validation:MaxLength=256
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Name string `json:"name"`

	// The value of the attribute.
	// +kubebuilder:validation:XPreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Value runtime.RawExtension `json:"value"`
}

// Style defines the styles that should be applied to the webcomponent.
type Style struct {
	// The name of the style.
	// +kubebuilder:validation:MaxLength=256
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Name string `json:"name"`

	// The value of the style.
	// +kubebuilder:validation:MaxLength=4096
	Value string `json:"value"`
}

// DisplayRules defines the conditions under which the web component should be loaded.
// There is an and opperation between AllOf, AnyOf and NoneOf lists.
type DisplayRules struct {
	// If all of the matchers in this list are matched, the web component will be loaded.
	// +kubebuilder:validation:MaxItems=16
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	AllOf []Matcher `json:"allOf,omitempty"`

	// If any of the matchers in this list are matched, the web component will be loaded.
	// +kubebuilder:validation:MaxItems=16
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	AnyOf []Matcher `json:"anyOf,omitempty"`

	// If none of the matchers in this list are matched, the web component will be loaded.
	// +kubebuilder:validation:MaxItems=16
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	NoneOf []Matcher `json:"noneOf,omitempty"`
}

// Matcher defines the conditions under which the web component should be loaded.
// A Matcher may contain scalar fields (context-name, path, role) and/or nested
// operator fields (allOf, anyOf, noneOf). All present conditions are combined with
// AND semantics: every condition must hold for the matcher to evaluate to true.
type Matcher struct {
	// This is a list of context names in which this element is intended to be shown.
	// +kubebuilder:validation:MaxLength=256
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ContextName string `json:"context-name,omitempty"`

	// The list of paths in which this element is intended to be shown.
	// +kubebuilder:example="/my-menu-item"
	// +kubebuilder:validation:MaxLength=2048
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Path string `json:"path,omitempty"`

	// The list of roles for which this element is intended to be shown.
	// +kubebuilder:example="admin"
	// +kubebuilder:validation:MaxLength=256
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Role string `json:"role,omitempty"`

	// Nested allOf: all sub-matchers must match.
	//
	// Note: +kubebuilder:pruning:PreserveUnknownFields and +kubebuilder:validation:Schemaless
	// are required here because Matcher references itself. controller-gen detects the cycle
	// and refuses to expand the schema, so Schemaless is the only escape hatch. As a
	// side-effect, +kubebuilder:validation:MaxItems=16 cannot be set — controller-gen rejects
	// array-specific markers on Schemaless fields. The semantically intended limit is 16,
	// matching the top-level DisplayRules fields, but it cannot be enforced at the CRD level.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	AllOf []Matcher `json:"allOf,omitempty"`

	// Nested anyOf: at least one sub-matcher must match.
	// See AllOf for an explanation of why MaxItems cannot be set here.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	AnyOf []Matcher `json:"anyOf,omitempty"`

	// Nested noneOf: none of the sub-matchers may match.
	// See AllOf for an explanation of why MaxItems cannot be set here.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	NoneOf []Matcher `json:"noneOf,omitempty"`
}

// ObjectReference contains information about a referenced object
type ObjectReference struct {
	// Name of the referenced object
	// +kubebuilder:validation:MaxLength=253
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Name string `json:"name"`

	// Namespace of the referenced object
	// +kubebuilder:validation:MaxLength=63
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Namespace string `json:"namespace"`

	// Found indicates if the referenced object was found
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Found bool `json:"found"`
}

// WebComponentStatus defines the observed state of WebComponent
type WebComponentStatus struct {
	// Conditions represent the latest available observations of the WebComponent's state
	// +optional
	// +listType=map
	// +listMapKey=type
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Phase represents the current lifecycle phase of the WebComponent
	// Possible values: Pending, Ready, Failed, MicroFrontendNotFound
	// +optional
	// +kubebuilder:validation:Enum=Pending;Ready;Failed;MicroFrontendNotFound
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Phase string `json:"phase,omitempty"`

	// MicroFrontendRef indicates the resolved MicroFrontend reference
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=status
	MicroFrontendRef *ObjectReference `json:"microFrontendRef,omitempty"`

	// ObservedGeneration reflects the generation of the most recently observed WebComponent
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=status
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// WebComponent is the Schema for the webcomponents API
type WebComponent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebComponentSpec   `json:"spec,omitempty"`
	Status WebComponentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// WebComponentList contains a list of WebComponent
type WebComponentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WebComponent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WebComponent{}, &WebComponentList{})
}
