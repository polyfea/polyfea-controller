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

// MicroFrontendSpec defines the desired state of MicroFrontend
type MicroFrontendSpec struct {
	// Reference to a service from which the modules or css would be served.
	// Fully qualified name of the service should be specified in the format <schema>://<service-name>.<namespace>.<cluster>.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Service *string `json:"service"`

	// This specifies whether the loading of web components should be proxied by the controller.
	// +kubebuilder:default=true
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Proxy *bool `json:"proxy,omitempty"`

	// TODO: Make this work
	// CachingStrategy defines the caching strategy for the micro frontend.
	// +kubebuilder:default=none
	// +kubebuilder:validation:Enum=none;cache;
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	CacheStrategy string `json:"cacheStrategy,omitempty"`

	// TODO: Make this work
	// CacheControl defines the cache control header for the micro frontend. This is only used if the caching strategy is set to 'cache'.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	CacheControl *string `json:"cacheControl,omitempty"`

	// Relative path to the module file within the service.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ModulePath *string `json:"modulePath"`

	// Relative path to the static files within the service.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	StaticResources []StaticResources `json:"staticPaths,omitempty"`

	// FrontendClass is the name of the frontend class that should be used for this micro frontend.
	// +kubebuilder:default=polyfea-controller-default
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	FrontendClass *string `json:"frontendClass"`

	// List of dependencies that should be loaded before this micro frontend.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	DependsOn []string `json:"dependsOn,omitempty"`

	// TODO: Test similar to MFC values
	// CacheOptions specifies the cache settings for the PWA, including pre-caching and runtime caching.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	CacheOptions PWACache `json:"cacheOptions"`
}

// StaticResources defines the static resources that should be loaded before this micro frontend.
type StaticResources struct {
	// Kind defines the kind of the static resource can be script, stylesheet, or any other `link` element.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Kind string `json:"kind"`

	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Path string `json:"path"`

	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Attributes []Attribute `json:"attributes,omitempty"`

	// WaitOnLoad defines whether the micro frontend should wait for the static resource to load before loading itself.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	WaitOnLoad bool `json:"waitOnLoad,omitempty"`

	// This specifies whether the loading of static resource components should be proxied by the controller.
	// +kubebuilder:default=true
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Proxy *bool `json:"proxy,omitempty"`
}

// Port is the service port being referenced.
// +kubebuilder:validation:MaxProperties=1
type Port struct {
	// Name is the name of the port on the Service. This is a mutually exclusive setting with "Number".
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Name string `json:"name,omitempty"`
	// Number is the numerical port number (e.g. 80) on the Service. This is a mutually exclusive setting with "Name".
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Number *int32 `json:"number,omitempty"`
}

// MicroFrontendStatus defines the observed state of MicroFrontend
type MicroFrontendStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MicroFrontend is the Schema for the microfrontends API
type MicroFrontend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MicroFrontendSpec   `json:"spec,omitempty"`
	Status MicroFrontendStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MicroFrontendList contains a list of MicroFrontend
type MicroFrontendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MicroFrontend `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MicroFrontend{}, &MicroFrontendList{})
}
