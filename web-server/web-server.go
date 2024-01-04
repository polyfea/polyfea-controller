package webserver

import (
	"github.com/gorilla/mux"
	"github.com/polyfea/polyfea-controller/web-server/api"
	"github.com/polyfea/polyfea-controller/web-server/internal/polyfea"
)

func SetupRouter() *mux.Router {
	polyfeaAPIService := polyfea.NewPolyfeaAPIService()
	polyfeaAPIController := polyfea.NewPolyfeaAPIController(polyfeaAPIService)

	router := polyfea.NewRouter(polyfeaAPIController)

	router.HandleFunc("/openapi", api.HandleOpenApi)

	return router
}
