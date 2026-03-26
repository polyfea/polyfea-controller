package polyfea

import (
	"strings"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
)

// hashOrDefault returns hash if non-empty, otherwise "nohash".
func hashOrDefault(hash string) string {
	if hash == "" {
		return "nohash"
	}
	return hash
}

// buildProxyPath constructs a proxy URL path for a microfrontend resource.
// The cache-busting hash segment is always present (defaults to "nohash" when empty).
// Leading slashes are stripped from path to avoid double slashes in the result.
func buildProxyPath(namespace, name, cacheBustingHash, path string) string {
	return "./polyfea/proxy/" + namespace + "/" + name + "/" +
		hashOrDefault(cacheBustingHash) + "/" + strings.TrimLeft(path, "/")
}

// buildScopeKey constructs the import map scope key for a microfrontend.
// Scope keys are hash-free URL prefixes: a script served from
// ./polyfea/proxy/ns/mf/{hash}/app.js still starts with ./polyfea/proxy/ns/mf/,
// so the scope applies regardless of which hash is active.
func buildScopeKey(namespace, name string) string {
	return "./polyfea/proxy/" + namespace + "/" + name + "/"
}

// joinURL joins a base URL and path, ensuring exactly one "/" separator.
func joinURL(baseURL, path string) string {
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] != '/' && len(path) > 0 && path[0] != '/' {
		return baseURL + "/" + path
	}
	return baseURL + path
}

// buildModulePath constructs the path for a microfrontend module, either proxied or direct.
func buildModulePath(namespace, name, path, cacheBustingHash string, proxy bool, service *v1alpha1.ServiceReference) *string {
	if proxy {
		result := buildProxyPath(namespace, name, cacheBustingHash, path)
		return &result
	}
	// For non-proxied services, combine service URL with path
	if service != nil {
		baseURL := service.ResolveServiceURL(namespace)
		if baseURL != "" {
			result := joinURL(baseURL, path)
			return &result
		}
	}
	// Fallback to just path if service is not provided
	return &path
}
