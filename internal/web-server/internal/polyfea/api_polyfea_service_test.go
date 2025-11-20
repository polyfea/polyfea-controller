package polyfea

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-logr/logr"
	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/internal/repository"
	"github.com/polyfea/polyfea-controller/internal/web-server/internal/polyfea/generated"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestPolyfeaApiServiceGetContextAreaReturnsContextAreaIfRepositoryContainsMatchingWebComponents(t *testing.T) {
	// Arrange
	webComponentRepo, microFrontendRepo := setupRepositories()
	err := webComponentRepo.Store(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test-path",
					},
					{
						ContextName: "test-name",
					},
				},
			},
		},
		&[]int32{1}[0]))

	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	err = microFrontendRepo.Store(createTestMicroFrontend("test-microfrontend", []string{}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	polyfeaApiService := NewPolyfeaAPIService(webComponentRepo, microFrontendRepo, &logr.Logger{})
	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("test-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"test-microfrontend": createTestMicroFrontendSpec("test-microfrontend", []string{}),
		})
	ctx := setupContext("/", "test-frontend-class")

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/context-area/test-name?path=test-path&take=10", nil)
	req = req.WithContext(ctx)
	take := 10
	polyfeaApiService.GetContextArea(w, req, "test-name", generated.GetContextAreaParams{Path: "test-path", Take: &take})

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %v", w.Code)
	}
	var actualContextArea generated.ContextArea
	if err := json.Unmarshal(w.Body.Bytes(), &actualContextArea); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsContextAreaWithExtraHeaders(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryRepository[*v1alpha1.WebComponent]()
	err := testWebComponentRepository.Store(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test-path",
					},
					{
						ContextName: "test-name",
					},
				},
			},
		},
		&[]int32{1}[0]))

	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}
	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	err = testMicroFrontendRepository.Store(createTestMicroFrontend("test-microfrontend", []string{}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	mfc := createTestMicroFrontendClass("test-frontend-class", "/")
	mfc.Spec.ExtraHeaders = []v1alpha1.Header{
		{
			Name:  "test-header",
			Value: "test-value",
		},
	}

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&logr.Logger{},
	)

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("test-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"test-microfrontend": createTestMicroFrontendSpec("test-microfrontend", []string{}),
		})
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, mfc)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/context-area/test-name", nil)
	req = req.WithContext(ctx)
	take := 10
	polyfeaApiService.GetContextArea(w, req, "test-name", generated.GetContextAreaParams{Path: "test-path", Take: &take})

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %v", w.Code)
	}

	var actualContextArea generated.ContextArea
	if err := json.Unmarshal(w.Body.Bytes(), &actualContextArea); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}

	// Check extra headers
	if w.Header().Get("test-header") != "test-value" {
		t.Errorf("Expected header 'test-header' with value 'test-value', got '%v'", w.Header().Get("test-header"))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsContextAreaIfNoneOfIsMatching(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryRepository[*v1alpha1.WebComponent]()
	err := testWebComponentRepository.Store(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				NoneOf: []v1alpha1.Matcher{
					{
						Path: "teset-path",
					},
					{
						ContextName: "etest-name",
					},
				},
			},
		},
		&[]int32{1}[0]))
	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	err = testMicroFrontendRepository.Store(createTestMicroFrontend("test-microfrontend", []string{}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&logr.Logger{},
	)

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("test-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"test-microfrontend": createTestMicroFrontendSpec("test-microfrontend", []string{}),
		})
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	// Act
	take := 10
	statusCode, actualContextArea := callGetContextArea(t, polyfeaApiService, ctx, "test-name", "test-path", &take, nil)

	// Assert
	if statusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %v", statusCode)
	}

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsContextAreaIfAnyOfIsMatching(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryRepository[*v1alpha1.WebComponent]()
	err := testWebComponentRepository.Store(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				AnyOf: []v1alpha1.Matcher{
					{
						Path: "teet-path",
					},
					{
						ContextName: "test-name",
					},
				},
			},
		},
		&[]int32{1}[0]))
	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	err = testMicroFrontendRepository.Store(createTestMicroFrontend("test-microfrontend", []string{}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&logr.Logger{},
	)

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("test-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"test-microfrontend": createTestMicroFrontendSpec("test-microfrontend", []string{}),
		})
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	// Act
	take := 10
	statusCode, actualContextArea := callGetContextArea(t, polyfeaApiService, ctx, "test-name", "test-path", &take, nil)

	// Assert
	if statusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %v", statusCode)
	}

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsContextAreaIfComplexCombinationIsMatching(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryRepository[*v1alpha1.WebComponent]()
	err := testWebComponentRepository.Store(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				NoneOf: []v1alpha1.Matcher{
					{
						Path: "teset-path",
					},
					{
						ContextName: "etest-name",
					},
				},
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test-path",
					},
					{
						ContextName: "test-name",
					},
				},
				AnyOf: []v1alpha1.Matcher{
					{
						Path: "teet-path",
					},
					{
						ContextName: "test-name",
					},
				},
			},
		},
		&[]int32{1}[0]))
	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	err = testMicroFrontendRepository.Store(createTestMicroFrontend("test-microfrontend", []string{}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}
	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&logr.Logger{},
	)

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("test-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"test-microfrontend": createTestMicroFrontendSpec("test-microfrontend", []string{}),
		})
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	// Act
	take := 10
	statusCode, actualContextArea := callGetContextArea(t, polyfeaApiService, ctx, "test-name", "test-path", &take, nil)

	// Assert
	if statusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %v", statusCode)
	}

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsContextAreaIfComplexMatcherIsMatching(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryRepository[*v1alpha1.WebComponent]()
	err := testWebComponentRepository.Store(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test*",
					},
					{
						ContextName: "test-name",
					},
					{
						Role: "test-role",
					},
					{
						Role: "test-other-role",
					},
				},
			},
		},
		&[]int32{1}[0]))
	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	err = testMicroFrontendRepository.Store(createTestMicroFrontend("test-microfrontend", []string{}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&logr.Logger{},
	)

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("test-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"test-microfrontend": createTestMicroFrontendSpec("test-microfrontend", []string{}),
		})
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	headers := http.Header{}
	headers.Set("test-user-roles-header", "some-different-role")
	headers.Add("test-user-roles-header", "test-role, test-other-role")

	// Act
	take := 10
	statusCode, actualContextArea := callGetContextArea(t, polyfeaApiService, ctx, "test-name", "test-path", &take, headers)

	// Assert
	if statusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %v", statusCode)
	}

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsElementWithoutMicrofrontendIfItHasNoMicrofrontends(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryRepository[*v1alpha1.WebComponent]()
	err := testWebComponentRepository.Store(createTestWebComponent(
		"test-name",
		"",
		[]v1alpha1.DisplayRules{
			{
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test*",
					},
					{
						ContextName: "test-name",
					},
					{
						Role: "test-role",
					},
					{
						Role: "test-other-role",
					},
				},
			},
		},
		&[]int32{1}[0]))
	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&logr.Logger{},
	)

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec(""),
		},
		map[string]generated.MicrofrontendSpec{})
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	headers := http.Header{}
	headers.Set("test-user-roles-header", "some-different-role")
	headers.Add("test-user-roles-header", "test-role, test-other-role")

	// Act
	take := 10
	statusCode, actualContextArea := callGetContextArea(t, polyfeaApiService, ctx, "test-name", "test-path", &take, headers)

	// Assert
	if statusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %v", statusCode)
	}

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsEmptyIfRoleMatcherIsNotMatching(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryRepository[*v1alpha1.WebComponent]()
	err := testWebComponentRepository.Store(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test-path",
					},
					{
						ContextName: "test-name",
					},
					{
						Role: "test-role-not-matching",
					},
					{
						Role: "test-other-role",
					},
				},
			},
		},
		&[]int32{1}[0]))
	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	err = testMicroFrontendRepository.Store(createTestMicroFrontend("test-microfrontend", []string{}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&logr.Logger{},
	)

	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	headers := http.Header{}
	headers.Set("test-user-roles-header", "some-different-role")
	headers.Add("test-user-roles-header", "test-role, test-other-role")

	// Act
	take := 10
	statusCode, contextArea := callGetContextArea(t, polyfeaApiService, ctx, "test-name", "test.*", &take, headers)

	// Assert
	if statusCode != http.StatusOK {
		t.Errorf("Expected 200, got %v", statusCode)
	}

	if len(contextArea.Elements) != 0 {
		t.Errorf("Expected 0 elements, got %v", len(contextArea.Elements))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsEmptyIfContextMatcherIsNotMatching(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryRepository[*v1alpha1.WebComponent]()
	err := testWebComponentRepository.Store(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test-path",
					},
					{
						ContextName: "test-name-not-matching",
					},
					{
						Role: "test-role",
					},
					{
						Role: "test-other-role",
					},
				},
			},
		},
		&[]int32{1}[0]))
	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	err = testMicroFrontendRepository.Store(createTestMicroFrontend("test-microfrontend", []string{}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&logr.Logger{},
	)

	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	headers := http.Header{}
	headers.Set("test-user-roles-header", "some-different-role")
	headers.Add("test-user-roles-header", "test-role, test-other-role")

	// Act
	take := 10
	statusCode, contextArea := callGetContextArea(t, polyfeaApiService, ctx, "test-name", "test.*", &take, headers)

	// Assert
	if statusCode != http.StatusOK {
		t.Errorf("Expected 200, got %v", statusCode)
	}

	if len(contextArea.Elements) != 0 {
		t.Errorf("Expected 0 elements, got %v", len(contextArea.Elements))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsEmptyIfPathIsNotMatching(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryRepository[*v1alpha1.WebComponent]()
	err := testWebComponentRepository.Store(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				AllOf: []v1alpha1.Matcher{
					{
						Path: "soe.*",
					},
					{
						ContextName: "test-name",
					},
					{
						Role: "test-role",
					},
					{
						Role: "test-other-role",
					},
				},
			},
		},
		&[]int32{1}[0]))
	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	err = testMicroFrontendRepository.Store(createTestMicroFrontend("test-microfrontend", []string{}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&logr.Logger{},
	)

	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	headers := http.Header{}
	headers.Set("test-user-roles-header", "some-different-role")
	headers.Add("test-user-roles-header", "test-role, test-other-role")

	// Act
	take := 10
	statusCode, contextArea := callGetContextArea(t, polyfeaApiService, ctx, "test-name", "sometest", &take, headers)

	// Assert
	if statusCode != http.StatusOK {
		t.Errorf("Expected 200, got %v", statusCode)
	}

	if len(contextArea.Elements) != 0 {
		t.Errorf("Expected 0 elements, got %v", len(contextArea.Elements))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsEmptyIfNoneOfIsMatching(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryRepository[*v1alpha1.WebComponent]()
	err := testWebComponentRepository.Store(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test-path",
					},
					{
						ContextName: "test-name",
					},
					{
						Role: "test-role",
					},
					{
						Role: "test-other-role",
					},
				},
				NoneOf: []v1alpha1.Matcher{
					{
						Role: "test-role",
					},
				},
			},
		},
		&[]int32{1}[0]))
	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	err = testMicroFrontendRepository.Store(createTestMicroFrontend("test-microfrontend", []string{}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&logr.Logger{},
	)

	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	headers := http.Header{}
	headers.Set("test-user-roles-header", "some-different-role")
	headers.Add("test-user-roles-header", "test-role, test-other-role")

	// Act
	take := 10
	statusCode, contextArea := callGetContextArea(t, polyfeaApiService, ctx, "test-name", "sometest.*", &take, headers)

	// Assert
	if statusCode != http.StatusOK {
		t.Errorf("Expected 200, got %v", statusCode)
	}

	if len(contextArea.Elements) != 0 {
		t.Errorf("Expected 0 elements, got %v", len(contextArea.Elements))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsEmptyIfAnyOfIsNotMatching(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryRepository[*v1alpha1.WebComponent]()
	err := testWebComponentRepository.Store(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				AnyOf: []v1alpha1.Matcher{
					{
						Path: "tes-path",
					},
					{
						ContextName: "test-nameer",
					},
					{
						Role: "rtest-role",
					},
					{
						Role: "wtest-other-role",
					},
				},
			},
		},
		&[]int32{1}[0]))
	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	err = testMicroFrontendRepository.Store(createTestMicroFrontend("test-microfrontend", []string{}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&logr.Logger{},
	)

	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	headers := http.Header{}
	headers.Set("test-user-roles-header", "some-different-role")
	headers.Add("test-user-roles-header", "test-role, test-other-role")

	// Act
	take := 10
	statusCode, contextArea := callGetContextArea(t, polyfeaApiService, ctx, "test-name", "sometest.*", &take, headers)

	// Assert
	if statusCode != http.StatusOK {
		t.Errorf("Expected 200, got %v", statusCode)
	}

	if len(contextArea.Elements) != 0 {
		t.Errorf("Expected 0 elements, got %v", len(contextArea.Elements))
	}
}

