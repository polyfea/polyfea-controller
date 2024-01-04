package webserver

import (
	"github.com/gorilla/mux"
	"github.com/polyfea/polyfea-controller/web-server/api"
	"github.com/polyfea/polyfea-controller/web-server/internal/polyfea"
)

func SetupRouter() *mux.Router {
	baseRouter := polyfea.NewRouter()

	baseRouter.HandleFunc("/openapi", api.HandleOpenApi)

	return baseRouter
}
