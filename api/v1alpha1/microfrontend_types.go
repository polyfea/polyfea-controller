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

// NamespacedReference refers to a named resource, optionally in a specific namespace.
// If Namespace is omitted, the referencing resource's own namespace is assumed.
type NamespacedReference struct {
	// Name of the referenced resource.
	// +kubebuilder:validation:MaxLength=253
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Name string `json:"name"`

	// Namespace of the referenced resource.
	// Defaults to the namespace of the referencing resource if not specified.
	// +optional
	// +kubebuilder:validation:MaxLength=63
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Namespace *string `json:"namespace,omitempty"`
}

// ServiceReference defines how to reach the service hosting the micro frontend
// +kubebuilder:validation:XValidation:rule="(has(self.name) && size(self.name) > 0) != (has(self.uri) && size(self.uri) > 0)",message="Either 'name' or 'uri' must be specified, but not both"
type ServiceReference struct {
	// Name of the Kubernetes service (mutually exclusive with URI)
	// +optional
	// +kubebuilder:validation:MaxLength=253
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Name *string `json:"name,omitempty"`

	// URI for external services (mutually exclusive with Name)
	// Should include schema (http:// or https://)
	// +optional
	// +kubebuilder:validation:MaxLength=2048
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	URI *string `json:"uri,omitempty"`

	// Namespace of the service. Defaults to the MicroFrontend's namespace if not specified.
	// Only used when Name is set.
	// +optional
	// +kubebuilder:validation:MaxLength=63
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
	// +kubebuilder:validation:MaxLength=5
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Scheme *string `json:"scheme,omitempty"`

	// Domain is the cluster domain suffix. Defaults to svc.cluster.local if not specified.
	// Only used when Name is set. Allows customization for different cluster implementations.
	// +optional
	// +kubebuilder:default=svc.cluster.local
	// +kubebuilder:validation:MaxLength=253
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

	// Relative path to the module file within the service.
	// +kubebuilder:validation:MaxLength=2048
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ModulePath *string `json:"modulePath"`

	// Relative path to the static files within the service.
	// +kubebuilder:validation:MaxItems=64
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	StaticResources []StaticResources `json:"staticPaths,omitempty"`

	// FrontendClass is a reference to the MicroFrontendClass to bind to.
	// If Namespace is omitted, the MicroFrontend's own namespace is assumed.
	// +kubebuilder:default={name: "polyfea-controller-default"}
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	FrontendClass NamespacedReference `json:"frontendClass"`

	// List of dependencies that should be loaded before this micro frontend.
	// +kubebuilder:validation:MaxItems=64
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	DependsOn []string `json:"dependsOn,omitempty"`

	// CacheOptions specifies the cache settings for the PWA, including pre-caching and runtime caching.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	CacheOptions *PWACache `json:"cacheOptions,omitempty"`

	// ImportMap defines module specifier mappings for this microfrontend.
	// Entries are merged at the MicroFrontendClass level using optional (global, skip-if-exists)
	// and scoped (isolated per microfrontend) strategies.
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ImportMap *ImportMap `json:"importMap,omitempty"`

	// CacheBustingHash is a token embedded as a path segment in all proxied URLs for this
	// microfrontend: ./polyfea/proxy/{namespace}/{name}/{CacheBustingHash}/{path}
	// The proxy always strips this segment before forwarding to the backend service.
	// Change this value to force browsers to re-fetch all proxied resources (cache busting).
	// Defaults to "nohash". Only applies to proxied resources (Proxy: true).
	// +kubebuilder:default=nohash
	// +kubebuilder:validation:MaxLength=64
	// +kubebuilder:validation:Pattern=`^[a-zA-Z0-9_.-]+$`
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	CacheBustingHash string `json:"cacheBustingHash,omitempty"`
}

// StaticResources defines the static resources that should be loaded before this micro frontend.
type StaticResources struct {
	// Kind defines the kind of the static resource can be script, stylesheet, or any other `link` element.
	// +kubebuilder:validation:MaxLength=64
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Kind string `json:"kind"`

	// +kubebuilder:validation:MaxLength=2048
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Path string `json:"path"`

	// +kubebuilder:validation:MaxItems=64
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

// ImportMap defines module resolution mappings for a microfrontend.
type ImportMap struct {
	// Optional maps bare module specifiers to paths (relative to this microfrontend's service)
	// or absolute URLs. Entries are added to the global import map. If a specifier is already
	// registered by another MicroFrontend, it is silently skipped (first-registered-wins).
	// Example: {"react": "./react.js", "lodash": "https://cdn.example.com/lodash.js"}
	// +optional
	// +kubebuilder:validation:MaxProperties=256
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Optional map[string]string `json:"optional,omitempty"`

	// Scoped maps bare module specifiers to paths (relative to this microfrontend's service)
	// or absolute URLs. Entries are placed under this MicroFrontend's scope (derived from its
	// proxy path), ensuring isolation. When scripts loaded from this MicroFrontend import a
	// specifier, the scoped entry is used instead of the global one.
	// Example: {"react": "./react-v17.js"}
	// +optional
	// +kubebuilder:validation:MaxProperties=256
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Scoped map[string]string `json:"scoped,omitempty"`
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
	// +kubebuilder:validation:MaxLength=253
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Name string `json:"name"`

	// Namespace of the MicroFrontendClass
	// +optional
	// +kubebuilder:validation:MaxLength=63
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
	// +kubebuilder:validation:MaxLength=2048
	// +operator-sdk:csv:customresourcedefinitions:type=status
	ResolvedServiceURL string `json:"resolvedServiceURL,omitempty"`

	// FrontendClassRef indicates which MicroFrontendClass this microfrontend is bound to
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=status
	FrontendClassRef *MicroFrontendClassReference `json:"frontendClassRef,omitempty"`

	// RejectionReason explains why the microfrontend was rejected (namespace policy violation, etc.)
	// +optional
	// +kubebuilder:validation:MaxLength=1024
	// +operator-sdk:csv:customresourcedefinitions:type=status
	RejectionReason string `json:"rejectionReason,omitempty"`

	// ObservedGeneration reflects the generation of the most recently observed MicroFrontend
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=status
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
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
