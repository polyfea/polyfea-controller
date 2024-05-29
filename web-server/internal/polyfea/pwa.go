package polyfea

import (
	"encoding/json"
	"net/http"
	"strings"

	_ "embed"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

//go:embed .resources/register.mjs
var register []byte

//go:embed .resources/sw.mjs
var serviceWorker []byte

type CacheRouteResponse struct {
	v1alpha1.CacheRoute
	Prefix *string `json:"prefix"`
}

type ProxyConfigResponse struct {
	PreCache []v1alpha1.PreCacheEntry `json:"precache"`
	Routes   []CacheRouteResponse     `json:"routes"`
}

type ProgressiveWebApplication struct {
	logger                  *zerolog.Logger
	microFrontendRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontend]
}

func NewProgressiveWebApplication(logger *zerolog.Logger, microFrontendRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontend]) *ProgressiveWebApplication {
	return &ProgressiveWebApplication{
		logger:                  logger,
		microFrontendRepository: microFrontendRepository,
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
	logger := pwa.logger.With().
		Str("function", "ServeServiceWorker").
		Str("method", r.Method).
		Str("path", r.URL.Path).Logger()

	_, span := telemetry().tracer.Start(
		r.Context(), "pwa_d.serve_service_worker",
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
	w.Header().Set("Content-Type", "application/javascript;")

	w.Write(serviceWorker)
}

func (pwa *ProgressiveWebApplication) ServeRegister(w http.ResponseWriter, r *http.Request) {
	logger := pwa.logger.With().
		Str("function", "ServeRegister").
		Str("method", r.Method).
		Str("path", r.URL.Path).Logger()

	_, span := telemetry().tracer.Start(
		r.Context(), "pwa_d.serve_register",
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
	w.Header().Set("Content-Type", "application/javascript;")

	w.Write(register)
}

func (pwa *ProgressiveWebApplication) ServeCaching(w http.ResponseWriter, r *http.Request) {
	logger := pwa.logger.With().
		Str("function", "ServeCaching").
		Str("method", r.Method).
		Str("path", r.URL.Path).Logger()

	_, span := telemetry().tracer.Start(
		r.Context(), "pwa_d.serve_caching",
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

	config, err := pwa.getProxyConfig(microFrontendClass)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get proxy config")
		w.WriteHeader(http.StatusInternalServerError)
		span.SetStatus(codes.Error, "proxy_config_failed")
		return
	}

	err = json.NewEncoder(w).Encode(config)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to encode JSON")
		w.WriteHeader(http.StatusInternalServerError)
		span.SetStatus(codes.Error, "json_encode_failed")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (pwa *ProgressiveWebApplication) getProxyConfig(microFrontendClass *v1alpha1.MicroFrontendClass) (*ProxyConfigResponse, error) {
	var preCache []v1alpha1.PreCacheEntry
	var routes []CacheRouteResponse

	if microFrontendClass.Spec.ProgressiveWebApp == nil {
		return nil, nil
	}

	if microFrontendClass.Spec.ProgressiveWebApp.CacheOptions == nil {
		return nil, nil
	}

	if microFrontendClass.Spec.ProgressiveWebApp.CacheOptions.CacheRoutes != nil {
		for _, route := range microFrontendClass.Spec.ProgressiveWebApp.CacheOptions.CacheRoutes {
			routes = append(routes, CacheRouteResponse{
				CacheRoute: route,
			})
		}
	}

	if microFrontendClass.Spec.ProgressiveWebApp.CacheOptions.PreCache != nil {
		for _, entry := range microFrontendClass.Spec.ProgressiveWebApp.CacheOptions.PreCache {
			preCache = append(preCache, v1alpha1.PreCacheEntry{
				URL:      entry.URL,
				Revision: entry.Revision,
			})
		}
	}

	relevantMicroFrontends, err := pwa.microFrontendRepository.GetItems(func(mf *v1alpha1.MicroFrontend) bool {
		return *mf.Spec.FrontendClass == microFrontendClass.Name && *mf.Spec.Proxy
	})

	if err != nil {
		pwa.logger.Error().Err(err).Msg("Failed to get microfrontends")
		return nil, err
	}

	for _, mf := range relevantMicroFrontends {
		if mf.Spec.CacheOptions == nil {
			continue
		}

		if mf.Spec.CacheOptions.PreCache != nil {
			for _, entry := range mf.Spec.CacheOptions.PreCache {
				preCache = append(preCache, v1alpha1.PreCacheEntry{
					URL:      buildPreCachePath(mf, *entry.URL),
					Revision: entry.Revision,
				})
			}
		}

		if mf.Spec.CacheOptions.CacheRoutes != nil {
			for _, route := range mf.Spec.CacheOptions.CacheRoutes {
				routes = append(routes, CacheRouteResponse{
					CacheRoute: route,
					Prefix:     buildPreCachePath(mf, ""),
				})
			}
		}
	}

	return &ProxyConfigResponse{
		PreCache: preCache,
		Routes:   routes,
	}, nil
}

func buildPreCachePath(mf *v1alpha1.MicroFrontend, url string) *string {
	path := strings.ReplaceAll("polyfea/proxy/"+mf.Namespace+"/"+mf.Name+"/"+url, "//", "/")
	return &path
}
