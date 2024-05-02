package polyfea

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"testing"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
	"github.com/polyfea/polyfea-controller/web-server/internal/polyfea/generated"
	"github.com/rs/zerolog"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestPolyfeaApiServiceGetContextAreaReturnsContextAreaIfRepositoryContainsMatchingWebComponents(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.WebComponent]()
	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		"test-tag-name",
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
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&zerolog.Logger{},
	)

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("test-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"test-microfrontend": createTestMicroFrontendSpec("test-microfrontend", []string{}, true),
		})
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	// Act
	actualContextAreaResponse, err := polyfeaApiService.GetContextArea(ctx, "test-name", "test-path", 10, map[string][]string{})

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	actualContextArea := actualContextAreaResponse.Body.(generated.ContextArea)

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsContextAreaWithExtraHeaders(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.WebComponent]()
	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		"test-tag-name",
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
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true))

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
		&zerolog.Logger{},
	)

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("test-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"test-microfrontend": createTestMicroFrontendSpec("test-microfrontend", []string{}, true),
		})
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, mfc)

	// Act
	actualContextAreaResponse, err := polyfeaApiService.GetContextArea(ctx, "test-name", "test-path", 10, map[string][]string{})

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	actualContextArea := actualContextAreaResponse.Body.(generated.ContextArea)

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}

	if !slices.Contains(actualContextAreaResponse.Headers["test-header"], "test-value") {
		t.Errorf("Expected %v, got %v", "test-value", actualContextAreaResponse.Headers["test-header"][0])
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsContextAreaIfNoneOfIsMatching(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.WebComponent]()
	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		"test-tag-name",
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
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&zerolog.Logger{},
	)

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("test-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"test-microfrontend": createTestMicroFrontendSpec("test-microfrontend", []string{}, true),
		})
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	// Act
	actualContextAreaResponse, err := polyfeaApiService.GetContextArea(ctx, "test-name", "test-path", 10, map[string][]string{})

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	actualContextArea := actualContextAreaResponse.Body.(generated.ContextArea)

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsContextAreaIfAnyOfIsMatching(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.WebComponent]()
	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		"test-tag-name",
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
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&zerolog.Logger{},
	)

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("test-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"test-microfrontend": createTestMicroFrontendSpec("test-microfrontend", []string{}, true),
		})
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	// Act
	actualContextAreaResponse, err := polyfeaApiService.GetContextArea(ctx, "test-name", "test-path", 10, map[string][]string{})

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	actualContextArea := actualContextAreaResponse.Body.(generated.ContextArea)

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsContextAreaIfComplexCombinationIsMatching(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.WebComponent]()
	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		"test-tag-name",
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
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&zerolog.Logger{},
	)

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("test-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"test-microfrontend": createTestMicroFrontendSpec("test-microfrontend", []string{}, true),
		})
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	// Act
	actualContextAreaResponse, err := polyfeaApiService.GetContextArea(ctx, "test-name", "test-path", 10, map[string][]string{})

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	actualContextArea := actualContextAreaResponse.Body.(generated.ContextArea)

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsContextAreaIfComplexMatcherIsMatching(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.WebComponent]()
	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		"test-tag-name",
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
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&zerolog.Logger{},
	)

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("test-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"test-microfrontend": createTestMicroFrontendSpec("test-microfrontend", []string{}, true),
		})
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	headers := http.Header{}
	headers.Set("test-user-roles-header", "some-different-role")
	headers.Add("test-user-roles-header", "test-role, test-other-role")

	// Act
	actualContextAreaResponse, err := polyfeaApiService.GetContextArea(ctx, "test-name", "test-path", 10, headers)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	actualContextArea := actualContextAreaResponse.Body.(generated.ContextArea)

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsElementWithoutMicrofrontendIfItHasNoMicrofrontends(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.WebComponent]()
	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-name",
		"",
		"test-tag-name",
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
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&zerolog.Logger{},
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
	actualContextAreaResponse, err := polyfeaApiService.GetContextArea(ctx, "test-name", "test-path", 10, headers)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	actualContextArea := actualContextAreaResponse.Body.(generated.ContextArea)

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsEmptyIfRoleMatcherIsNotMatching(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.WebComponent]()
	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		"test-tag-name",
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
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&zerolog.Logger{},
	)

	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	headers := http.Header{}
	headers.Set("test-user-roles-header", "some-different-role")
	headers.Add("test-user-roles-header", "test-role, test-other-role")

	// Act
	response, err := polyfeaApiService.GetContextArea(ctx, "test-name", "test.*", 10, headers)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Code != 200 {
		t.Errorf("Expected 200, got %v", response.Code)
	}

	if len(response.Body.(generated.ContextArea).Elements) != 0 {
		t.Errorf("Expected 0 elements, got %v", len(response.Body.(generated.ContextArea).Elements))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsEmptyIfContextMatcherIsNotMatching(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.WebComponent]()
	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		"test-tag-name",
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
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&zerolog.Logger{},
	)

	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	headers := http.Header{}
	headers.Set("test-user-roles-header", "some-different-role")
	headers.Add("test-user-roles-header", "test-role, test-other-role")

	// Act
	response, err := polyfeaApiService.GetContextArea(ctx, "test-name", "test.*", 10, headers)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Code != 200 {
		t.Errorf("Expected 200, got %v", response.Code)
	}

	if len(response.Body.(generated.ContextArea).Elements) != 0 {
		t.Errorf("Expected 0 elements, got %v", len(response.Body.(generated.ContextArea).Elements))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsEmptyIfPathIsNotMatching(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.WebComponent]()
	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		"test-tag-name",
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
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&zerolog.Logger{},
	)

	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	headers := http.Header{}
	headers.Set("test-user-roles-header", "some-different-role")
	headers.Add("test-user-roles-header", "test-role, test-other-role")

	// Act
	response, err := polyfeaApiService.GetContextArea(ctx, "test-name", "sometest", 10, headers)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Code != 200 {
		t.Errorf("Expected 200, got %v", response.Code)
	}

	if len(response.Body.(generated.ContextArea).Elements) != 0 {
		t.Errorf("Expected 0 elements, got %v", len(response.Body.(generated.ContextArea).Elements))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsEmptyIfNoneOfIsMatching(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.WebComponent]()
	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		"test-tag-name",
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
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&zerolog.Logger{},
	)

	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	headers := http.Header{}
	headers.Set("test-user-roles-header", "some-different-role")
	headers.Add("test-user-roles-header", "test-role, test-other-role")

	// Act
	response, err := polyfeaApiService.GetContextArea(ctx, "test-name", "sometest.*", 10, headers)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Code != 200 {
		t.Errorf("Expected 200, got %v", response.Code)
	}

	if len(response.Body.(generated.ContextArea).Elements) != 0 {
		t.Errorf("Expected 0 elements, got %v", len(response.Body.(generated.ContextArea).Elements))
	}
}

