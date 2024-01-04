package api

import (
	"net/http"
)

func SetupRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/openapi", HandleOpenApi)

	return router
}
