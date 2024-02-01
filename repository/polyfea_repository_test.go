package repository

import (
	"encoding/json"
	"testing"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestInMemoryPolyfeaRepositoryItemStoredCanBeRetrieved(t *testing.T) {

	// Arrange
	repository := NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	expectedMicrofrontend := createTestMicrofrontend()

	// Act
	repository.StoreItem(expectedMicrofrontend)
	result, _ := repository.GetItems(func(mf *v1alpha1.MicroFrontend) bool {
		return mf.Name == expectedMicrofrontend.Name
	})

	// Assert
	expectedMicrofrontendBytes, _ := json.Marshal(expectedMicrofrontend)
	resultMicrofrontendBytes, _ := json.Marshal(result[0])

	if string(expectedMicrofrontendBytes) != string(resultMicrofrontendBytes) {
		t.Errorf("Expected microfrontend %v, but got %v", string(expectedMicrofrontendBytes), string(resultMicrofrontendBytes))
	}
}

func TestInMemoryPolyfeaRepositoryItemStoredCanBeRetrievedAsSingleItem(t *testing.T) {

	// Arrange
	repository := NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	expectedMicrofrontend := createTestMicrofrontend()

	// Act
	repository.StoreItem(expectedMicrofrontend)
	result, _ := repository.GetItem(expectedMicrofrontend)

	// Assert
	expectedMicrofrontendBytes, _ := json.Marshal(expectedMicrofrontend)
	resultMicrofrontendBytes, _ := json.Marshal(result)

	if string(expectedMicrofrontendBytes) != string(resultMicrofrontendBytes) {
		t.Errorf("Expected microfrontend %v, but got %v", string(expectedMicrofrontendBytes), string(resultMicrofrontendBytes))
	}
}

func TestInMemoryPolyfeaRepositoryItemsNotFoundReturnsEmptySlice(t *testing.T) {

	// Arrange
	repository := NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()

	// Act
	result, _ := repository.GetItems(func(mf *v1alpha1.MicroFrontend) bool {
		return mf.Name == "test"
	})

	// Assert
	if len(result) != 0 {
		t.Errorf("Expected empty slice, but got %v", result)
	}
}

func TestInMemoryPolyfeaRepositoryItemDeletedCannotBeRetrieved(t *testing.T) {

	// Arrange
	repository := NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	expectedMicrofrontend := createTestMicrofrontend()

	// Act
	repository.StoreItem(expectedMicrofrontend)
	repository.DeleteItem(expectedMicrofrontend)
	result, _ := repository.GetItems(func(mf *v1alpha1.MicroFrontend) bool {
		return mf.Name == expectedMicrofrontend.Name
	})

	// Assert
	if len(result) != 0 {
		t.Errorf("Expected empty slice, but got %v", result)
	}
}

func createTestMicrofrontend() *v1alpha1.MicroFrontend {
	return &v1alpha1.MicroFrontend{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
		Spec: v1alpha1.MicroFrontendSpec{
			Service:       &[]string{"http://test-service.test-namespace.svc.cluster.local"}[0],
			CacheStrategy: "none",
			CacheControl:  &[]string{"no-cache"}[0],
			ModulePath:    &[]string{"test"}[0],
			StaticResources: []v1alpha1.StaticResources{{
				Kind: "test",
				Path: "test",
			}},
			FrontendClass: &[]string{"test"}[0],
			DependsOn:     []string{"test"},
		},
	}
}
