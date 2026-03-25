package polyfea

import (
	"testing"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
)

func TestBuildProxyPath(t *testing.T) {
	tests := []struct {
		name, namespace, mfName, path, want string
	}{
		{"basic", "default", "mf", "module.js", "./polyfea/proxy/default/mf/module.js"},
		{"nested path", "ns", "app", "assets/style.css", "./polyfea/proxy/ns/app/assets/style.css"},
		{"leading slash stripped", "ns", "mf", "/imports/@lit/context/", "./polyfea/proxy/ns/mf/imports/@lit/context/"},
		{"empty path", "ns", "mf", "", "./polyfea/proxy/ns/mf/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildProxyPath(tt.namespace, tt.mfName, tt.path)
			if got != tt.want {
				t.Errorf("buildProxyPath() = %q, want %q", got, tt.want)
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
	t.Run("proxy enabled", func(t *testing.T) {
		got := buildModulePath("ns", "mf", "app.js", "", true, nil)
		want := "./polyfea/proxy/ns/mf/app.js"
		if got == nil || *got != want {
			t.Errorf("buildModulePath(proxy=true) = %v, want %q", got, want)
		}
	})

	t.Run("proxy enabled with version", func(t *testing.T) {
		got := buildModulePath("ns", "mf", "app.js", "abc123def456", true, nil)
		want := "./polyfea/proxy/ns/mf/app.js#abc123def456"
		if got == nil || *got != want {
			t.Errorf("buildModulePath(proxy=true, version) = %v, want %q", got, want)
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

func TestAppendVersionFragment(t *testing.T) {
	tests := []struct {
		name, path, version, want string
	}{
		{"empty version", "./polyfea/proxy/ns/mf/app.js", "", "./polyfea/proxy/ns/mf/app.js"},
		{"with version", "./polyfea/proxy/ns/mf/app.js", "abc123", "./polyfea/proxy/ns/mf/app.js#abc123"},
		{"trailing slash not versioned", "./polyfea/proxy/ns/mf/imports/@lit/context/", "abc123", "./polyfea/proxy/ns/mf/imports/@lit/context/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := appendVersionFragment(tt.path, tt.version)
			if got != tt.want {
				t.Errorf("appendVersionFragment(%q, %q) = %q, want %q", tt.path, tt.version, got, tt.want)
			}
		})
	}
}

func strPtr(s string) *string { return &s }
