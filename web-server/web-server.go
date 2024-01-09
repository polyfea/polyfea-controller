package webserver

import (
	"net/http"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
	"github.com/polyfea/polyfea-controller/web-server/api"
	"github.com/polyfea/polyfea-controller/web-server/internal/polyfea"
	"github.com/polyfea/polyfea-controller/web-server/internal/polyfea/generated"
)

func SetupRouter(
	microFrontendClassRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontendClass],
	microFrontendRepoistory repository.PolyfeaRepository[*v1alpha1.MicroFrontend],
	webComponentRepository repository.PolyfeaRepository[*v1alpha1.WebComponent]) http.Handler {

	polyfeaAPIService := polyfea.NewPolyfeaAPIService(
		webComponentRepository,
		microFrontendRepoistory,
		microFrontendClassRepository)

	polyfeaAPIController := generated.NewPolyfeaAPIController(polyfeaAPIService)

	router := generated.NewRouter(polyfeaAPIController)

	router.HandleFunc("/openapi", api.HandleOpenApi)

	return polyfea.BasePathStrippingMiddleware(router)
}
