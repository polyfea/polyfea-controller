package repository

import (
	"encoding/json"
	"testing"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestInMemoryRepository(t *testing.T) {
	t.Run("Store and retrieve item with List", func(t *testing.T) {
		repo := NewInMemoryRepository[*v1alpha1.MicroFrontend]()
		expected := createTestMicrofrontend()

		if err := repo.Store(expected); err != nil {
			t.Fatalf("Store failed: %v", err)
		}

		result, err := repo.List(func(mf *v1alpha1.MicroFrontend) bool {
			return mf.Name == expected.Name
		})
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if len(result) != 1 {
			t.Fatalf("Expected 1 item, got %d", len(result))
		}

		assertEqualMicrofrontend(t, expected, result[0])
	})

	t.Run("Store and retrieve single item with Get", func(t *testing.T) {
		repo := NewInMemoryRepository[*v1alpha1.MicroFrontend]()
		expected := createTestMicrofrontend()

		if err := repo.Store(expected); err != nil {
			t.Fatalf("Store failed: %v", err)
		}

		result, err := repo.Get(expected)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		assertEqualMicrofrontend(t, expected, result)
	})

	t.Run("List returns empty slice when not found", func(t *testing.T) {
		repo := NewInMemoryRepository[*v1alpha1.MicroFrontend]()
		result, err := repo.List(func(mf *v1alpha1.MicroFrontend) bool {
			return mf.Name == "notfound"
		})
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("Delete removes item", func(t *testing.T) {
		repo := NewInMemoryRepository[*v1alpha1.MicroFrontend]()
		expected := createTestMicrofrontend()

		if err := repo.Store(expected); err != nil {
			t.Fatalf("Store failed: %v", err)
		}
		if err := repo.Delete(expected); err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		result, err := repo.List(func(mf *v1alpha1.MicroFrontend) bool {
			return mf.Name == expected.Name
		})
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("Expected empty slice after delete, got %v", result)
		}
	})
}

// assertEqualMicrofrontend compares two MicroFrontend objects using JSON marshaling.
func assertEqualMicrofrontend(t *testing.T, expected, actual *v1alpha1.MicroFrontend) {
	t.Helper()
	expBytes, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("Failed to marshal expected: %v", err)
	}
	actBytes, err := json.Marshal(actual)
	if err != nil {
		t.Fatalf("Failed to marshal actual: %v", err)
	}
	if string(expBytes) != string(actBytes) {
		t.Errorf("Expected %v, got %v", string(expBytes), string(actBytes))
	}
}

func createTestMicrofrontend() *v1alpha1.MicroFrontend {
	return &v1alpha1.MicroFrontend{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
		Spec: v1alpha1.MicroFrontendSpec{
			Service:       ptr("http://test-service.test-namespace.svc.cluster.local"),
			CacheStrategy: "none",
			CacheControl:  ptr("no-cache"),
			ModulePath:    ptr("test"),
			StaticResources: []v1alpha1.StaticResources{{
				Kind: "test",
				Path: "test",
			}},
			FrontendClass: ptr("test"),
			DependsOn:     []string{"test"},
		},
	}
}

func ptr[T any](v T) *T {
	return &v
}
