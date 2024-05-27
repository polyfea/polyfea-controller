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
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// MicroFrontendClassSpec defines the desired state of MicroFrontendClass
type MicroFrontendClassSpec struct {
	// BaseUri for which the frontend class will be used
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	BaseUri *string `json:"baseUri"`

	// Title that will be used for the frontend class.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Title *string `json:"title"`

	// CspHeader that will be used for the frontend class, a default will be used if not set.
	// +kubebuilder:default="default-src 'self'; font-src 'self'; script-src 'strict-dynamic' 'nonce-{NONCE_VALUE}'; worker-src 'self'; manifest-src 'self'; style-src 'self' 'strict-dynamic' 'nonce-{NONCE_VALUE}'; style-src-attr 'self' 'unsafe-inline';"
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	CspHeader string `json:"cspHeader,omitempty"`

	// ExtraMetaTags that will be used for the frontend class, none if not set.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ExtraMetaTags []MetaTag `json:"extraMetaTags,omitempty"`

	// ExtraHeaders that will be used for the frontend class, none if not set.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ExtraHeaders []Header `json:"extraHeaders,omitempty"`

	// UserRolesHeader is the name of the header that contains the roles of the user. Defaults to 'x-auth-request-roles'.
	// +kubebuilder:default=x-auth-request-roles
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	UserRolesHeader string `json:"rolesHeader,omitempty"`

	// UserHeader is the name of the header that contains the user id. Defaults to 'x-auth-request-user'.
	// +kubebuilder:default=x-auth-request-user
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	UserHeader string `json:"userHeader,omitempty"`

	// Routing defines the routing for the frontend class from outside of the cluster you can either use a Gateway API or an Ingress.
	// You can also define your own routing by not specifying any of the fields.
	// You can either use a Gateway API or an Ingress.
	// We currently support only basic path prefix routing any customization requires creation of HTTPRoute or Ingress manually.
	// You need to have a service for the operator with label 'app' set to 'polyfea-webserver' and a port with name webserver for the routing to work.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Routing *Routing `json:"routing,omitempty"`

	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ProgressiveWebApp ProgressiveWebApp `json:"progressiveWebApp,omitempty"`
}

// Routing defines the routing for the frontend class from outside of the cluster you can either use a Gateway API or an Ingress.
// +kubebuilder:validation:MaxProperties=1
// +kubebuilder:validation:MinProperties=1
type Routing struct {
	// ParentRefs is the name of the parent refs that the created HTTPRoute will be attached to.
	// If specified an HttpRoute will be created for the frontend class automatically.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ParentRefs []gatewayv1.ParentReference `json:"parentRefs,omitempty"`

	// IngressClassName is the name of the ingress class that will be used for the frontend class.
	// If specified an Ingress will be created for the frontend class automatically.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	IngressClassName *string `json:"ingressClassName,omitempty"`
}

// MetaTag defines the meta tag of the frontend class
type MetaTag struct {
	// Name of the meta tag
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Name string `json:"name"`

	// Content of the meta tag
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Content string `json:"content"`
}

// Header defines the header of the frontend class
type Header struct {
	// Name of the header
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Name string `json:"name"`

	// Value of the header
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Value string `json:"value"`
}

type ProgressiveWebApp struct {

	// TODO: MAKE IT LITTLE BIT MORE SPECIFIC
	// WebAppManifest represents the web app manifest file for the PWA.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	WebAppManifest interface{} `json:"webAppManifest"`

	// CacheOptions specifies the cache settings for the PWA, including pre-caching and runtime caching.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	CacheOptions PWACache `json:"cacheOptions"`

	// Time for reconciliation of the strategies from the frontend side.
	// +kubebuilder:default=1800000
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	PolyfeaSWReconcileInterval *int32 `json:"polyfeaSWReconcileInterval,omitempty"`
}

type PWACache struct {

	// PreCache lists the URLs or resources to be pre-cached when the PWA is installed.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	PreCache []PreCacheEntry `json:"preCache"`

	// CacheRoutes specifies the caching strategies for different URL patterns.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	CacheRoutes []CacheRoute `json:"cacheRoutes"`
}

type PreCacheEntry struct {
	URL      string `json:"url"`
	Revision string `json:"revision"`
}

type CacheRoute struct {

	// Pattern is the URL pattern to which this caching strategy applies.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Pattern string `json:"pattern"`

	// Destination is the optional destination URL for this caching strategy.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Destination *string `json:"destination,omitempty"`

	// Strategy defines the caching strategy to be used for this URL pattern. It defaults to "cache-first".
	// +kubebuilder:default=cache-first
	// +kubebuilder:validation:Enum=cache-first;network-first;cache-only;network-only;stale-while-revalidate;
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Strategy *string `json:"strategy,omitempty"`

	// MaxAgeSeconds specifies the maximum age (in seconds) for cached content.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	MaxAgeSeconds *int32 `json:"maxAgeSeconds,omitempty"`

	// SyncRetentionMinutes specifies the duration (in minutes) to retain synced content in the cache.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	SyncRetentionMinutes *int32 `json:"syncRetentionMinutes,omitempty"`

	// Method specifies the HTTP method to be used with this caching strategy. It defaults to "GET".
	// +kubebuilder:default=GET
	// +kubebuilder:validation:Enum=DELETE;GET;HEAD;PATCH;POST;PUT;
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Method *string `json:"method,omitempty"`

	// Statuses lists the HTTP status codes to be cached. It defaults to [0, 200, 201, 202, 204].
	// +kubebuilder:default={0,200,201,202,204}
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Statuses []int32 `json:"statuses,omitempty"`
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
