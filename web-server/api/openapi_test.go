package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/polyfea/polyfea-controller/web-server/internal/polyfea"
)

var (
	testServer *httptest.Server
)

func TestMain(t *testing.M) {
	// Setup
	r := polyfea.NewRouter()
	r.HandleFunc("/openapi", HandleOpenApi)

	// Test server
	testServer = httptest.NewServer(r)
	defer testServer.Close()

	t.Run()
}

func TestHandleOpenApi(t *testing.T) {
	// Arrange
	testServerUrl := testServer.URL
	testRoute := "/openapi"
	openapiSpec, _ := os.ReadFile("v1alpha1.openapi.yaml")

	// Act
	response, err := http.Get(testServerUrl + testRoute)
	println(testServerUrl + testRoute)
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Could not read response body: %v", err)
	}

	bodyString := string(body)
	if bodyString != string(openapiSpec) {
		t.Errorf("Expected body to be equal to openapi spec, got %s", bodyString)
	}
}
