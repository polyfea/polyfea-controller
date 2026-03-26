package polyfea

import (
	"testing"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
)

func TestBuildProxyPath(t *testing.T) {
	tests := []struct {
		name, namespace, mfName, hash, path, want string
	}{
		{"basic no hash", "default", "mf", "", "module.js", "./polyfea/proxy/default/mf/nohash/module.js"},
		{"basic with hash", "default", "mf", "v1", "module.js", "./polyfea/proxy/default/mf/v1/module.js"},
		{"nested path", "ns", "app", "abc123", "assets/style.css", "./polyfea/proxy/ns/app/abc123/assets/style.css"},
		{"leading slash stripped", "ns", "mf", "v2", "/imports/@lit/context/", "./polyfea/proxy/ns/mf/v2/imports/@lit/context/"},
		{"empty path", "ns", "mf", "v1", "", "./polyfea/proxy/ns/mf/v1/"},
		{"trailing slash preserved", "ns", "mf", "v1", "imports/@lit/context/", "./polyfea/proxy/ns/mf/v1/imports/@lit/context/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildProxyPath(tt.namespace, tt.mfName, tt.hash, tt.path)
			if got != tt.want {
				t.Errorf("buildProxyPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildScopeKey(t *testing.T) {
	tests := []struct {
		name, namespace, mfName, want string
	}{
		{"basic", "default", "mf1", "./polyfea/proxy/default/mf1/"},
		{"cross-namespace", "ns-a", "my-app", "./polyfea/proxy/ns-a/my-app/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildScopeKey(tt.namespace, tt.mfName)
			if got != tt.want {
				t.Errorf("buildScopeKey() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHashOrDefault(t *testing.T) {
	tests := []struct {
		name, hash, want string
	}{
		{"empty hash uses nohash", "", "nohash"},
		{"explicit hash returned as-is", "v1", "v1"},
		{"nohash passthrough", "nohash", "nohash"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hashOrDefault(tt.hash)
			if got != tt.want {
				t.Errorf("hashOrDefault(%q) = %q, want %q", tt.hash, got, tt.want)
			}
		})
	}
}

func TestJoinURL(t *testing.T) {
	tests := []struct {
		name, base, path, want string
	}{
		{"both without slash", "http://svc", "path", "http://svc/path"},
		{"base with trailing slash", "http://svc/", "path", "http://svc/path"},
		{"path with leading slash", "http://svc", "/path", "http://svc/path"},
		{"both with slash", "http://svc/", "/path", "http://svc//path"},
		{"empty base", "", "path", "path"},
		{"empty path", "http://svc", "", "http://svc"},
		{"both empty", "", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := joinURL(tt.base, tt.path)
			if got != tt.want {
				t.Errorf("joinURL(%q, %q) = %q, want %q", tt.base, tt.path, got, tt.want)
			}
		})
	}
}

func TestBuildModulePath(t *testing.T) {
	t.Run("proxy enabled no hash", func(t *testing.T) {
		got := buildModulePath("ns", "mf", "app.js", "", true, nil)
		want := "./polyfea/proxy/ns/mf/nohash/app.js"
		if got == nil || *got != want {
			t.Errorf("buildModulePath(proxy=true, no hash) = %v, want %q", got, want)
		}
	})

	t.Run("proxy enabled with hash", func(t *testing.T) {
		got := buildModulePath("ns", "mf", "app.js", "abc123def456", true, nil)
		want := "./polyfea/proxy/ns/mf/abc123def456/app.js"
		if got == nil || *got != want {
			t.Errorf("buildModulePath(proxy=true, hash) = %v, want %q", got, want)
		}
	})

	t.Run("proxy disabled with service", func(t *testing.T) {
		svc := &v1alpha1.ServiceReference{
			URI: strPtr("https://cdn.example.com"),
		}
		got := buildModulePath("ns", "mf", "app.js", "", false, svc)
		want := "https://cdn.example.com/app.js"
		if got == nil || *got != want {
			t.Errorf("buildModulePath(proxy=false, service) = %v, want %q", got, want)
		}
	})

	t.Run("proxy disabled without service", func(t *testing.T) {
		got := buildModulePath("ns", "mf", "app.js", "", false, nil)
		want := "app.js"
		if got == nil || *got != want {
			t.Errorf("buildModulePath(proxy=false, nil service) = %v, want %q", got, want)
		}
	})
}

func strPtr(s string) *string { return &s }
