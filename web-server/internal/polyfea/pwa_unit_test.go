package polyfea

import (
	"testing"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
	"github.com/rs/zerolog"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestServeAppWebManifestReturnsExpectedManifest(t *testing.T) {
	// Arrange
	pwa := NewProgressiveWebApplication(&zerolog.Logger{}, repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]())

	mfc := &v1alpha1.MicroFrontendClass{
		ObjectMeta: v1.ObjectMeta{
			Name:      "polyfea",
			Namespace: "polyfea",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			Title:     &[]string{"Test MicroFrontendClass"}[0],
			BaseUri:   &[]string{"/someother"}[0],
			CspHeader: "default-src 'self';",
			ExtraHeaders: []v1alpha1.Header{
				{
					Name:  "X-Frame-Options",
					Value: "DENY",
				},
			},
			UserRolesHeader: "X-User-Roles",
			UserHeader:      "X-User-Id",
			ProgressiveWebApp: &v1alpha1.ProgressiveWebApp{
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
			},
		},
	}

	expected := &v1alpha1.WebAppManifest{
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
	}

	// Act
	actual := pwa.serveAppWebManifest(mfc)

	// Assert
	if *actual.Name != *expected.Name {
		t.Errorf("Expected %s, got %s", *expected.Name, *actual.Name)
	}
	if *actual.Icons[0].Src != *expected.Icons[0].Src {
		t.Errorf("Expected %s, got %s", *expected.Icons[0].Src, *actual.Icons[0].Src)
	}
	if *actual.Icons[0].Sizes != *expected.Icons[0].Sizes {
		t.Errorf("Expected %s, got %s", *expected.Icons[0].Sizes, *actual.Icons[0].Sizes)
	}
	if *actual.Icons[0].Type != *expected.Icons[0].Type {
		t.Errorf("Expected %s, got %s", *expected.Icons[0].Type, *actual.Icons[0].Type)
	}
	if *actual.StartUrl != *expected.StartUrl {
		t.Errorf("Expected %s, got %s", *expected.StartUrl, *actual.StartUrl)
	}
	if *actual.Display != *expected.Display {
		t.Errorf("Expected %s, got %s", *expected.Display, *actual.Display)
	}
}

func TestServeProxyConfigReturnsExpectedConfigForAllRelevantMicrofrontends(t *testing.T) {
	// Arrange
	microFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()

	mfc := &v1alpha1.MicroFrontendClass{
		ObjectMeta: v1.ObjectMeta{
			Name:      "polyfea",
			Namespace: "default",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			Title:     &[]string{"Test MicroFrontendClass"}[0],
			BaseUri:   &[]string{"/someother"}[0],
			CspHeader: "default-src 'self';",
			ExtraHeaders: []v1alpha1.Header{
				{
					Name:  "X-Frame-Options",
					Value: "DENY",
				},
			},
			UserRolesHeader: "X-User-Roles",
			UserHeader:      "X-User-Id",
			ProgressiveWebApp: &v1alpha1.ProgressiveWebApp{
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
				CacheOptions: &v1alpha1.PWACache{
					PreCache: []v1alpha1.PreCacheEntry{
						{
							URL: &[]string{"/test-class"}[0],
						},
					},
					CacheRoutes: []v1alpha1.CacheRoute{
						{
							Pattern:  &[]string{"/cache-route-class"}[0],
							Strategy: &[]string{"cache-first"}[0],
						},
					},
				},
			},
		},
	}

	mf := createTestMicroFrontend("polyfea1", []string{}, "test-module", "polyfea", true)
	mf.Spec.CacheOptions = &v1alpha1.PWACache{
		PreCache: []v1alpha1.PreCacheEntry{
			{
				URL: &[]string{"/test"}[0],
			},
		},
	}
	microFrontendRepository.StoreItem(mf)

	mf2 := createTestMicroFrontend("polyfea2", []string{}, "test-module", "polyfea", true)
	mf2.Spec.CacheOptions = &v1alpha1.PWACache{
		PreCache: []v1alpha1.PreCacheEntry{
			{
				URL: &[]string{"/test2"}[0],
			},
		},
		CacheRoutes: []v1alpha1.CacheRoute{
			{
				Pattern:  &[]string{"/cache-route"}[0],
				Strategy: &[]string{"network-first"}[0],
			},
		},
	}
	microFrontendRepository.StoreItem(mf2)

	mf3 := createTestMicroFrontend("polyfea3", []string{}, "test-module", "polyfea", false)
	mf3.Spec.CacheOptions = &v1alpha1.PWACache{
		PreCache: []v1alpha1.PreCacheEntry{
			{
				URL: &[]string{"/test3"}[0],
			},
		},
	}
	microFrontendRepository.StoreItem(mf3)

	mf4 := createTestMicroFrontend("polyfea4", []string{}, "test-module", "someother", false)
	mf4.Spec.CacheOptions = &v1alpha1.PWACache{
		PreCache: []v1alpha1.PreCacheEntry{
			{
				URL: &[]string{"/test4"}[0],
			},
		},
	}
	microFrontendRepository.StoreItem(mf4)

	pwa := NewProgressiveWebApplication(&zerolog.Logger{}, microFrontendRepository)

	expected := &ProxyConfigResponse{
		PreCache: []v1alpha1.PreCacheEntry{
			{
				URL: &[]string{"/test-class"}[0],
			},
			{
				URL: buildPreCachePath(mf, *mf.Spec.CacheOptions.PreCache[0].URL),
			},
			{
				URL: buildPreCachePath(mf2, *mf2.Spec.CacheOptions.PreCache[0].URL),
			},
		},
		Routes: []CacheRouteResponse{
			{
				CacheRoute: mfc.Spec.ProgressiveWebApp.CacheOptions.CacheRoutes[0],
			},
			{
				CacheRoute: mf2.Spec.CacheOptions.CacheRoutes[0],
				Prefix:     buildPreCachePath(mf2, ""), // TODO: Should there be a leading / here?
			},
		},
	}

	// Act
	actual, err := pwa.getProxyConfig(mfc)

	if err != nil {
		t.Fatal(err)
	}

	// Assert
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
