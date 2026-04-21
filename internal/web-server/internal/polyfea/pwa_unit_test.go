package polyfea

import (
	"sort"
	"testing"

	"github.com/go-logr/logr"
	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/internal/repository"
)

func TestServeAppWebManifestReturnsExpectedManifest(t *testing.T) {
	// Test that the `serveAppWebManifest` method returns the expected manifest
	// Arrange
	pwa := NewProgressiveWebApplication(&logr.Logger{}, repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]())

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
		ServiceWorker: &v1alpha1.ServiceWorker{
			PreCache: []v1alpha1.PreCacheEntry{
				{URL: &[]string{"/test-class"}[0]},
			},
			CacheRoutes: []v1alpha1.CacheRoute{
				{Pattern: &[]string{"/cache-route-class"}[0], Strategy: &[]string{"cache-first"}[0]},
			},
		},
	}

	mf := createTestMicroFrontend("polyfea1", []string{}, "polyfea", true)
	mf.Spec.ServiceWorker = &v1alpha1.MicroFrontendServiceWorker{
		PreCache: []v1alpha1.PreCacheEntry{
			{URL: &[]string{"/test"}[0]},
		},
	}
	err := microFrontendRepository.Store(mf)
	if err != nil {
		t.Fatal(err)
	}

	mf2 := createTestMicroFrontend("polyfea2", []string{}, "polyfea", true)
	mf2.Spec.ServiceWorker = &v1alpha1.MicroFrontendServiceWorker{
		PreCache: []v1alpha1.PreCacheEntry{
			{URL: &[]string{"/test2"}[0]},
		},
		CacheRoutes: []v1alpha1.CacheRoute{
			{Pattern: &[]string{"/cache-route"}[0], Strategy: &[]string{"network-first"}[0]},
		},
	}

	err = microFrontendRepository.Store(mf2)
	if err != nil {
		t.Fatal(err)
	}

	pwa := NewProgressiveWebApplication(&logr.Logger{}, microFrontendRepository)

	expected := &ProxyConfigResponse{
		PreCache: []v1alpha1.PreCacheEntry{
			{URL: &[]string{"/test-class"}[0]},
			{URL: buildPreCachePath(mf, *mf.Spec.ServiceWorker.PreCache[0].URL)},
			{URL: buildPreCachePath(mf2, *mf2.Spec.ServiceWorker.PreCache[0].URL)},
		},
		Routes: []CacheRouteResponse{
			{CacheRoute: mfc.Spec.ProgressiveWebApp.ServiceWorker.CacheRoutes[0]},
			{CacheRoute: mf2.Spec.ServiceWorker.CacheRoutes[0], Prefix: buildPreCachePath(mf2, "")},
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

func TestServeProxyConfigIncludesAndOrdersServiceWorkerInterceptors(t *testing.T) {
	microFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()

	mfc := createTestMicroFrontendClass("polyfea", "/someother")
	mfc.Spec.ProgressiveWebApp = &v1alpha1.ProgressiveWebApp{
		ServiceWorker: &v1alpha1.ServiceWorker{
			Interceptors: []v1alpha1.SWInterceptor{
				{
					Name:      &[]string{"class-high"}[0],
					ModuleUrl: "/class-high.mjs",
					Priority:  &[]int32{10}[0],
				},
				{
					Name:      &[]string{"class-low"}[0],
					ModuleUrl: "/class-low.mjs",
				},
			},
		},
	}

	mf := createTestMicroFrontend("polyfea-interceptors", []string{}, "polyfea", true)
	mf.Spec.ServiceWorker = &v1alpha1.MicroFrontendServiceWorker{
		Interceptors: []v1alpha1.SWInterceptor{
			{
				Name:      &[]string{"mf-mid"}[0],
				ModuleUrl: "/mf-mid.mjs",
				Priority:  &[]int32{5}[0],
			},
		},
	}

	err := microFrontendRepository.Store(mf)
	if err != nil {
		t.Fatal(err)
	}

	pwa := NewProgressiveWebApplication(&logr.Logger{}, microFrontendRepository)

	actual, err := pwa.getProxyConfig(mfc)
	if err != nil {
		t.Fatal(err)
	}

	if len(actual.Interceptors) != 3 {
		t.Fatalf("Expected %d interceptors, got %d", 3, len(actual.Interceptors))
	}

	if *actual.Interceptors[0].Name != "class-high" {
		t.Errorf("Expected first interceptor %s, got %s", "class-high", *actual.Interceptors[0].Name)
	}
	if *actual.Interceptors[1].Name != "mf-mid" {
		t.Errorf("Expected second interceptor %s, got %s", "mf-mid", *actual.Interceptors[1].Name)
	}
	if *actual.Interceptors[2].Name != "class-low" {
		t.Errorf("Expected third interceptor %s, got %s", "class-low", *actual.Interceptors[2].Name)
	}

	expectedMfModuleURL := buildPreCachePath(mf, "/mf-mid.mjs")
	if actual.Interceptors[1].ModuleUrl != *expectedMfModuleURL {
		t.Errorf("Expected module URL %s, got %s", *expectedMfModuleURL, actual.Interceptors[1].ModuleUrl)
	}

	for i, interceptor := range actual.Interceptors {
		if interceptor.Priority != nil {
			t.Errorf("Expected interceptor at index %d to have nil priority", i)
		}
	}
}

func TestServeProxyConfigSkipsMicrofrontendInterceptorsWhenDisabled(t *testing.T) {
	microFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()

	mfc := createTestMicroFrontendClass("polyfea", "/someother")
	mfc.Spec.ProgressiveWebApp = &v1alpha1.ProgressiveWebApp{
		ServiceWorker: &v1alpha1.ServiceWorker{
			NoMicroFrontEndInterceptors: &[]bool{true}[0],
			Interceptors: []v1alpha1.SWInterceptor{
				{
					Name:      &[]string{"class-only"}[0],
					ModuleUrl: "/class-only.mjs",
				},
			},
		},
	}

	mf := createTestMicroFrontend("polyfea-interceptors", []string{}, "polyfea", true)
	mf.Spec.ServiceWorker = &v1alpha1.MicroFrontendServiceWorker{
		Interceptors: []v1alpha1.SWInterceptor{
			{
				Name:      &[]string{"mf-ignored"}[0],
				ModuleUrl: "/mf-ignored.mjs",
			},
		},
	}

	err := microFrontendRepository.Store(mf)
	if err != nil {
		t.Fatal(err)
	}

	pwa := NewProgressiveWebApplication(&logr.Logger{}, microFrontendRepository)

	actual, err := pwa.getProxyConfig(mfc)
	if err != nil {
		t.Fatal(err)
	}

	if len(actual.Interceptors) != 1 {
		t.Fatalf("Expected %d interceptor, got %d", 1, len(actual.Interceptors))
	}

	if *actual.Interceptors[0].Name != "class-only" {
		t.Errorf("Expected interceptor %s, got %s", "class-only", *actual.Interceptors[0].Name)
	}
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
	// Sort both PreCache slices by URL to ensure consistent comparison
	sort.Slice(expected.PreCache, func(i, j int) bool {
		return *expected.PreCache[i].URL < *expected.PreCache[j].URL
	})
	sort.Slice(actual.PreCache, func(i, j int) bool {
		return *actual.PreCache[i].URL < *actual.PreCache[j].URL
	})

	if len(actual.PreCache) != len(expected.PreCache) {
		t.Errorf("Expected %d, got %d", len(expected.PreCache), len(actual.PreCache))
	}
	for i, entry := range actual.PreCache {
		if *entry.URL != *expected.PreCache[i].URL {
			t.Errorf("Expected %s, got %s", *expected.PreCache[i].URL, *entry.URL)
		}
	}

	// Sort both Routes slices by Pattern to ensure consistent comparison
	sort.Slice(expected.Routes, func(i, j int) bool {
		return *expected.Routes[i].Pattern < *expected.Routes[j].Pattern
	})
	sort.Slice(actual.Routes, func(i, j int) bool {
		return *actual.Routes[i].Pattern < *actual.Routes[j].Pattern
	})

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
