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
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Service *ServiceReference `json:"service"`

	// This specifies whether the loading of web components should be proxied by the controller. This is useful if the web component is served from within the cluster and cannot be accessed from outside the cluster network. The module will be served from the URL base_controller_url/web-components/web_component_name.jsm. This is the recommended approach for the standard assumed use-case.
	// +kubebuilder:default=true
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Proxy *bool `json:"proxy,omitempty"`

	// CachingStrategy defines the caching strategy for the micro frontend.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	// +kubebuilder:validation:Enum=none;
	CacheStrategy string `json:"cacheStrategy,omitempty"`

	// Relative path to the module file within the service.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ModulePath string `json:"modulePath,omitempty"`

	// Relative path to the static files within the service.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	StaticPath string `json:"staticPath,omitempty"`

	// The modules are not preloaded by default but only when navigating to some of the subpaths mentioned in the 'navigation' list. Setting this property to true ensures that the module is loaded when the application starts.
	// +kubebuilder:default=false
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Preload *bool `json:"preload,omitempty"`

	// FrontendClass is the name of the frontend class that should be used for this micro frontend.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	FrontendClass *string `json:"frontendClass,omitempty"`

	// List of dependencies that should be loaded before this micro frontend.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	DependsOn []string `json:"dependsOn,omitempty"`
}

// ServiceReference references a Kubernetes Service as a Backend.
type ServiceReference struct {
	// Name is the name of the service being referenced.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Name *string `json:"name,omitempty"`
	// Port is the port of the service being referenced.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Port *Port `json:"port,omitempty"`
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
