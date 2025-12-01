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
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceReference defines how to reach the service hosting the micro frontend
// +kubebuilder:validation:XValidation:rule="(has(self.name) && size(self.name) > 0) != (has(self.uri) && size(self.uri) > 0)",message="Either 'name' or 'uri' must be specified, but not both"
type ServiceReference struct {
	// Name of the Kubernetes service (mutually exclusive with URI)
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Name *string `json:"name,omitempty"`

	// URI for external services (mutually exclusive with Name)
	// Should include schema (http:// or https://)
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	URI *string `json:"uri,omitempty"`

	// Namespace of the service. Defaults to the MicroFrontend's namespace if not specified.
	// Only used when Name is set.
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Namespace *string `json:"namespace,omitempty"`

	// Port of the service. Defaults to 80 if not specified.
	// Only used when Name is set.
	// +optional
	// +kubebuilder:default=80
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Port *int32 `json:"port,omitempty"`

	// Scheme to use for connection (http or https). Defaults to http.
	// Only used when Name is set.
	// +optional
	// +kubebuilder:default=http
	// +kubebuilder:validation:Enum=http;https
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Scheme *string `json:"scheme,omitempty"`

	// Domain is the cluster domain suffix. Defaults to svc.cluster.local if not specified.
	// Only used when Name is set. Allows customization for different cluster implementations.
	// +optional
	// +kubebuilder:default=svc.cluster.local
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Domain *string `json:"domain,omitempty"`
}

// MicroFrontendSpec defines the desired state of MicroFrontend
type MicroFrontendSpec struct {
	// Reference to a service from which the modules or css would be served.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Service *ServiceReference `json:"service"`

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

	// CacheOptions specifies the cache settings for the PWA, including pre-caching and runtime caching.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	CacheOptions *PWACache `json:"cacheOptions,omitempty"`

	// ImportMap defines module specifier mappings for this microfrontend.
	// These entries are merged at the MicroFrontendClass level. In case of conflicts between
	// microfrontends, first-registered wins (based on creation timestamp).
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ImportMap *ImportMap `json:"importMap,omitempty"`
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

// ImportMap defines module resolution mappings following the Import Map specification.
type ImportMap struct {
	// Imports maps bare module specifiers to paths (relative to this microfrontend's service)
	// or absolute URLs.
	// Example: {"angular": "./bundle/angular.mjs", "react": "https://cdn.example.com/react.js"}
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Imports map[string]string `json:"imports,omitempty"`

	// Scopes allows different module resolutions based on the referrer URL path.
	// Keys are URL path prefixes, values are import maps that apply within that scope.
	// Example: {"/legacy/": {"angular": "./old-angular.mjs"}}
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Scopes map[string]map[string]string `json:"scopes,omitempty"`
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

// MicroFrontendClassReference contains information about the MicroFrontendClass binding
type MicroFrontendClassReference struct {
	// Name of the MicroFrontendClass
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Name string `json:"name"`

	// Namespace of the MicroFrontendClass
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Namespace string `json:"namespace,omitempty"`

	// Accepted indicates if this MicroFrontend is accepted by the class's namespace policy
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Accepted bool `json:"accepted"`
}

// MicroFrontendStatus defines the observed state of MicroFrontend
type MicroFrontendStatus struct {
	// Conditions represent the latest available observations of the MicroFrontend's state
	// +optional
	// +listType=map
	// +listMapKey=type
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Phase represents the current lifecycle phase of the MicroFrontend
	// Possible values: Pending, Ready, Failed, Rejected
	// +optional
	// +kubebuilder:validation:Enum=Pending;Ready;Failed;Rejected
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Phase string `json:"phase,omitempty"`

	// ResolvedServiceURL is the computed URL where the microfrontend is served from
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=status
	ResolvedServiceURL string `json:"resolvedServiceURL,omitempty"`

	// FrontendClassRef indicates which MicroFrontendClass this microfrontend is bound to
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=status
	FrontendClassRef *MicroFrontendClassReference `json:"frontendClassRef,omitempty"`

	// RejectionReason explains why the microfrontend was rejected (namespace policy violation, etc.)
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=status
	RejectionReason string `json:"rejectionReason,omitempty"`

	// ObservedGeneration reflects the generation of the most recently observed MicroFrontend
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=status
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// ImportMapConflicts lists module specifiers that couldn't be registered
	// due to conflicts with other microfrontends (first-registered wins)
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=status
	ImportMapConflicts []ImportMapConflict `json:"importMapConflicts,omitempty"`
}

// ImportMapConflict describes a module specifier conflict in the import map
type ImportMapConflict struct {
	// Specifier is the module name that has a conflict
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Specifier string `json:"specifier"`

	// RequestedPath is what this microfrontend requested
	// +operator-sdk:csv:customresourcedefinitions:type=status
	RequestedPath string `json:"requestedPath"`

	// RegisteredPath is what's actually registered (from another microfrontend)
	// +operator-sdk:csv:customresourcedefinitions:type=status
	RegisteredPath string `json:"registeredPath"`

	// RegisteredBy indicates which microfrontend registered it first (namespace/name)
	// +operator-sdk:csv:customresourcedefinitions:type=status
	RegisteredBy string `json:"registeredBy"`

	// Scope is the scope path where the conflict occurred (empty string for top-level imports)
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Scope string `json:"scope,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MicroFrontend is the Schema for the microfrontends API
type MicroFrontend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MicroFrontendSpec   `json:"spec,omitempty"`
	Status MicroFrontendStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MicroFrontendList contains a list of MicroFrontend
type MicroFrontendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MicroFrontend `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MicroFrontend{}, &MicroFrontendList{})
}

// ResolveServiceURL resolves the ServiceReference to a complete URL
// For in-cluster services (when Name is set), it constructs the URL from name, namespace, port, and scheme
// For external services (when URI is set), it returns the URI directly
func (sr *ServiceReference) ResolveServiceURL(defaultNamespace string) string {
	if sr == nil {
		return ""
	}

	// If URI is specified, use it directly (external service)
	if sr.URI != nil && *sr.URI != "" {
		return *sr.URI
	}

	// If Name is specified, construct in-cluster service URL
	if sr.Name != nil && *sr.Name != "" {
		// Determine namespace (use provided or default)
		namespace := defaultNamespace
		if sr.Namespace != nil && *sr.Namespace != "" {
			namespace = *sr.Namespace
		}

		// Determine scheme (default to http)
		scheme := "http"
		if sr.Scheme != nil && *sr.Scheme != "" {
			scheme = *sr.Scheme
		}

		// Determine port (default to 80)
		port := int32(80)
		if sr.Port != nil {
			port = *sr.Port
		}

		// Determine domain (default to svc.cluster.local)
		domain := "svc.cluster.local"
		if sr.Domain != nil && *sr.Domain != "" {
			domain = *sr.Domain
		}

		// Construct the service URL
		// Format: scheme://service-name.namespace.domain:port
		return scheme + "://" + *sr.Name + "." + namespace + "." + domain + ":" + fmt.Sprint(port)
	}

	return ""
}
