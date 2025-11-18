package polyfea

import (
	_ "embed"
	"testing"
)

// TestGetMicrofrontendAndBase tests the retrieval of microfrontend class and base URI.
func TestGetMicrofrontendAndBase(t *testing.T) {
	mfcRepository := setupMfcRepository()

	tests := []struct {
		name            string
		requestPath     string
		expectBasePath  string
		expectClassName string
	}{
		{
			name:            "Nonexistent path should default to default class",
			requestPath:     "/nonexistent",
			expectBasePath:  "/",
			expectClassName: "default",
		},
		{
			name:            "Path matching fea with asset",
			requestPath:     "/fea/asset",
			expectBasePath:  "/fea/",
			expectClassName: "fea",
		},
		{
			name:            "Path matching fea without asset",
			requestPath:     "/fea",
			expectBasePath:  "/fea/",
			expectClassName: "fea",
		},
		{
			name:            "Path matching feature",
			requestPath:     "/feature",
			expectBasePath:  "/feature/",
			expectClassName: "feature",
		},
		{
			name:            "Path matching feature with asset",
			requestPath:     "/feature/asset",
			expectBasePath:  "/feature/",
			expectClassName: "feature",
		},
		{
			name:            "Path not matching any class should default to default",
			requestPath:     "/fea-nix",
			expectBasePath:  "/",
			expectClassName: "default",
		},
		{
			name:            "Nested path not matching any class",
			requestPath:     "/other/qweqwesop",
			expectBasePath:  "/",
			expectClassName: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			basePath, microfrontend, err := getMicrofrontendClassAndBase(tt.requestPath, mfcRepository)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if basePath != tt.expectBasePath {
				t.Errorf("expected base path %q, got %q", tt.expectBasePath, basePath)
			}
			if microfrontend.Name != tt.expectClassName {
				t.Errorf("expected class name %q, got %q", tt.expectClassName, microfrontend.Name)
			}
		})
	}
}

// TestGetMicrofrontendAndBaseEdgeCases tests edge cases and error handling.
func TestGetMicrofrontendAndBaseEdgeCases(t *testing.T) {
	mfcRepository := setupMfcRepository()

	tests := []struct {
		name            string
		requestPath     string
		expectBasePath  string
		expectClassName string
	}{
		{
			name:            "Empty requestPath",
			requestPath:     "",
			expectBasePath:  "/",
			expectClassName: "default",
		},
		{
			name:            "Special characters in requestPath",
			requestPath:     "/!@#$%^&*()",
			expectBasePath:  "/",
			expectClassName: "default",
		},
		{
			name:            "Path longer than any base URI",
			requestPath:     "/this/path/is/way/too/long/to/match/anything",
			expectBasePath:  "/",
			expectClassName: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			basePath, microfrontend, err := getMicrofrontendClassAndBase(tt.requestPath, mfcRepository)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if basePath != tt.expectBasePath {
				t.Errorf("expected base path %q, got %q", tt.expectBasePath, basePath)
			}
			if microfrontend.Name != tt.expectClassName {
				t.Errorf("expected class name %q, got %q", tt.expectClassName, microfrontend.Name)
			}
		})
	}
}

func ptr(s string) *string {
	return &s
}
