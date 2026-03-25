package polyfea

import (
	"strings"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
)

// buildProxyPath constructs a proxy URL path for a microfrontend resource.
// Leading slashes are stripped from path to avoid double slashes in the result.
func buildProxyPath(namespace, name, path string) string {
	return "./polyfea/proxy/" + namespace + "/" + name + "/" + strings.TrimLeft(path, "/")
}

// joinURL joins a base URL and path, ensuring exactly one "/" separator.
func joinURL(baseURL, path string) string {
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] != '/' && len(path) > 0 && path[0] != '/' {
		return baseURL + "/" + path
	}
	return baseURL + path
}

// appendVersionFragment appends a query-string version (?v=version) to path when version is non-empty.
// The proxy strips query strings before forwarding to the backend, so the version token is invisible
// to the origin server but changes the URL seen by the browser's HTTP cache, busting it on update.
// Paths ending in "/" are left unchanged: they are import map namespace-like entries and the
// browser requires the address to also end in "/" — appending anything would break that.
func appendVersionFragment(path, version string) string {
	if version == "" || strings.HasSuffix(path, "/") {
		return path
	}
	return path + "?v=" + version
}

// buildModulePath constructs the path for a microfrontend module, either proxied or direct.
func buildModulePath(namespace, name, path, version string, proxy bool, service *v1alpha1.ServiceReference) *string {
	if proxy {
		result := appendVersionFragment(buildProxyPath(namespace, name, path), version)
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
