package polyfea

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
	"github.com/polyfea/polyfea-controller/web-server/api"
	"github.com/polyfea/polyfea-controller/web-server/internal/polyfea/generated"
	"github.com/rs/zerolog"
)

var apiPolyfeaTestSuite = IntegrationTestSuite{
	TestRouter: polyfeaApiSetupRouter(),
	TestSet: []Test{
		{
			Name: "PolyfeaApiGetContextAreaMultipleElementsTakeOneCorrectComponentIsSelected",
			Func: PolyfeaApiGetContextAreaMultipleElementsTakeOneCorrectComponentIsSelected,
		},
		{
			Name: "PolyfeaApiGetContextAreaMultipleElementsNotMatchingReturnNotFound",
			Func: PolyfeaApiGetContextAreaMultipleElementsNotMatchingReturnNotFound,
		},
		{
			Name: "PolyfeaApiGetContextAreaMultipleElementsMatchingReturned",
			Func: PolyfeaApiGetContextAreaMultipleElementsMatchingReturned,
		},
		{
			Name: "PolyfeaApiGetStaticConfigReturnsNotImplemented",
			Func: PolyfeaApiGetStaticConfigReturnsNotImplemented,
		},
	},
}

func PolyfeaApiGetContextAreaMultipleElementsTakeOneCorrectComponentIsSelected(t *testing.T) {
	// Arrange
	testServerUrl := os.Getenv(TestServerUrlName)
	testRoute := "/polyfea/context-area/test-name?path=test-path&take=1"

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("other-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"other-microfrontend": createTestMicroFrontendSpec("other-microfrontend", []string{}, true),
		})

	// Act
	req, err := http.NewRequest("GET", testServerUrl+testRoute, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	req.Header.Set("test-user-roles-header", "some-different-role")
	req.Header.Add("test-user-roles-header", "test-role, test-other-role")

	client := &http.Client{}
	response, err := client.Do(req)
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer response.Body.Close()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	var actualContextArea generated.ContextArea
	err = json.NewDecoder(response.Body).Decode(&actualContextArea)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}

	if response.Header.Get("test-header") != "test-value" {
		t.Errorf("Expected %v, got %v", "test-value", response.Header.Get("test-header"))
	}
}

func PolyfeaApiGetContextAreaMultipleElementsNotMatchingReturnNotFound(t *testing.T) {
	// Arrange
	testServerUrl := os.Getenv(TestServerUrlName)
	testRoute := "/polyfea/context-area/tt-name?path=test.*&take=1"

	// Act
	req, err := http.NewRequest("GET", testServerUrl+testRoute, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	req.Header.Set("test-user-roles-header", "some-different-role")
	req.Header.Add("test-user-roles-header", "test-role, test-other-role")

	client := &http.Client{}
	response, err := client.Do(req)
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer response.Body.Close()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	var actualContextArea generated.ContextArea
	err = json.NewDecoder(response.Body).Decode(&actualContextArea)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	if len(actualContextArea.Elements) != 0 {
		t.Errorf("Expected no elements, got %v", actualContextArea.Elements)
	}
}

func PolyfeaApiGetContextAreaMultipleElementsMatchingReturned(t *testing.T) {
	// Arrange
	testServerUrl := os.Getenv(TestServerUrlName)
	testRoute := "/polyfea/context-area/test-name?path=test-path"

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("other-microfrontend"),
			createTestElementSpec("test-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"other-microfrontend": createTestMicroFrontendSpec("other-microfrontend", []string{}, true),
			"test-microfrontend":  createTestMicroFrontendSpec("test-microfrontend", []string{}, true),
		})

	// Act
	req, err := http.NewRequest("GET", testServerUrl+testRoute, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	req.Header.Set("test-user-roles-header", "some-different-role")
	req.Header.Add("test-user-roles-header", "test-role, test-other-role")

	client := &http.Client{}
	response, err := client.Do(req)
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer response.Body.Close()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	var actualContextArea generated.ContextArea
	err = json.NewDecoder(response.Body).Decode(&actualContextArea)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
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

	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		"test-tag-name",
		[]v1alpha1.DisplayRules{
			{
				NoneOf: []v1alpha1.Matcher{
					{
						Path: "tt.*",
					},
					{
						ContextName: "tt-name",
					},
				},
			},
		},
		&[]int32{1}[0]))

	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-other-name",
		"other-microfrontend",
		"test-tag-name",
		[]v1alpha1.DisplayRules{
			{
				NoneOf: []v1alpha1.Matcher{
					{
						Path: "tes-path",
					},
					{
						ContextName: "tt-name",
					},
				},
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test.*",
					},
					{
						ContextName: "test-name",
					},
				},
				AnyOf: []v1alpha1.Matcher{
					{
						Path: "t-path",
					},
					{
						Role: "test-role",
					},
				},
			},
		},
		&[]int32{10}[0]))

	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()

	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true))
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("other-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	testMicroFrontendClassRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontendClass]()

	mfc := createTestMicroFrontendClass("test-frontend-class", "/")
	mfc.Spec.ExtraHeaders = []v1alpha1.Header{
		{
			Name:  "test-header",
			Value: "test-value",
		},
	}

	testMicroFrontendClassRepository.StoreItem(mfc)
	testMicroFrontendClassRepository.StoreItem(createTestMicroFrontendClass("other-frontend-class", "other"))

	polyfeaAPIService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		testMicroFrontendClassRepository,
		&zerolog.Logger{},
	)

	polyfeaAPIController := generated.NewPolyfeaAPIController(polyfeaAPIService)

	router := generated.NewRouter(polyfeaAPIController)

	router.HandleFunc("/openapi", api.HandleOpenApi)

	return router
}
