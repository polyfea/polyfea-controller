package polyfea

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
	"github.com/polyfea/polyfea-controller/web-server/api"
)

var (
	testServer *httptest.Server
)

func TestMain(t *testing.M) {
	// Setup
	r := setupRouter()

	// Test server
	testServer = httptest.NewServer(r)
	defer testServer.Close()

	t.Run()
}

func TestPolyfeaApiGetContextAreaReturnsNotImplemented(t *testing.T) {
	// Arrange
	testServerUrl := testServer.URL
	testRoute := "/polyfea/context-area/test?path=test&take=0"

	// Act
	response, err := http.Get(testServerUrl + testRoute)
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNotImplemented {
		t.Errorf("Expected status code %d, got %d", http.StatusNotImplemented, response.StatusCode)
	}
}

func TestPolyfeaApiGetStaticConfigReturnsNotImplemented(t *testing.T) {
	// Arrange
	testServerUrl := testServer.URL
	testRoute := "/polyfea/static-config"

	// Act
	response, err := http.Get(testServerUrl + testRoute)
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNotImplemented {
		t.Errorf("Expected status code %d, got %d", http.StatusNotImplemented, response.StatusCode)
	}
}

func setupRouter() *mux.Router {
	testWebComponentRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.WebComponent]()
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	testMicroFrontedClassRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontendClass]()

	polyfeaAPIService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		testMicroFrontedClassRepository,
		map[string]string{})

	polyfeaAPIController := NewPolyfeaAPIController(polyfeaAPIService)

	router := NewRouter(polyfeaAPIController)

	router.HandleFunc("/openapi", api.HandleOpenApi)

	return router
}
