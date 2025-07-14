package polyfea

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
	"github.com/polyfea/polyfea-controller/web-server/internal/polyfea/generated"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Improved readability and fixed issues in middlewaresTestSuite
var middlewaresTestSuite = IntegrationTestSuite{
	TestRouter: basePathStrippingMiddlewareRouter(),
	TestSet: []Test{
		{
			Name: "BasePathStrippingMiddlewareStripsTheBasePathAndForwardItInContext",
			Func: BasePathStrippingMiddlewareStripsTheBasePathAndForwardItInContext,
		},
		{
			Name: "BasePathStrippingMiddlewareNoBasePathIsFoundUseDefault",
			Func: BasePathStrippingMiddlewareNoBasePathIsFoundUseDefault,
		},
		{
			Name: "BasePathStrippingMiddlewareBasePathIsFound",
			Func: BasePathStrippingMiddlewareBasePathIsFound,
		},
		{
			Name: "BasePathStrippingMiddlewareWithoutPolyfea",
			Func: BasePathStrippingMiddlewareWithoutPolyfea,
		},
		{
			Name: "BasePathStrippingMiddlewareRewritesURLPathCorrectly",
			Func: BasePathStrippingMiddlewareRewritesURLPathCorrectly,
		},
	},
}

// Added comments to clarify the purpose of each test
func BasePathStrippingMiddlewareStripsTheBasePathAndForwardItInContext(t *testing.T) {
	// Test that the middleware correctly strips the base path and forwards it in the context
	// Arrange
	testServerUrl := os.Getenv(TestServerUrlName)
	// Act
	response, err := http.Get(testServerUrl + "/someBasePath/polyfea")
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	if response.Header.Get("X-Base-Path") != "/someBasePath/" {
		t.Errorf("Expected header X-Base-Path to be %s, got %s", "/someBasePath/", response.Header.Get("X-Base-Path"))
	}
}

func BasePathStrippingMiddlewareNoBasePathIsFoundUseDefault(t *testing.T) {
	// Test that the middleware uses the default base path when no base path is found
	// Arrange
	testServerUrl := os.Getenv(TestServerUrlName)
	// Act
	response, err := http.Get(testServerUrl + "/polyfea")
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	if response.Header.Get("X-Base-Path") != "/" {
		t.Errorf("Expected header X-Base-Path to be %s, got %s", "/", response.Header.Get("X-Base-Path"))
	}
}

func BasePathStrippingMiddlewareBasePathIsFound(t *testing.T) {
	// Test that the middleware correctly identifies and forwards the base path
	// Arrange
	testServerUrl := os.Getenv(TestServerUrlName)
	expectedBasePath := "/expected/base/path/"
	// Act
	response, err := http.Get(testServerUrl + expectedBasePath + "polyfea/some/path")
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	if response.Header.Get("X-Base-Path") != expectedBasePath {
		t.Errorf("Expected header X-Base-Path to be %s, got %s", expectedBasePath, response.Header.Get("X-Base-Path"))
	}
}

func BasePathStrippingMiddlewareWithoutPolyfea(t *testing.T) {
	// Test that the middleware handles paths without "polyfea" correctly
	// Arrange
	testServerUrl := os.Getenv(TestServerUrlName)
	// Act
	response, err := http.Get(testServerUrl + "/some/other/path")
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Assuming the middleware sets a default base path or doesn't modify it
	expectedBasePath := "/"
	if response.Header.Get("X-Base-Path") != expectedBasePath {
		t.Errorf("Expected header X-Base-Path to be %s, got %s", expectedBasePath, response.Header.Get("X-Base-Path"))
	}
}

func BasePathStrippingMiddlewareRewritesURLPathCorrectly(t *testing.T) {
	// Test that the middleware correctly rewrites the URL path
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/polyfea/some/path" {
			t.Errorf("Expected URL path to be rewritten to %s, got %s", "/polyfea/some/path", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})

	mfcRepository := setupMfcRepository()
	middlewareHandler := BasePathStrippingMiddleware(handler, mfcRepository)
	testServer := httptest.NewServer(middlewareHandler)
	defer testServer.Close()

	// Act
	response, err := http.Get(testServer.URL + "/any/base/path/polyfea/some/path")
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func basePathStrippingMiddlewareRouter() http.Handler {
	router := generated.NewRouter()
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		basePathValue := r.Context().Value(PolyfeaContextKeyBasePath)
		w.Header().Set("X-Base-Path", basePathValue.(string))
		w.WriteHeader(http.StatusOK)
	})
	return BasePathStrippingMiddleware(router, setupMfcRepository())
}

func setupMfcRepository() repository.Repository[*v1alpha1.MicroFrontendClass] {
	testData := []struct {
		name      string
		namespace string
		baseUri   *string
	}{
		{"default", "default", ptr("/")},
		{"some", "somes", ptr("/someBasePath/")},
		{"expected", "expecteds", ptr("/expected/base/path")},
		{"fea", "feas", ptr("/fea/")},
		{"feature", "features", ptr("/feature/")},
	}

	mfcRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontendClass]()

	for _, data := range testData {
		mfcRepository.Store(&v1alpha1.MicroFrontendClass{
			ObjectMeta: metav1.ObjectMeta{
				Name:      data.name,
				Namespace: data.namespace,
			},
			Spec: v1alpha1.MicroFrontendClassSpec{
				BaseUri: data.baseUri,
			},
		})
	}

	return mfcRepository
}
