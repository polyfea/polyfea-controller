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

func BasePathStrippingMiddlewareStripsTheBasePathAndForwardItInContext(t *testing.T) {
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
		t.Errorf("Expected header X-Base-Path to be %s, got %s", "someBasePath", response.Header.Get("X-Base-Path"))
	}
}

func BasePathStrippingMiddlewareNoBasePathIsFoundUseDefault(t *testing.T) {
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
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/polyfea/some/path" {
			t.Errorf("Expected URL path to be rewritten to %s, got %s", "/polyfea/some/path", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the handler with the BasePathStrippingMiddleware
	// Assuming setupMfcRepository() returns a configured repository needed by the middleware
	mfcRepository := setupMfcRepository() // This function needs to be defined according to your application's requirements
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
	mfcRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontendClass]()
	defaultBase := "/"
	mfcRepository.StoreItem(&v1alpha1.MicroFrontendClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "default",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			BaseUri: &defaultBase,
		},
	})

	someBasePath := "/someBasePath/"
	mfcRepository.StoreItem(&v1alpha1.MicroFrontendClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "some",
			Namespace: "somes",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			BaseUri: &someBasePath,
		},
	})

	expectedBasePath := "/expected/base/path"
	mfcRepository.StoreItem(&v1alpha1.MicroFrontendClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "expected",
			Namespace: "expecteds",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			BaseUri: &expectedBasePath,
		},
	})

	feaBase := "/fea/"
	mfcRepository.StoreItem(&v1alpha1.MicroFrontendClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fea",
			Namespace: "feas",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			BaseUri: &feaBase,
		},
	})

	featureBase := "/feature/"
	mfcRepository.StoreItem(&v1alpha1.MicroFrontendClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "feature",
			Namespace: "features",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			BaseUri: &featureBase,
		},
	})

	return mfcRepository
}
