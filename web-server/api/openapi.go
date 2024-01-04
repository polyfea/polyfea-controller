package api

import (
	_ "embed"
	"net/http"
)

//go:embed v1alpha1.openapi.yaml
var openapiSpec []byte

func HandleOpenApi(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	w.Write(openapiSpec)
}