func TestPolyfeaApiServiceGetContextAreaReturnsEmptyIfAnyOfIsNotMatching(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.WebComponent]()
	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		"test-tag-name",
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
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&zerolog.Logger{},
	)

	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	headers := http.Header{}
	headers.Set("test-user-roles-header", "some-different-role")
	headers.Add("test-user-roles-header", "test-role, test-other-role")

	// Act
	response, err := polyfeaApiService.GetContextArea(ctx, "test-name", "sometest.*", 10, headers)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Code != 200 {
		t.Errorf("Expected 200, got %v", response.Code)
	}

	if len(response.Body.(generated.ContextArea).Elements) != 0 {
		t.Errorf("Expected 0 elements, got %v", len(response.Body.(generated.ContextArea).Elements))
	}
}

func TestPolyfeaApiServiceGetContextAreaMultipleElementsTakeOneOnlyOneElementReturned(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.WebComponent]()

	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		"test-tag-name",
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

	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-other-name",
		"other-microfrontend",
		"test-tag-name",
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
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true))
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("other-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&zerolog.Logger{},
	)

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("test-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"test-microfrontend": createTestMicroFrontendSpec("test-microfrontend", []string{}, true),
		})
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))
	// Act
	actualContextAreaResponse, err := polyfeaApiService.GetContextArea(ctx, "test-name", "test-path", 1, map[string][]string{})

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	actualContextArea := actualContextAreaResponse.Body.(generated.ContextArea)

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func TestPolyfeaApiServiceGetContextAreaMultipleElementsTakeOneCorrectComponentIsSelected(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.WebComponent]()

	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		"test-tag-name",
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

	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-other-name",
		"other-microfrontend",
		"test-tag-name",
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
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true))
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("other-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&zerolog.Logger{},
	)

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("other-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"other-microfrontend": createTestMicroFrontendSpec("other-microfrontend", []string{}, true),
		})
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	// Act
	actualContextAreaResponse, err := polyfeaApiService.GetContextArea(ctx, "test-name", "test-path", 1, map[string][]string{})

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	actualContextArea := actualContextAreaResponse.Body.(generated.ContextArea)

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func TestPolyfeaApiServiceGetContextAreaMicroFrontendDependsOnIsEvaluated(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.WebComponent]()

	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		"test-tag-name",
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

	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-other-name",
		"other-microfrontend",
		"test-tag-name",
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
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-microfrontend", []string{"test-dependency"}, "test-module", "test-frontend-class", true))
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("other-microfrontend", []string{"test-dependency"}, "test-module", "test-frontend-class", true))

	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-dependency", []string{"yet-another-test-dependency"}, "test-module", "test-frontend-class", true))
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("yet-another-test-dependency", []string{}, "test-module", "test-frontend-class", true))

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&zerolog.Logger{},
	)

	expectedContextArea := createTestContextArea(
		[]generated.ElementSpec{
			createTestElementSpec("test-microfrontend"),
			createTestElementSpec("other-microfrontend"),
		},
		map[string]generated.MicrofrontendSpec{
			"test-microfrontend":          createTestMicroFrontendSpec("test-microfrontend", []string{"test-dependency"}, true),
			"other-microfrontend":         createTestMicroFrontendSpec("other-microfrontend", []string{"test-dependency"}, true),
			"test-dependency":             createTestMicroFrontendSpec("test-dependency", []string{"yet-another-test-dependency"}, true),
			"yet-another-test-dependency": createTestMicroFrontendSpec("yet-another-test-dependency", []string{}, true),
		})
	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	// Act
	actualContextAreaResponse, err := polyfeaApiService.GetContextArea(ctx, "test-name", "test-path", 0, map[string][]string{})

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	actualContextArea := actualContextAreaResponse.Body.(generated.ContextArea)

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func TestPolyfeaApiServiceGetContextAreaMicroFrontendDependencyMissingErrorIsReturned(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.WebComponent]()

	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		"test-tag-name",
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

	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-other-name",
		"other-microfrontend",
		"test-tag-name",
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
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-microfrontend", []string{"test-dependency"}, "test-module", "test-frontend-class", true))
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("other-microfrontend", []string{"test-dependency"}, "test-module", "test-frontend-class", true))

	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-dependency", []string{"yet-another-test-dependency"}, "test-module", "test-frontend-class", true))

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&zerolog.Logger{},
	)

	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	// Act
	_, err := polyfeaApiService.GetContextArea(ctx, "test-name", "test-path", 0, map[string][]string{})

	// Assert
	if err == nil {
		t.Errorf("Expected error, got %v", err)
	}

	if err.Error() != "Microfrontend yet-another-test-dependency not found" {
		t.Errorf("Expected error, got %v", err)
	}
}

