package polyfea

import (
	"testing"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
	"github.com/rs/zerolog"
)

func TestServeAppWebManifestReturnsExpectedManifest(t *testing.T) {
	// Test that the `serveAppWebManifest` method returns the expected manifest
	// Arrange
	pwa := NewProgressiveWebApplication(&zerolog.Logger{}, repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]())

	mfc := createTestMicroFrontendClass("polyfea", "/someother")
	mfc.Spec.ProgressiveWebApp = &v1alpha1.ProgressiveWebApp{
		WebAppManifest: &v1alpha1.WebAppManifest{
			Name: &[]string{"Test"}[0],
			Icons: []v1alpha1.PWAIcon{
				{
					Type:  &[]string{"image/png"}[0],
					Sizes: &[]string{"192x192"}[0],
					Src:   &[]string{"icon.png"}[0],
				},
			},
			StartUrl: &[]string{"/"}[0],
			Display:  &[]string{"standalone"}[0],
		},
	}

	expected := mfc.Spec.ProgressiveWebApp.WebAppManifest

	// Act
	actual := pwa.serveAppWebManifest(mfc)

	// Assert
	assertWebAppManifestEqual(t, expected, actual)
}

func TestServeProxyConfigReturnsExpectedConfigForAllRelevantMicrofrontends(t *testing.T) {
	// Test that the `getProxyConfig` method returns the expected config for all relevant microfrontends
	// Arrange
	microFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()

	mfc := createTestMicroFrontendClass("polyfea", "/someother")
	mfc.Spec.ProgressiveWebApp = &v1alpha1.ProgressiveWebApp{
		CacheOptions: &v1alpha1.PWACache{
			PreCache: []v1alpha1.PreCacheEntry{
				{URL: &[]string{"/test-class"}[0]},
			},
			CacheRoutes: []v1alpha1.CacheRoute{
				{Pattern: &[]string{"/cache-route-class"}[0], Strategy: &[]string{"cache-first"}[0]},
			},
		},
	}

	mf := createTestMicroFrontend("polyfea1", []string{}, "test-module", "polyfea", true)
	mf.Spec.CacheOptions = &v1alpha1.PWACache{
		PreCache: []v1alpha1.PreCacheEntry{
			{URL: &[]string{"/test"}[0]},
		},
	}
	microFrontendRepository.Store(mf)

	mf2 := createTestMicroFrontend("polyfea2", []string{}, "test-module", "polyfea", true)
	mf2.Spec.CacheOptions = &v1alpha1.PWACache{
		PreCache: []v1alpha1.PreCacheEntry{
			{URL: &[]string{"/test2"}[0]},
		},
		CacheRoutes: []v1alpha1.CacheRoute{
			{Pattern: &[]string{"/cache-route"}[0], Strategy: &[]string{"network-first"}[0]},
		},
	}
	microFrontendRepository.Store(mf2)

	pwa := NewProgressiveWebApplication(&zerolog.Logger{}, microFrontendRepository)

	expected := &ProxyConfigResponse{
		PreCache: []v1alpha1.PreCacheEntry{
			{URL: &[]string{"/test-class"}[0]},
			{URL: buildPreCachePath(mf, *mf.Spec.CacheOptions.PreCache[0].URL)},
			{URL: buildPreCachePath(mf2, *mf2.Spec.CacheOptions.PreCache[0].URL)},
		},
		Routes: []CacheRouteResponse{
			{CacheRoute: mfc.Spec.ProgressiveWebApp.CacheOptions.CacheRoutes[0]},
			{CacheRoute: mf2.Spec.CacheOptions.CacheRoutes[0], Prefix: buildPreCachePath(mf2, "")},
		},
	}

	// Act
	actual, err := pwa.getProxyConfig(mfc)

	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assertProxyConfigEqual(t, expected, actual)
}

// Helper functions for assertions
func assertWebAppManifestEqual(t *testing.T, expected, actual *v1alpha1.WebAppManifest) {
	if *actual.Name != *expected.Name {
		t.Errorf("Expected %s, got %s", *expected.Name, *actual.Name)
	}
	if len(actual.Icons) != len(expected.Icons) {
		t.Errorf("Expected %d icons, got %d", len(expected.Icons), len(actual.Icons))
	}
	for i := range actual.Icons {
		if *actual.Icons[i].Src != *expected.Icons[i].Src {
			t.Errorf("Expected %s, got %s", *expected.Icons[i].Src, *actual.Icons[i].Src)
		}
		if *actual.Icons[i].Sizes != *expected.Icons[i].Sizes {
			t.Errorf("Expected %s, got %s", *expected.Icons[i].Sizes, *actual.Icons[i].Sizes)
		}
		if *actual.Icons[i].Type != *expected.Icons[i].Type {
			t.Errorf("Expected %s, got %s", *expected.Icons[i].Type, *actual.Icons[i].Type)
		}
	}
	if *actual.StartUrl != *expected.StartUrl {
		t.Errorf("Expected %s, got %s", *expected.StartUrl, *actual.StartUrl)
	}
	if *actual.Display != *expected.Display {
		t.Errorf("Expected %s, got %s", *expected.Display, *actual.Display)
	}
}

func assertProxyConfigEqual(t *testing.T, expected, actual *ProxyConfigResponse) {
	if len(actual.PreCache) != len(expected.PreCache) {
		t.Errorf("Expected %d, got %d", len(expected.PreCache), len(actual.PreCache))
	}
	for i, entry := range actual.PreCache {
		if *entry.URL != *expected.PreCache[i].URL {
			t.Errorf("Expected %s, got %s", *expected.PreCache[i].URL, *entry.URL)
		}
	}
	if len(actual.Routes) != len(expected.Routes) {
		t.Errorf("Expected %d, got %d", len(expected.Routes), len(actual.Routes))
	}
	for i, route := range actual.Routes {
		if (route.Prefix == nil && expected.Routes[i].Prefix != nil) ||
			(route.Prefix != nil && expected.Routes[i].Prefix == nil) ||
			(route.Prefix != nil && expected.Routes[i].Prefix != nil && *route.Prefix != *expected.Routes[i].Prefix) {
			t.Errorf("Expected %s, got %s", *expected.Routes[i].Prefix, *route.Prefix)
		}
		if *route.Pattern != *expected.Routes[i].Pattern {
			t.Errorf("Expected %s, got %s", *expected.Routes[i].Pattern, *route.Pattern)
		}
		if *route.Strategy != *expected.Routes[i].Strategy {
			t.Errorf("Expected %s, got %s", *expected.Routes[i].Strategy, *route.Strategy)
		}
	}
}
