package repository

import (
	"encoding/json"
	"testing"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestInMemoryPolyfeaRepositoryMicrofrontendStoredCanBeRetrieved(t *testing.T) {

	// Arrange
	repository := NewInMemoryPolyfeaRepository()
	expectedMicrofrontend := createTestMicrofrontend()

	// Act
	repository.StoreMicrofrontend(expectedMicrofrontend)
	result, _ := repository.GetMicrofrontends(func(mf v1alpha1.MicroFrontend) bool {
		return mf.Name == expectedMicrofrontend.Name
	})

	// Assert
	expectedMicrofrontendBytes, _ := json.Marshal(expectedMicrofrontend)
	resultMicrofrontendBytes, _ := json.Marshal(result[0])

	if string(expectedMicrofrontendBytes) != string(resultMicrofrontendBytes) {
		t.Errorf("Expected microfrontend %v, but got %v", string(expectedMicrofrontendBytes), string(resultMicrofrontendBytes))
	}
}

func TestInMemoryPolyfeaRepositoryMicrofrontendsNotFoundReturnsEmptySlice(t *testing.T) {

	// Arrange
	repository := NewInMemoryPolyfeaRepository()

	// Act
	result, _ := repository.GetMicrofrontends(func(mf v1alpha1.MicroFrontend) bool {
		return mf.Name == "test"
	})

	// Assert
	if len(result) != 0 {
		t.Errorf("Expected empty slice, but got %v", result)
	}
}

func TestInMemoryPolyfeaRepositoryWebcomponentStoredCanBeRetrieved(t *testing.T) {

	// Arrange
	repository := NewInMemoryPolyfeaRepository()
	expectedWebComponent := createTestWebComponent()

	// Act
	repository.StoreWebComponent(expectedWebComponent)
	result, _ := repository.GetWebComponents(func(mf v1alpha1.WebComponent) bool {
		return mf.Name == expectedWebComponent.Name
	})

	// Assert
	expectedWebComponentBytes, _ := json.Marshal(expectedWebComponent)
	resultWebComponentBytes, _ := json.Marshal(result[0])

	if string(expectedWebComponentBytes) != string(resultWebComponentBytes) {
		t.Errorf("Expected microfrontend %v, but got %v", string(expectedWebComponentBytes), string(resultWebComponentBytes))
	}
}

func TestInMemoryPolyfeaRepositoryWebcomponentNotFoundReturnsEmptySlice(t *testing.T) {

	// Arrange
	repository := NewInMemoryPolyfeaRepository()

	// Act
	result, _ := repository.GetWebComponents(func(mf v1alpha1.WebComponent) bool {
		return mf.Name == "test"
	})

	// Assert
	if len(result) != 0 {
		t.Errorf("Expected empty slice, but got %v", result)
	}
}

func TestInMemoryPolyfeaRepositoryMicrofrontendClassStoredCanBeRetrieved(t *testing.T) {

	// Arrange
	repository := NewInMemoryPolyfeaRepository()
	expectedMicrofrontendClass := createTestMicrofrontendClass()

	// Act
	repository.StoreMicrofrontendClass(expectedMicrofrontendClass)
	result, _ := repository.GetMicrofrontendClasses(func(mf v1alpha1.MicroFrontendClass) bool {
		return mf.Name == expectedMicrofrontendClass.Name
	})

	// Assert
	expectedMicrofrontendClassBytes, _ := json.Marshal(expectedMicrofrontendClass)
	resultMicrofrontendClassBytes, _ := json.Marshal(result[0])

	if string(expectedMicrofrontendClassBytes) != string(resultMicrofrontendClassBytes) {
		t.Errorf("Expected microfrontend %v, but got %v", string(expectedMicrofrontendClassBytes), string(resultMicrofrontendClassBytes))
	}
}

func TestInMemoryPolyfeaRepositoryMicrofrontendClassNotFoundReturnsEmptySlice(t *testing.T) {

	// Arrange
	repository := NewInMemoryPolyfeaRepository()

	// Act
	result, _ := repository.GetMicrofrontendClasses(func(mf v1alpha1.MicroFrontendClass) bool {
		return mf.Name == "test"
	})

	// Assert
	if len(result) != 0 {
		t.Errorf("Expected empty slice, but got %v", result)
	}
}

func createTestMicrofrontend() v1alpha1.MicroFrontend {
	return v1alpha1.MicroFrontend{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
		Spec: v1alpha1.MicroFrontendSpec{
			Service: &v1alpha1.ServiceReference{
				Name: "test",
			},
			CacheStrategy: "none",
			CacheControl:  &[]string{"no-cache"}[0],
			ModulePath:    &[]string{"test"}[0],
			StaticPaths:   []string{"test"},
			PreloadPaths:  []string{"test"},
			FrontendClass: &[]string{"test"}[0],
			DependsOn:     []string{"test"},
		},
	}
}

func createTestWebComponent() v1alpha1.WebComponent {
	return v1alpha1.WebComponent{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
		Spec: v1alpha1.WebComponentSpec{
			MicroFrontend: &[]string{"test"}[0],
			Element:       &[]string{"test"}[0],
			Attributes: []v1alpha1.Attribute{
				{
					Name:  "test",
					Value: runtime.RawExtension{Raw: []byte("test")},
				},
			},
			DisplayRules: []v1alpha1.DisplayRules{
				{
					AllOf: []v1alpha1.Matcher{
						{
							ContextName: "test",
						},
					},
				},
			},
			Priority: &[]int32{1}[0],
			Style:    &[]string{"test"}[0],
		},
	}
}

func createTestMicrofrontendClass() v1alpha1.MicroFrontendClass {
	return v1alpha1.MicroFrontendClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			BaseUri:   &[]string{"test"}[0],
			CspHeader: "test",
			ExtraHeaders: []v1alpha1.Header{
				{
					Name:  "test",
					Value: "test",
				},
			},
			UserRolesHeader: "test",
			UserHeader:      "test",
		},
	}
}
