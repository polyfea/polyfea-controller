package webserver

import (
	"github.com/gorilla/mux"
	"github.com/polyfea/polyfea-controller/web-server/api"
	"github.com/polyfea/polyfea-controller/web-server/internal/polyfea"
)

func SetupRouter() *mux.Router {
	router := polyfea.NewRouter()

	router.HandleFunc("/openapi", api.HandleOpenApi)

	return router
}
