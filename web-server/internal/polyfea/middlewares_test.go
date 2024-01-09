package polyfea

import (
	"net/http"
	"os"
	"testing"

	"github.com/polyfea/polyfea-controller/web-server/internal/polyfea/generated"
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

func basePathStrippingMiddlewareRouter() http.Handler {
	router := generated.NewRouter()
	router.HandleFunc("/polyfea", func(w http.ResponseWriter, r *http.Request) {
		basePathValue := r.Context().Value(PolyfeaContextKeyBasePath)
		w.Header().Set("X-Base-Path", basePathValue.(string))
		w.WriteHeader(http.StatusOK)
	})
	return BasePathStrippingMiddleware(router)
}
