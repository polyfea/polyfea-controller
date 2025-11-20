package webserver

import (
	"net/http"

	"github.com/go-logr/logr"
	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/internal/repository"
	"github.com/polyfea/polyfea-controller/internal/web-server/api"
	"github.com/polyfea/polyfea-controller/internal/web-server/internal/polyfea"
	"github.com/polyfea/polyfea-controller/internal/web-server/internal/polyfea/generated"
)

func SetupRouter(
	microFrontendClassRepository repository.Repository[*v1alpha1.MicroFrontendClass],
	microFrontendRepository repository.Repository[*v1alpha1.MicroFrontend],
	webComponentRepository repository.Repository[*v1alpha1.WebComponent],
	logger *logr.Logger,
) http.Handler {

	polyfeaAPIService := polyfea.NewPolyfeaAPIService(
		webComponentRepository,
		microFrontendRepository,
		logger,
	)

	// Create a new mux and add handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/openapi", api.HandleOpenApi)

	// Create the polyfea handler with base URL "/polyfea"
	_ = generated.HandlerFromMuxWithBaseURL(polyfeaAPIService, mux, "/polyfea")

	proxy := polyfea.NewPolyfeaProxy(microFrontendClassRepository, microFrontendRepository, &http.Client{}, logger)

	mux.HandleFunc("/polyfea/proxy/{"+polyfea.NamespacePathParamName+"}/{"+polyfea.MicrofrontendPathParamName+"}/{"+polyfea.PathPathParamName+"...}", proxy.HandleProxy)

	spa := polyfea.NewSinglePageApplication(logger)
	pwa := polyfea.NewProgressiveWebApplication(logger, microFrontendRepository)

	mux.HandleFunc("/polyfea/boot.mjs", spa.HandleBootJs)

	mux.HandleFunc("/polyfea/app.webmanifest", pwa.ServeAppWebManifest)
	mux.HandleFunc("/polyfea/register.mjs", pwa.ServeRegister)

	mux.HandleFunc("/polyfea/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	mux.HandleFunc("/sw.mjs", pwa.ServeServiceWorker)
	mux.HandleFunc("/polyfea-caching.json", pwa.ServeCaching)
	mux.HandleFunc("/", spa.HandleSinglePageApplication)

	return polyfea.BasePathStrippingMiddleware(mux, microFrontendClassRepository)
}
