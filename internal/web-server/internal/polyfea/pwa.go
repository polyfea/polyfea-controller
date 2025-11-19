package polyfea

import (
	"encoding/json"
	"net/http"
	"strings"

	_ "embed"

	"github.com/go-logr/logr"
	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/internal/repository"
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
	logger                  *logr.Logger
	microFrontendRepository repository.Repository[*v1alpha1.MicroFrontend]
}

func NewProgressiveWebApplication(logger *logr.Logger, microFrontendRepository repository.Repository[*v1alpha1.MicroFrontend]) *ProgressiveWebApplication {
	return &ProgressiveWebApplication{
		logger:                  logger,
		microFrontendRepository: microFrontendRepository,
	}
}

func (pwa *ProgressiveWebApplication) ServeAppWebManifest(w http.ResponseWriter, r *http.Request) {
	pwa.serveResource(w, r, "ServeAppWebManifest", "application/manifest+json", func(microFrontendClass *v1alpha1.MicroFrontendClass) ([]byte, error) {
		manifest := pwa.serveAppWebManifest(microFrontendClass)
		return json.Marshal(manifest)
	})
}

func (pwa *ProgressiveWebApplication) serveAppWebManifest(microFrontendClass *v1alpha1.MicroFrontendClass) *v1alpha1.WebAppManifest {
	return microFrontendClass.Spec.ProgressiveWebApp.WebAppManifest
}

func (pwa *ProgressiveWebApplication) ServeServiceWorker(w http.ResponseWriter, r *http.Request) {
	pwa.serveResource(w, r, "ServeServiceWorker", "application/javascript", func(microFrontendClass *v1alpha1.MicroFrontendClass) ([]byte, error) {
		return serviceWorker, nil
	})
}

func (pwa *ProgressiveWebApplication) ServeRegister(w http.ResponseWriter, r *http.Request) {
	pwa.serveResource(w, r, "ServeRegister", "application/javascript", func(microFrontendClass *v1alpha1.MicroFrontendClass) ([]byte, error) {
		return register, nil
	})
}

func (pwa *ProgressiveWebApplication) ServeCaching(w http.ResponseWriter, r *http.Request) {
	pwa.serveResource(w, r, "ServeCaching", "application/json", func(microFrontendClass *v1alpha1.MicroFrontendClass) ([]byte, error) {
		config, err := pwa.getProxyConfig(microFrontendClass)
		if err != nil {
			return nil, err
		}
		return json.Marshal(config)
	})
}

func (pwa *ProgressiveWebApplication) serveResource(w http.ResponseWriter, r *http.Request, functionName, contentType string, resourceProvider func(*v1alpha1.MicroFrontendClass) ([]byte, error)) {
	logger := pwa.logger.WithValues("function", functionName, "method", r.Method, "path", r.URL.Path)

	_, span := telemetry().tracer.Start(
		r.Context(), "pwa_d."+functionName,
		trace.WithAttributes(
			attribute.String("path", r.URL.Path),
			attribute.String("method", r.Method),
		))
	defer span.End()

	basePath := r.Context().Value(PolyfeaContextKeyBasePath).(string)
	microFrontendClass := r.Context().Value(PolyfeaContextKeyMicroFrontendClass).(*v1alpha1.MicroFrontendClass)

	if microFrontendClass == nil {
		logger.Error(nil, "Microfrontend class not found")
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Microfrontend class not found"))
		if err != nil {
			logger.Error(err, "Failed to write response")
		}
		telemetry().not_found.Add(r.Context(), 1)
		span.SetStatus(codes.Error, "microfrontend_class_not_found")
		return
	}

	logger = logger.WithValues("base", basePath, "frontendClass", microFrontendClass.Name)
	span.SetAttributes(
		attribute.String("base", basePath),
		attribute.String("frontendClass", microFrontendClass.Name),
	)

	for _, header := range microFrontendClass.Spec.ExtraHeaders {
		w.Header().Set(header.Name, header.Value)
	}
	w.Header().Set("Content-Type", contentType)

	resource, err := resourceProvider(microFrontendClass)
	if err != nil {
		logger.Error(err, "Failed to provide resource")
		w.WriteHeader(http.StatusInternalServerError)
		span.SetStatus(codes.Error, "resource_provider_failed")
		return
	}

	_, err = w.Write(resource)
	if err != nil {
		logger.Error(err, "Failed to write response")
	}
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

	relevantMicroFrontends, err := pwa.microFrontendRepository.List(func(mf *v1alpha1.MicroFrontend) bool {
		return *mf.Spec.FrontendClass == microFrontendClass.Name && *mf.Spec.Proxy
	})

	if err != nil {
		pwa.logger.Error(err, "Failed to get microfrontends")
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
