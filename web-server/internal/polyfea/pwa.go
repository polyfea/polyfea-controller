package polyfea

import (
	"encoding/json"
	"net/http"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ProgressiveWebApplication struct {
	logger *zerolog.Logger
}

func NewProgressiveWebApplication(logger *zerolog.Logger) *ProgressiveWebApplication {
	return &ProgressiveWebApplication{
		logger: logger,
	}
}

func (pwa *ProgressiveWebApplication) ServeAppWebManifest(w http.ResponseWriter, r *http.Request) {
	logger := pwa.logger.With().
		Str("function", "ServeAppWebManifest").
		Str("method", r.Method).
		Str("path", r.URL.Path).Logger()

	_, span := telemetry().tracer.Start(
		r.Context(), "pwa_d.serve_web_manifest",
		trace.WithAttributes(
			attribute.String("path", r.URL.Path),
			attribute.String("method", r.Method),
		))
	defer span.End()

	basePath := r.Context().Value(PolyfeaContextKeyBasePath).(string)
	microFrontendClass := r.Context().Value(PolyfeaContextKeyMicroFrontendClass).(*v1alpha1.MicroFrontendClass)

	if microFrontendClass == nil {
		logger.Warn().Msg("Microfrontend class not found")
		w.Write([]byte("Microfrontend class not found"))
		w.WriteHeader(http.StatusNotFound)
		telemetry().not_found.Add(r.Context(), 1)
		span.SetStatus(codes.Error, "microfrontend_class_not_found")
		return
	}

	logger = logger.With().Str("base", basePath).Str("frontendClass", microFrontendClass.Name).Logger()
	span.SetAttributes(
		attribute.String("base", basePath),
		attribute.String("frontendClass", microFrontendClass.Name),
	)

	for _, header := range microFrontendClass.Spec.ExtraHeaders {
		w.Header().Set(header.Name, header.Value)
	}

	err := json.NewEncoder(w).Encode(pwa.serveAppWebManifest(microFrontendClass))
	if err != nil {
		logger.Error().Err(err).Msg("Failed to encode JSON")
		w.WriteHeader(http.StatusInternalServerError)
		span.SetStatus(codes.Error, "json_encode_failed")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (pwa *ProgressiveWebApplication) serveAppWebManifest(microFrontendClass *v1alpha1.MicroFrontendClass) *v1alpha1.WebAppManifest {
	return microFrontendClass.Spec.ProgressiveWebApp.WebAppManifest
}

func (pwa *ProgressiveWebApplication) ServeServiceWorker(w http.ResponseWriter, r *http.Request) {

}

func (pwa *ProgressiveWebApplication) ServeRegister(w http.ResponseWriter, r *http.Request) {

}

func (pwa *ProgressiveWebApplication) ServeCaching(w http.ResponseWriter, r *http.Request) {

}
