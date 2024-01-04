package polyfea

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/polyfea/polyfea-controller/web-server/api"
)

var (
	testServer *httptest.Server
)

func TestPolyfeaApiServiceGetContextAreaReturnsNotImplemented(t *testing.T) {

	// Arrange
	s := NewPolyfeaAPIService()
	ctx := context.Background()

	// Act
	response, err := s.GetContextArea(ctx, "test", "test", 0)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response.Code != 501 {
		t.Errorf("Expected response code to be %d, got %d", 501, response.Code)
	}
}

func TestPolyfeaApiServiceGetStaticConfigReturnsNotImplemented(t *testing.T) {

	// Arrange
	s := NewPolyfeaAPIService()
	ctx := context.Background()

	// Act
	response, err := s.GetStaticConfig(ctx)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response.Code != 501 {
		t.Errorf("Expected response code to be %d, got %d", 501, response.Code)
	}
}

func setupRouter() *mux.Router {
	polyfeaAPIService := NewPolyfeaAPIService()
	polyfeaAPIController := NewPolyfeaAPIController(polyfeaAPIService)

	router := NewRouter(polyfeaAPIController)

	router.HandleFunc("/openapi", api.HandleOpenApi)

	return router
}
