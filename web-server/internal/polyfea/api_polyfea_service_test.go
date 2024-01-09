package polyfea

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
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
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("test-microfrontend", []string{"test-dependency"}, "test-module"))
	testMicroFrontedClassRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontendClass]()

	polyfeaApiService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		testMicroFrontedClassRepository,
		map[string]string{})

	expectedContextArea := createTestContextArea()

	// Act
	actualContextAreaResponse, err := polyfeaApiService.GetContextArea(context.TODO(), "test-name", "test-path", 10, map[string][]string{})

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	actualContextArea := actualContextAreaResponse.Body.(ContextArea)

	expectedContextAreaBytes, _ := json.Marshal(expectedContextArea)
	actualContextAreaBytes, _ := json.Marshal(actualContextArea)

	if string(expectedContextAreaBytes) != string(actualContextAreaBytes) {
		t.Errorf("Expected %v, got %v", string(expectedContextAreaBytes), string(actualContextAreaBytes))
	}
}

func createTestContextArea() ContextArea {
	return ContextArea{
		Elements: []ElementSpec{
			{
				Microfrontend: "test-microfrontend",
				TagName:       "test-tag-name",
				Attributes: map[string]string{
					"test-attribute-name": "test-attribute-value",
				},
				Style: map[string]string{
					"test-style-name": "test-style-value",
				},
			},
		},
		Microfrontends: map[string]MicrofrontendSpec{
			"test-microfrontend": {
				DependsOn: []string{
					"test-dependency",
				},
				Module: "test-module",
				Resources: []MicrofrontendResource{
					{
						Kind: "test-type",
						Href: "test-uri",
						Attributes: map[string]string{
							"test-attribute-name": "test-attribute-value",
						},
						WaitOnLoad: true,
					},
				},
			},
		},
	}
}

func createTestWebComponent(objecName string, microFrontendName string, element string, displayRules []v1alpha1.DisplayRules, priority *int32) *v1alpha1.WebComponent {
	return &v1alpha1.WebComponent{
		ObjectMeta: v1.ObjectMeta{
			Name: objecName,
		},
		Spec: v1alpha1.WebComponentSpec{
			MicroFrontend: &microFrontendName,
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

func createTestMicroFrontend(objecName string, dependsOn []string, modulePath string) *v1alpha1.MicroFrontend {
	return &v1alpha1.MicroFrontend{
		ObjectMeta: v1.ObjectMeta{
			Name: objecName,
		},
		Spec: v1alpha1.MicroFrontendSpec{
			DependsOn:  dependsOn,
			ModulePath: &modulePath,
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
			}},
		},
	}
}