func TestPolyfeaApiServiceGetContextAreaMultipleElementsTakeOneOnlyOneElementReturned(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryRepository[*v1alpha1.WebComponent]()

	err := testWebComponentRepository.Store(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test-path",
					},
					{
						ContextName: "test-name",
					},
				},
			},
		},
		&[]int32{1}[0]))
	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	err = testWebComponentRepository.Store(createTestWebComponent(
		"test-other-name",
		"other-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test-path",
					},
					{
						ContextName: "test-name",
					},
				},
			},
		},
		&[]int32{0}[0]))

	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	err = testMicroFrontendRepository.Store(createTestMicroFrontend("test-microfrontend", []string{}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	err = testMicroFrontendRepository.Store(createTestMicroFrontend("other-microfrontend", []string{}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&logr.Logger{},
	)

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("test-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"test-microfrontend": createTestMicroFrontendSpec("test-microfrontend", []string{}),
		})
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))
	// Act
	take := 1
	statusCode, actualContextArea := callGetContextArea(t, polyfeaApiService, ctx, "test-name", "test-path", &take, nil)

	// Assert
	if statusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %v", statusCode)
	}

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func TestPolyfeaApiServiceGetContextAreaMultipleElementsTakeOneCorrectComponentIsSelected(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryRepository[*v1alpha1.WebComponent]()

	err := testWebComponentRepository.Store(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test-path",
					},
					{
						ContextName: "test-name",
					},
				},
			},
		},
		&[]int32{1}[0]))
	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	err = testWebComponentRepository.Store(createTestWebComponent(
		"test-other-name",
		"other-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test-path",
					},
					{
						ContextName: "test-name",
					},
				},
			},
		},
		&[]int32{10}[0]))
	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	err = testMicroFrontendRepository.Store(createTestMicroFrontend("test-microfrontend", []string{}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	err = testMicroFrontendRepository.Store(createTestMicroFrontend("other-microfrontend", []string{}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&logr.Logger{},
	)

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("other-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"other-microfrontend": createTestMicroFrontendSpec("other-microfrontend", []string{}),
		})
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	// Act
	take := 1
	statusCode, actualContextArea := callGetContextArea(t, polyfeaApiService, ctx, "test-name", "test-path", &take, nil)

	// Assert
	if statusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %v", statusCode)
	}

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func TestPolyfeaApiServiceGetContextAreaMicroFrontendDependsOnIsEvaluated(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryRepository[*v1alpha1.WebComponent]()

	err := testWebComponentRepository.Store(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test-path",
					},
					{
						ContextName: "test-name",
					},
				},
			},
		},
		&[]int32{1}[0]))
	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	err = testWebComponentRepository.Store(createTestWebComponent(
		"test-other-name",
		"other-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test-path",
					},
					{
						ContextName: "test-name",
					},
				},
			},
		},
		&[]int32{0}[0]))
	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	err = testMicroFrontendRepository.Store(createTestMicroFrontend("test-microfrontend", []string{"test-dependency"}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	err = testMicroFrontendRepository.Store(createTestMicroFrontend("other-microfrontend", []string{"test-dependency"}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	err = testMicroFrontendRepository.Store(createTestMicroFrontend("test-dependency", []string{"yet-another-test-dependency"}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	err = testMicroFrontendRepository.Store(createTestMicroFrontend("yet-another-test-dependency", []string{}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&logr.Logger{},
	)

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("test-microfrontend"),
			createTestElementSpec("other-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"test-microfrontend":          createTestMicroFrontendSpec("test-microfrontend", []string{"test-dependency"}),
			"other-microfrontend":         createTestMicroFrontendSpec("other-microfrontend", []string{"test-dependency"}),
			"test-dependency":             createTestMicroFrontendSpec("test-dependency", []string{"yet-another-test-dependency"}),
			"yet-another-test-dependency": createTestMicroFrontendSpec("yet-another-test-dependency", []string{}),
		})
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	// Act
	take := 0
	statusCode, actualContextArea := callGetContextArea(t, polyfeaApiService, ctx, "test-name", "test-path", &take, nil)

	// Assert
	if statusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %v", statusCode)
	}

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func TestPolyfeaApiServiceGetContextAreaMicroFrontendDependencyMissingErrorIsReturned(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryRepository[*v1alpha1.WebComponent]()

	err := testWebComponentRepository.Store(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test-path",
					},
					{
						ContextName: "test-name",
					},
				},
			},
		},
		&[]int32{1}[0]))
	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	err = testWebComponentRepository.Store(createTestWebComponent(
		"test-other-name",
		"other-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test-path",
					},
					{
						ContextName: "test-name",
					},
				},
			},
		},
		&[]int32{0}[0]))
	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	err = testMicroFrontendRepository.Store(createTestMicroFrontend("test-microfrontend", []string{"test-dependency"}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	err = testMicroFrontendRepository.Store(createTestMicroFrontend("other-microfrontend", []string{"test-dependency"}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	err = testMicroFrontendRepository.Store(createTestMicroFrontend("test-dependency", []string{"yet-another-test-dependency"}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&logr.Logger{},
	)

	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	// Act
	take := 0
	statusCode, _ := callGetContextArea(t, polyfeaApiService, ctx, "test-name", "test-path", &take, nil)

	// Assert
	// Errors are now returned as HTTP 500, not error values
	if statusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %v", statusCode)
	}
}

func TestPolyfeaApiServiceGetContextAreaMicroFrontendCircularDependencyErrorIsReturned(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryRepository[*v1alpha1.WebComponent]()

	err := testWebComponentRepository.Store(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test-path",
					},
					{
						ContextName: "test-name",
					},
				},
			},
		},
		&[]int32{1}[0]))
	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	err = testWebComponentRepository.Store(createTestWebComponent(
		"test-other-name",
		"other-microfrontend",
		[]v1alpha1.DisplayRules{
			{
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test-path",
					},
					{
						ContextName: "test-name",
					},
				},
			},
		},
		&[]int32{0}[0]))
	if err != nil {
		t.Fatalf("Failed to store WebComponent in repository: %v", err)
	}

	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()

	err = testMicroFrontendRepository.Store(createTestMicroFrontend("test-microfrontend", []string{"test-dependency"}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	err = testMicroFrontendRepository.Store(createTestMicroFrontend("other-microfrontend", []string{"test-dependency"}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	err = testMicroFrontendRepository.Store(createTestMicroFrontend("test-dependency", []string{"test-microfrontend"}, "test-frontend-class", true))
	if err != nil {
		t.Fatalf("Failed to store MicroFrontend in repository: %v", err)
	}

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&logr.Logger{},
	)

	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	// Act
	take := 0
	statusCode, _ := callGetContextArea(t, polyfeaApiService, ctx, "test-name", "test-path", &take, nil)

	// Assert
	// Errors are now returned as HTTP 500, not error values
	if statusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %v", statusCode)
	}
}

func TestPolyfeaApiServiceGetContextAreaInvalidHeaders(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryRepository[*v1alpha1.WebComponent]()
	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&logr.Logger{},
	)

	expectedContextArea := createTestContextArea([]generated.ElementSpec{}, map[string]generated.MicrofrontendSpec{})
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	headers := http.Header{}
	headers.Set("invalid-header", "invalid-value")

	// Act
	take := 10
	statusCode, actualContextArea := callGetContextArea(t, polyfeaApiService, ctx, "test-name", "test-path", &take, headers)

	// Assert
	if statusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %v", statusCode)
	}

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func setupRepositories() (repository.Repository[*v1alpha1.WebComponent], repository.Repository[*v1alpha1.MicroFrontend]) {
	webComponentRepo := repository.NewInMemoryRepository[*v1alpha1.WebComponent]()
	microFrontendRepo := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	return webComponentRepo, microFrontendRepo
}

