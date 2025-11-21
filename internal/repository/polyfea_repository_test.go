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

		storeItem(t, repo, expected)
		result := listItems(t, repo, func(mf *v1alpha1.MicroFrontend) bool {
			return mf.Name == expected.Name
		})
		assertListLength(t, result, 1)
		assertEqualMicrofrontend(t, expected, result[0])
	})

	t.Run("Store and retrieve single item with Get", func(t *testing.T) {
		repo := NewInMemoryRepository[*v1alpha1.MicroFrontend]()
		expected := createTestMicrofrontend()

		storeItem(t, repo, expected)
		result := getItem(t, repo, expected)
		assertEqualMicrofrontend(t, expected, result)
	})

	t.Run("List returns empty slice when not found", func(t *testing.T) {
		repo := NewInMemoryRepository[*v1alpha1.MicroFrontend]()
		result := listItems(t, repo, func(mf *v1alpha1.MicroFrontend) bool {
			return mf.Name == "notfound"
		})
		assertListLength(t, result, 0)
	})

	t.Run("Delete removes item", func(t *testing.T) {
		repo := NewInMemoryRepository[*v1alpha1.MicroFrontend]()
		expected := createTestMicrofrontend()

		storeItem(t, repo, expected)
		deleteItem(t, repo, expected)
		result := listItems(t, repo, func(mf *v1alpha1.MicroFrontend) bool {
			return mf.Name == expected.Name
		})
		assertListLength(t, result, 0)
	})
}

func storeItem(t *testing.T, repo *InMemoryRepository[*v1alpha1.MicroFrontend], item *v1alpha1.MicroFrontend) {
	if err := repo.Store(item); err != nil {
		t.Fatalf("Store failed: %v", err)
	}
}

func getItem(t *testing.T, repo *InMemoryRepository[*v1alpha1.MicroFrontend], item *v1alpha1.MicroFrontend) *v1alpha1.MicroFrontend {
	result, err := repo.Get(item)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	return result
}

func listItems(t *testing.T, repo *InMemoryRepository[*v1alpha1.MicroFrontend], filter func(*v1alpha1.MicroFrontend) bool) []*v1alpha1.MicroFrontend {
	result, err := repo.List(filter)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	return result
}

func deleteItem(t *testing.T, repo *InMemoryRepository[*v1alpha1.MicroFrontend], item *v1alpha1.MicroFrontend) {
	if err := repo.Delete(item); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func assertListLength(t *testing.T, list []*v1alpha1.MicroFrontend, expectedLength int) {
	if len(list) != expectedLength {
		t.Fatalf("Expected %d items, got %d", expectedLength, len(list))
	}
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
	serviceName := "test-service"
	serviceNamespace := "test-namespace"
	servicePort := int32(80)
	serviceScheme := "http"
	return &v1alpha1.MicroFrontend{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
		Spec: v1alpha1.MicroFrontendSpec{
			Service: &v1alpha1.ServiceReference{
				Name:      &serviceName,
				Namespace: &serviceNamespace,
				Port:      &servicePort,
				Scheme:    &serviceScheme,
			},
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
