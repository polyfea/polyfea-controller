package polyfea

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sort"
	"testing"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
	"github.com/polyfea/polyfea-controller/web-server/internal/polyfea/generated"
	"github.com/rs/zerolog"
)

var pwaTestSuite = IntegrationTestSuite{
	TestRouter: polyfeaPWAApiSetupRouter(),
	TestSet: []Test{
		{
			Name: "ServeAppWebManifestReturnsExpectedManifest",
			Func: ServeAppWebManifestReturnsExpectedManifest,
		},
		{
			Name: "ServeRegisterReturnsExpectedFile",
			Func: ServeRegisterReturnsExpectedFile,
		},
		{
			Name: "ServeServiceWorkerReturnsExpectedFile",
			Func: ServeServiceWorkerReturnsExpectedFile,
		},
		{
			Name: "ServeCachingReturnsExpectedConfig",
			Func: ServeCachingReturnsExpectedConfig,
		},
	},
}

func ServeAppWebManifestReturnsExpectedManifest(t *testing.T) {
	testServerUrl := os.Getenv(TestServerUrlName)
	expected := createExpectedWebAppManifest()

	response, err := http.Get(testServerUrl + "/polyfea/app.webmanifest")
	handleHTTPError(t, err)
	defer response.Body.Close()

	actual := parseResponseBodyToManifest(t, response.Body)
	assertWebAppManifestEquality(t, expected, actual)
}

func ServeRegisterReturnsExpectedFile(t *testing.T) {
	testServerUrl := os.Getenv(TestServerUrlName)
	expected := readFileContent(t, ".resources/register.mjs")

	response, err := http.Get(testServerUrl + "/polyfea/register.mjs")
	handleHTTPError(t, err)
	defer response.Body.Close()

	body := readResponseBody(t, response.Body)
	assertStringEquality(t, expected, string(body))
}

func ServeServiceWorkerReturnsExpectedFile(t *testing.T) {
	testServerUrl := os.Getenv(TestServerUrlName)
	expected := readFileContent(t, ".resources/sw.mjs")

	response, err := http.Get(testServerUrl + "/sw.mjs")
	handleHTTPError(t, err)
	defer response.Body.Close()

	body := readResponseBody(t, response.Body)
	assertStringEquality(t, expected, string(body))
}

func ServeCachingReturnsExpectedConfig(t *testing.T) {
	testServerUrl := os.Getenv(TestServerUrlName)
	expected := createExpectedProxyConfigResponse()

	response, err := http.Get(testServerUrl + "/polyfea-caching.json")
	handleHTTPError(t, err)
	defer response.Body.Close()

	actual := parseResponseBodyToProxyConfig(t, response.Body)
	sortProxyConfigEntries(actual, expected)
	assertProxyConfigEquality(t, expected, actual)
}

// Helper functions
func handleHTTPError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func readFileContent(t *testing.T, path string) string {
	file, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(file)
}

func readResponseBody(t *testing.T, body io.ReadCloser) []byte {
	content, err := io.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}
	return content
}

func parseResponseBodyToManifest(t *testing.T, body io.ReadCloser) *v1alpha1.WebAppManifest {
	content := readResponseBody(t, body)
	manifest := &v1alpha1.WebAppManifest{}
	err := json.Unmarshal(content, manifest)
	if err != nil {
		t.Fatal(err)
	}
	return manifest
}

func parseResponseBodyToProxyConfig(t *testing.T, body io.ReadCloser) *ProxyConfigResponse {
	content := readResponseBody(t, body)
	config := &ProxyConfigResponse{}
	err := json.Unmarshal(content, config)
	if err != nil {
		t.Fatal(err)
	}
	return config
}

func assertStringEquality(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func assertWebAppManifestEquality(t *testing.T, expected, actual *v1alpha1.WebAppManifest) {
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

func assertProxyConfigEquality(t *testing.T, expected, actual *ProxyConfigResponse) {
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

func sortProxyConfigEntries(actual, expected *ProxyConfigResponse) {
	sort.Slice(actual.PreCache, func(i, j int) bool {
		return *actual.PreCache[i].URL < *actual.PreCache[j].URL
	})
	sort.Slice(expected.PreCache, func(i, j int) bool {
		return *expected.PreCache[i].URL < *expected.PreCache[j].URL
	})

	sort.Slice(actual.Routes, func(i, j int) bool {
		if actual.Routes[i].Pattern == nil {
			return true
		}
		if actual.Routes[j].Pattern == nil {
			return false
		}
		return *actual.Routes[i].Pattern < *actual.Routes[j].Pattern
	})
	sort.Slice(expected.Routes, func(i, j int) bool {
		if actual.Routes[i].Pattern == nil {
			return true
		}
		if actual.Routes[j].Pattern == nil {
			return false
		}
		return *expected.Routes[i].Pattern < *expected.Routes[j].Pattern
	})
}

func createExpectedWebAppManifest() *v1alpha1.WebAppManifest {
	return &v1alpha1.WebAppManifest{
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
}

func createExpectedProxyConfigResponse() *ProxyConfigResponse {
	mfc := createTestMicroFrontendClass("polyfea", "/some")
	mfc.Spec.ProgressiveWebApp = &v1alpha1.ProgressiveWebApp{
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
	}

	mf := createTestMicroFrontend("polyfea1", []string{}, "test-module", "polyfea", true)
	mf.Spec.CacheOptions = &v1alpha1.PWACache{
		PreCache: []v1alpha1.PreCacheEntry{
			{
				URL: &[]string{"/test"}[0],
			},
		},
	}

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

	return &ProxyConfigResponse{
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
				Prefix:     buildPreCachePath(mf2, ""),
			},
		},
	}
}

func polyfeaPWAApiSetupRouter() http.Handler {
	mfc := createTestMicroFrontendClass("polyfea", "/some")
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
	}

	router := generated.NewRouter()

	microFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()

	mf := createTestMicroFrontend("polyfea1", []string{}, "test-module", "polyfea", true)
	mf.Spec.CacheOptions = &v1alpha1.PWACache{
		PreCache: []v1alpha1.PreCacheEntry{
			{
				URL: &[]string{"/test"}[0],
			},
		},
	}
	microFrontendRepository.Store(mf)

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
	microFrontendRepository.Store(mf2)

	mf3 := createTestMicroFrontend("polyfea3", []string{}, "test-module", "polyfea", false)
	mf3.Spec.CacheOptions = &v1alpha1.PWACache{
		PreCache: []v1alpha1.PreCacheEntry{
			{
				URL: &[]string{"/test3"}[0],
			},
		},
	}
	microFrontendRepository.Store(mf3)

	mf4 := createTestMicroFrontend("polyfea4", []string{}, "test-module", "someother", false)
	mf4.Spec.CacheOptions = &v1alpha1.PWACache{
		PreCache: []v1alpha1.PreCacheEntry{
			{
				URL: &[]string{"/test4"}[0],
			},
		},
	}
	microFrontendRepository.Store(mf4)

	spa := NewProgressiveWebApplication(&zerolog.Logger{}, microFrontendRepository)

	router.HandleFunc("/polyfea/app.webmanifest", spa.ServeAppWebManifest)
	router.HandleFunc("/polyfea/register.mjs", spa.ServeRegister)
	router.HandleFunc("/sw.mjs", spa.ServeServiceWorker)
	router.HandleFunc("/polyfea-caching.json", spa.ServeCaching)

	return addDummyMiddleware(router, "/", mfc)
}