func setupContext(basePath string, frontendClassName string) context.Context {
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, basePath)
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass(frontendClassName, basePath))
	return ctx
}

// Helper to call GetContextArea with the new API
func callGetContextArea(t *testing.T, service *PolyfeaApiService, ctx context.Context, name, path string, take *int, headers http.Header) (int, generated.ContextArea) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/context-area/"+name, nil)
	req = req.WithContext(ctx)
	if headers != nil {
		req.Header = headers
	}
	service.GetContextArea(w, req, name, generated.GetContextAreaParams{Path: path, Take: take})

	var result generated.ContextArea
	if w.Code == http.StatusOK && w.Body.Len() > 0 {
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
	}
	return w.Code, result
}

func createTestContextArea(expectedElements []generated.ElementSpec, expectedMicroFrontends map[string]generated.MicrofrontendSpec) generated.ContextArea {
	return generated.ContextArea{
		Elements:       expectedElements,
		Microfrontends: &expectedMicroFrontends,
	}
}

func createTestWebComponent(objecName string, microFrontendName string, displayRules []v1alpha1.DisplayRules, priority *int32) *v1alpha1.WebComponent {

	var mfn *string
	if len(microFrontendName) == 0 {
		mfn = nil
	} else {
		mfn = &microFrontendName
	}

	return &v1alpha1.WebComponent{
		ObjectMeta: v1.ObjectMeta{
			Name:      objecName,
			Namespace: "default",
		},
		Spec: v1alpha1.WebComponentSpec{
			MicroFrontend: mfn,
			Element:       &[]string{"test-tag-name"}[0],
			Attributes: []v1alpha1.Attribute{
				{
					Name:  "test-attribute-name",
					Value: runtime.RawExtension{Raw: []byte(`"test-attribute-value"`)},
				},
			},
			DisplayRules: displayRules,
			Priority:     priority,
			Style: []v1alpha1.Style{
				{
					Name:  "test-style-name",
					Value: "test-style-value",
				},
			},
		},
	}
}

