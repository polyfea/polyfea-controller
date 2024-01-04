package polyfea

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
