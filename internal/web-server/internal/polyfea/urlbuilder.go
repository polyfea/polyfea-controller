package polyfea

import (
	"github.com/polyfea/polyfea-controller/api/v1alpha1"
)

// buildProxyPath constructs a proxy URL path for a microfrontend resource.
func buildProxyPath(namespace, name, path string) string {
	return "./polyfea/proxy/" + namespace + "/" + name + "/" + path
}

// joinURL joins a base URL and path, ensuring exactly one "/" separator.
func joinURL(baseURL, path string) string {
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] != '/' && len(path) > 0 && path[0] != '/' {
		return baseURL + "/" + path
	}
	return baseURL + path
}

// appendVersionFragment appends a URL fragment (#version) to path when version is non-empty.
// Used to bust the browser ES module registry cache without changing the HTTP request URL.
func appendVersionFragment(path, version string) string {
	if version == "" {
		return path
	}
	return path + "#" + version
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