func createTestMicroFrontend(objecName string, dependsOn []string, frontendClass string, proxy bool) *v1alpha1.MicroFrontend {
	return &v1alpha1.MicroFrontend{
		ObjectMeta: v1.ObjectMeta{
			Name:      objecName,
			Namespace: "default",
		},
		Spec: v1alpha1.MicroFrontendSpec{
			Service:       &[]string{"http://test-service.test-namespace.svc.cluster.local"}[0],
			Proxy:         &[]bool{proxy}[0],
			CacheStrategy: "none",
			CacheControl:  &[]string{"no-cache"}[0],
			FrontendClass: &[]string{frontendClass}[0],
			DependsOn:     dependsOn,
			ModulePath:    &[]string{"test-module"}[0],
			StaticResources: []v1alpha1.StaticResources{{
				Kind: "test-type",
				Path: "test-uri",
				Attributes: []v1alpha1.Attribute{
					{
						Name:  "test-attribute-name",
						Value: runtime.RawExtension{Raw: []byte(`"test-attribute-value"`)},
					},
				},
				WaitOnLoad: true,
				Proxy:      &[]bool{proxy}[0],
			}},
		},
	}
}

func createTestMicroFrontendClass(frontendClassName string, baseUri string) *v1alpha1.MicroFrontendClass {
	return &v1alpha1.MicroFrontendClass{
		ObjectMeta: v1.ObjectMeta{
			Name:      frontendClassName,
			Namespace: "default",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			BaseUri:         &[]string{baseUri}[0],
			UserRolesHeader: "test-user-roles-header",
		},
	}
}

