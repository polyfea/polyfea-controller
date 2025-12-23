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
)

// NamespacePolicyType defines namespace selection behavior
// +kubebuilder:validation:Enum=All;Same;FromNamespaces
type NamespacePolicyType string

const (
	// NamespaceFromAll allows MicroFrontends from all namespaces
	NamespaceFromAll NamespacePolicyType = "All"

	// NamespaceFromSame allows only MicroFrontends from the same namespace as the MicroFrontendClass
	NamespaceFromSame NamespacePolicyType = "Same"

	// NamespaceFromNamespaces allows MicroFrontends from specific namespaces listed in Namespaces field
	NamespaceFromNamespaces NamespacePolicyType = "FromNamespaces"
)

// NamespacePolicy defines which namespaces can attach MicroFrontends to this class
type NamespacePolicy struct {
	// From defines namespace selection behavior
	// +kubebuilder:validation:Enum=All;Same;FromNamespaces
	// +kubebuilder:default=All
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	From NamespacePolicyType `json:"from"`

	// Namespaces is a list of namespaces from which MicroFrontends can be attached
	// Only used when From is "FromNamespaces"
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Namespaces []string `json:"namespaces,omitempty"`
}

// MicroFrontendClassSpec defines the desired state of MicroFrontendClass
type MicroFrontendClassSpec struct {
	// BaseUri for which the frontend class will be used
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	BaseUri *string `json:"baseUri"`

	// Title that will be used for the frontend class.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Title *string `json:"title"`

	// NamespacePolicy defines which namespaces can attach MicroFrontends to this class
	// Defaults to allowing all namespaces
	// +optional
	// +kubebuilder:default={from: "All"}
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	NamespacePolicy *NamespacePolicy `json:"namespacePolicy,omitempty"`

	// CspHeader that will be used for the frontend class, a default will be used if not set.
	// +kubebuilder:default="default-src 'self'; font-src 'self'; script-src 'strict-dynamic' 'nonce-{NONCE_VALUE}'; worker-src 'self'; manifest-src 'self'; style-src 'self' 'nonce-{NONCE_VALUE}'; style-src-attr 'self' 'unsafe-inline';"
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

	// ProgressiveWebApp defines the configuration settings for a Progressive Web Application (PWA).
	// It includes specifications for the web app manifest and cache options, which are crucial for the PWA's functionality and performance.
	// This field is optional and can be omitted if not needed.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ProgressiveWebApp *ProgressiveWebApp `json:"progressiveWebApp,omitempty"`
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

// ProgressiveWebApp defines the configuration settings for a Progressive Web Application (PWA).
// This struct includes specifications for the web app manifest, caching options, and reconciliation interval,
// which are critical for the PWA's functionality, performance, and synchronization with frontend updates.
type ProgressiveWebApp struct {
	// WebAppManifest represents the web app manifest file for the PWA.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	WebAppManifest *WebAppManifest `json:"webAppManifest"`

	// CacheOptions specifies the cache settings for the PWA, including pre-caching and runtime caching.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	CacheOptions *PWACache `json:"cacheOptions,omitempty"`

	// Time for reconciliation of the strategies from the frontend side.
	// +kubebuilder:default=1800000
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	PolyfeaSWReconcileInterval *int32 `json:"polyfeaSWReconcileInterval,omitempty"`
}

// WebAppManifest represents the web app manifest file for the PWA.
type WebAppManifest struct {
	// Read more here: https://developer.mozilla.org/en-US/docs/Web/Manifest/name
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Name *string `json:"name"`

	// Read more here: https://developer.mozilla.org/en-US/docs/Web/Manifest/icons
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Icons []PWAIcon `json:"icons"`

	// Read more here: https://developer.mozilla.org/en-US/docs/Web/Manifest/start_url
	// URL needs to be relative to the base URL of the frontend class.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	StartUrl *string `json:"start_url"`

	// Read more here: https://developer.mozilla.org/en-US/docs/Web/Manifest/display
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Display *string `json:"display"`

	// Read more here: https://developer.mozilla.org/en-US/docs/Web/Manifest/display_override
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	DisplayOverride []string `json:"display_override,omitempty"`
}

// Read more here: https://developer.mozilla.org/en-US/docs/Web/Manifest/icons
type PWAIcon struct {
	// A string containing space-separated image dimensions using the same syntax as the sizes attribute.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Sizes *string `json:"sizes"`

	// The path to the image file. If src is a relative URL, the base URL will be the URL of the manifest.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Src *string `json:"src"`

	// A hint as to the media type of the image. The purpose of this member is to allow a user agent to quickly ignore images with media types it does not support.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Type *string `json:"type"`

	// Defines the purpose of the image, for example if the image is intended to serve some special purpose in the context of the host OS (i.e., for better integration).
	// purpose can have one or more of the following values, separated by spaces:
	// 	monochrome: A user agent can present this icon where a monochrome icon with a solid fill is needed. The color information in the icon is discarded and only the alpha data is used. The icon can then be used by the user agent like a mask over any solid fill.
	// 	maskable: The image is designed with icon masks and safe zone in mind, such that any part of the image outside the safe zone can safely be ignored and masked away by the user agent.
	// 	any: The user agent is free to display the icon in any context (this is the default value).
	// +kubebuilder:validation:Enum=monochrome;maskable;any;
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Purpose *string `json:"purpose,omitempty"`
}

// PWACache defines the caching options for a Progressive Web Application (PWA).
// This struct includes configurations for both pre-caching and runtime caching strategies, which are essential for improving the performance and offline capabilities of the PWA.
type PWACache struct {
	// PreCache lists the URLs or resources to be pre-cached when the PWA is installed.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	PreCache []PreCacheEntry `json:"preCache,omitempty"`

	// CacheRoutes specifies the caching strategies for different URL patterns.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	CacheRoutes []CacheRoute `json:"cacheRoutes,omitempty"`
}

// PreCacheEntry represents an individual entry in the pre-cache list for a Progressive Web Application (PWA).
// Each entry specifies a URL to be cached and an optional revision identifier to manage cache updates and invalidation.
type PreCacheEntry struct {
	// URL specifies the resource URL that should be pre-cached. This URL points to the asset that needs to be available offline, ensuring it is cached during the installation of the PWA.
	// URL needs to be relative to the base URL of the frontend class.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	URL *string `json:"url"`

	// Revision is an optional field that specifies a revision identifier for the resource.
	// The revision helps in cache management by allowing the service worker to recognize and update cached assets when their content changes. This ensures users always have access to the most up-to-date resources.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Revision *string `json:"revision,omitempty"`
}

// CacheRoute defines the caching strategy for a specific URL pattern within a Progressive Web Application (PWA).
// This struct allows for fine-tuned control over how different network requests are handled, enhancing performance, reliability, and offline capabilities based on the application's requirements.
type CacheRoute struct {
	// Pattern is the URL pattern to which this caching strategy applies.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Pattern *string `json:"pattern"`

	// Destination is the optional destination URL for this caching strategy.
	// You can find the list of possible values here: https://developer.mozilla.org/en-US/docs/Web/API/Request/destination
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
	// Conditions represent the latest available observations of the MicroFrontendClass's state
	// +optional
	// +listType=map
	// +listMapKey=type
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Phase represents the current lifecycle phase of the MicroFrontendClass
	// Possible values: Ready, Invalid
	// +optional
	// +kubebuilder:validation:Enum=Ready;Invalid
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Phase string `json:"phase,omitempty"`

	// AcceptedMicroFrontends counts how many MicroFrontends are currently bound to this class
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=status
	AcceptedMicroFrontends int32 `json:"acceptedMicroFrontends,omitempty"`

	// RejectedMicroFrontends counts how many MicroFrontends were rejected by namespace policy
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=status
	RejectedMicroFrontends int32 `json:"rejectedMicroFrontends,omitempty"`

	// ObservedGeneration reflects the generation of the most recently observed MicroFrontendClass
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=status
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MicroFrontendClass is the Schema for the microfrontendclasses API
type MicroFrontendClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MicroFrontendClassSpec   `json:"spec,omitempty"`
	Status MicroFrontendClassStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MicroFrontendClassList contains a list of MicroFrontendClass
type MicroFrontendClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MicroFrontendClass `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MicroFrontendClass{}, &MicroFrontendClassList{})
}

// IsNamespaceAllowed checks if a namespace is allowed by the NamespacePolicy
func (mfc *MicroFrontendClass) IsNamespaceAllowed(namespace string) bool {
	// If no policy is set, default to allowing all namespaces
	if mfc.Spec.NamespacePolicy == nil {
		return true
	}

	switch mfc.Spec.NamespacePolicy.From {
	case NamespaceFromAll:
		return true
	case NamespaceFromSame:
		return namespace == mfc.Namespace
	case NamespaceFromNamespaces:
		for _, ns := range mfc.Spec.NamespacePolicy.Namespaces {
			if ns == namespace {
				return true
			}
		}
		return false
	default:
		// Default to All if an unknown value is provided
		return true
	}
}
