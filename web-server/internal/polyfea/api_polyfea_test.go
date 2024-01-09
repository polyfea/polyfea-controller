package polyfea

import (
	"net/http"
	"os"
	"testing"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
	"github.com/polyfea/polyfea-controller/web-server/api"
)

var apiPolyfeaTestSuite = IntegrationTestSuite{
	TestRouter: polyfeaApiSetupRouter(),
	TestSet: []Test{
		{
			Name: "PolyfeaApiGetContextAreaReturnsNotImplemented",
			Func: PolyfeaApiGetContextAreaReturnsNotImplemented,
		},
		{
			Name: "PolyfeaApiGetStaticConfigReturnsNotImplemented",
			Func: PolyfeaApiGetStaticConfigReturnsNotImplemented,
		},
	},
}

func PolyfeaApiGetContextAreaReturnsNotImplemented(t *testing.T) {
	// Arrange
	testServerUrl := os.Getenv(TestServerUrlName)
	testRoute := "/polyfea/context-area/test?path=test&take=0"

	// Act
	response, err := http.Get(testServerUrl + testRoute)
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.StatusCode)
	}
}

func PolyfeaApiGetStaticConfigReturnsNotImplemented(t *testing.T) {
	// Arrange
	testServerUrl := os.Getenv(TestServerUrlName)
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

func polyfeaApiSetupRouter() http.Handler {
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
