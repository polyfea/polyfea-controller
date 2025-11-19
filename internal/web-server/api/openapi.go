package api

import (
	_ "embed"
	"net/http"
)

//go:embed v1alpha1.openapi.yaml
var openapiSpec []byte

func HandleOpenApi(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	_, err := w.Write(openapiSpec)
	if err != nil {
		http.Error(w, "Failed to respond with OpenAPI spec", http.StatusInternalServerError)
	}
}