func createTestElementSpec(microFrontendName string) generated.ElementSpec {
	var mfPtr *string
	if microFrontendName != "" {
		mfPtr = &microFrontendName
	}
	attr := map[string]string{
		"test-attribute-name": "test-attribute-value",
	}
	style := map[string]string{
		"test-style-name": "test-style-value",
	}
	return generated.ElementSpec{
		Microfrontend: mfPtr,
		TagName:       "test-tag-name",
		Attributes:    &attr,
		Style:         &style,
	}
}

func createTestMicroFrontendSpec(microfrontendName string, dependsOn []string) generated.MicrofrontendSpec {
	href := "./polyfea/proxy/default/" + microfrontendName + "/test-uri"
	module := "./polyfea/proxy/default/" + microfrontendName + "/test-module"
	kind := generated.MicrofrontendResourceKind("test-type")
	attr := map[string]string{
		"test-attribute-name": "test-attribute-value",
	}
	waitOnLoad := true
	resources := []generated.MicrofrontendResource{
		{
			Kind:       &kind,
			Href:       &href,
			Attributes: &attr,
			WaitOnLoad: &waitOnLoad,
		},
	}
	return generated.MicrofrontendSpec{
		DependsOn: &dependsOn,
		Module:    &module,
		Resources: &resources,
	}
}
