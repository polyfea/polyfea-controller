package webserver

import (
	"net/http"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
	"github.com/polyfea/polyfea-controller/web-server/api"
	"github.com/polyfea/polyfea-controller/web-server/internal/polyfea"
	"github.com/polyfea/polyfea-controller/web-server/internal/polyfea/generated"
	"github.com/rs/zerolog"
)

func SetupRouter(
	microFrontendClassRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontendClass],
	microFrontendRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontend],
	webComponentRepository repository.PolyfeaRepository[*v1alpha1.WebComponent],
	logger *zerolog.Logger,
) http.Handler {

	polyfeaAPIService := polyfea.NewPolyfeaAPIService(
		webComponentRepository,
		microFrontendRepository,
		microFrontendClassRepository,
		logger,
	)

	polyfeaAPIController := generated.NewPolyfeaAPIController(polyfeaAPIService)

	router := generated.NewRouter(polyfeaAPIController)

	router.HandleFunc("/openapi", api.HandleOpenApi)

	proxy := polyfea.NewPolyfeaProxy(microFrontendClassRepository, microFrontendRepository, &http.Client{}, logger)

	router.HandleFunc("/polyfea/proxy/{"+polyfea.NamespacePathParamName+"}/{"+polyfea.MicrofrontendPathParamName+"}/{"+polyfea.PathPathParamName+":.*}", proxy.HandleProxy)

	spa := polyfea.NewSinglePageApplication(microFrontendClassRepository, logger)

	router.HandleFunc("/polyfea/boot.mjs", spa.HandleBootJs)

	router.PathPrefix("/polyfea/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	router.PathPrefix("/").HandlerFunc(spa.HandleSinglePageApplication)

	return polyfea.BasePathStrippingMiddleware(router)
}
