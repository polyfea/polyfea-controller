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

// buildModulePath constructs the path for a microfrontend module, either proxied or direct.
func buildModulePath(namespace, name, path string, proxy bool, service *v1alpha1.ServiceReference) *string {
	if proxy {
		result := buildProxyPath(namespace, name, path)
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