func TestPolyfeaApiServiceGetContextAreaMicroFrontendCircularDependencyErrorIsReturned(t *testing.T) {
	// Arrange
	testWebComponentRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.WebComponent]()

	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		"test-tag-name",
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

	testWebComponentRepository.StoreItem(createTestWebComponent(
		"test-other-name",
		"other-microfrontend",
		"test-tag-name",
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
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-microfrontend", []string{"test-dependency"}, "test-module", "test-frontend-class", true))
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("other-microfrontend", []string{"test-dependency"}, "test-module", "test-frontend-class", true))

	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-dependency", []string{"test-microfrontend"}, "test-module", "test-frontend-class", true))

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&zerolog.Logger{},
	)

	ctx := context.WithValue(context.TODO(), PolyfeaContextKeyBasePath, "/")
	ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, createTestMicroFrontendClass("test-frontend-class", "/"))

	// Act
	_, err := polyfeaApiService.GetContextArea(ctx, "test-name", "test-path", 0, map[string][]string{})

	// Assert
	if err == nil {
		t.Errorf("Expected error, got %v", err)
	}

	if !strings.Contains(err.Error(), "Circular dependency detected") {
		t.Errorf("Expected error, got %v", err)
	}
}

func createTestContextArea(expectedElements []generated.ElementSpec, expectedMicroFrontends map[string]generated.MicrofrontendSpec) generated.ContextArea {
	return generated.ContextArea{
		Elements:       expectedElements,
		Microfrontends: expectedMicroFrontends,
	}
}

func createTestWebComponent(objecName string, microFrontendName string, element string, displayRules []v1alpha1.DisplayRules, priority *int32) *v1alpha1.WebComponent {

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
			Element:       &element,
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

func createTestMicroFrontend(objecName string, dependsOn []string, modulePath string, frontendClass string, proxy bool) *v1alpha1.MicroFrontend {
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
			ModulePath:    &modulePath,
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
	return generated.ElementSpec{
		Microfrontend: microFrontendName,
		TagName:       "test-tag-name",
		Attributes: map[string]string{
			"test-attribute-name": "test-attribute-value",
		},
		Style: map[string]string{
			"test-style-name": "test-style-value",
		},
	}
}

func createTestMicroFrontendSpec(microfrontendName string, dependsOn []string, withProxy bool) generated.MicrofrontendSpec {
	var href, module string
	if withProxy {
		href = "./polyfea/proxy/default/" + microfrontendName + "/test-uri"
		module = "./polyfea/proxy/default/" + microfrontendName + "/test-module"
	} else {
		href = "test-uri"
		module = "test-module"
	}

	return generated.MicrofrontendSpec{
		DependsOn: dependsOn,
		Module:    module,
		Resources: []generated.MicrofrontendResource{
			{
				Kind: "test-type",
				Href: href,
				Attributes: map[string]string{
					"test-attribute-name": "test-attribute-value",
				},
				WaitOnLoad: true,
			},
		},
	}
}
