package polyfea

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

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
	PreCache     []v1alpha1.PreCacheEntry `json:"precache"`
	Routes       []CacheRouteResponse     `json:"routes"`
	Interceptors []v1alpha1.SWInterceptor `json:"interceptors"`
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
	pwa.serveResource(w, r, "ServeAppWebManifest", "application/manifest+json", "no-cache", func(microFrontendClass *v1alpha1.MicroFrontendClass) ([]byte, error) {
		manifest := pwa.serveAppWebManifest(microFrontendClass)
		return json.Marshal(manifest)
	})
}

func (pwa *ProgressiveWebApplication) serveAppWebManifest(microFrontendClass *v1alpha1.MicroFrontendClass) *v1alpha1.WebAppManifest {
	return microFrontendClass.Spec.ProgressiveWebApp.WebAppManifest
}

func (pwa *ProgressiveWebApplication) ServeServiceWorker(w http.ResponseWriter, r *http.Request) {
	pwa.serveResource(w, r, "ServeServiceWorker", "application/javascript", "no-store", func(microFrontendClass *v1alpha1.MicroFrontendClass) ([]byte, error) {
		return serviceWorker, nil
	})
}

func (pwa *ProgressiveWebApplication) ServeRegister(w http.ResponseWriter, r *http.Request) {
	pwa.serveResource(w, r, "ServeRegister", "application/javascript", "no-cache", func(microFrontendClass *v1alpha1.MicroFrontendClass) ([]byte, error) {
		return register, nil
	})
}

func (pwa *ProgressiveWebApplication) ServeCaching(w http.ResponseWriter, r *http.Request) {
	pwa.serveResource(w, r, "ServeCaching", "application/json", "no-cache", func(microFrontendClass *v1alpha1.MicroFrontendClass) ([]byte, error) {
		config, err := pwa.getProxyConfig(microFrontendClass)
		if err != nil {
			return nil, err
		}
		return json.Marshal(config)
	})
}

func (pwa *ProgressiveWebApplication) serveResource(w http.ResponseWriter, r *http.Request, functionName, contentType, cacheControl string, resourceProvider func(*v1alpha1.MicroFrontendClass) ([]byte, error)) {
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
	w.Header().Set("Cache-Control", cacheControl)

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
	empty := &ProxyConfigResponse{PreCache: []v1alpha1.PreCacheEntry{}, Routes: []CacheRouteResponse{}, Interceptors: []v1alpha1.SWInterceptor{}}

	if microFrontendClass.Spec.ProgressiveWebApp == nil {
		return empty, nil
	}

	sw := microFrontendClass.Spec.ProgressiveWebApp.ServiceWorker
	if sw == nil {
		return empty, nil
	}

	relevantMicroFrontends, err := pwa.microFrontendRepository.List(func(mf *v1alpha1.MicroFrontend) bool {
		return mf.Spec.FrontendClass != nil && *mf.Spec.FrontendClass == microFrontendClass.Name &&
			mf.Spec.Proxy != nil && *mf.Spec.Proxy
	})
	if err != nil {
		pwa.logger.Error(err, "Failed to get microfrontends")
		return nil, err
	}

	preCache := pwa.collectPreCache(sw, relevantMicroFrontends)
	preCache = append(preCache, pwa.collectPreCacheByJson(sw, relevantMicroFrontends)...)
	routes := pwa.collectCacheRoutes(sw, relevantMicroFrontends)
	interceptors := pwa.collectInterceptors(sw, relevantMicroFrontends)

	return &ProxyConfigResponse{
		PreCache:     preCache,
		Routes:       routes,
		Interceptors: interceptors,
	}, nil
}

func (pwa *ProgressiveWebApplication) collectPreCache(sw *v1alpha1.ServiceWorker, microFrontends []*v1alpha1.MicroFrontend) []v1alpha1.PreCacheEntry {
	var preCache []v1alpha1.PreCacheEntry

	for _, entry := range sw.PreCache {
		preCache = append(preCache, v1alpha1.PreCacheEntry{
			URL:      entry.URL,
			Revision: entry.Revision,
		})
	}

	for _, mf := range microFrontends {
		if mf.Spec.ServiceWorker == nil {
			continue
		}
		for _, entry := range mf.Spec.ServiceWorker.PreCache {
			preCache = append(preCache, v1alpha1.PreCacheEntry{
				URL:      buildPreCachePath(mf, *entry.URL),
				Revision: entry.Revision,
			})
		}
	}

	return preCache
}

func (pwa *ProgressiveWebApplication) collectPreCacheByJson(sw *v1alpha1.ServiceWorker, microFrontends []*v1alpha1.MicroFrontend) []v1alpha1.PreCacheEntry {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var precacheMap sync.Map
	var wg sync.WaitGroup
	fetchIndex := 0

	if sw.PrecacheFromJson != "" {
		wg.Add(1)
		idx := fetchIndex
		fetchIndex++
		go func() {
			defer wg.Done()
			entries, err := fetchPrecacheFromJsonURL(ctx, sw.PrecacheFromJson)
			if err != nil {
				pwa.logger.Error(err, "Failed to fetch precache from JSON", "url", sw.PrecacheFromJson)
				return
			}
			var result []v1alpha1.PreCacheEntry
			for _, p := range entries {
				result = append(result, v1alpha1.PreCacheEntry{URL: &p})
			}
			precacheMap.Store(idx, result)
		}()
	}

	for _, mf := range microFrontends {
		if mf.Spec.ServiceWorker == nil || mf.Spec.ServiceWorker.PrecacheFromJson == "" {
			continue
		}
		wg.Add(1)
		idx := fetchIndex
		fetchIndex++
		go func() {
			defer wg.Done()
			resolvedURL := mf.Spec.Service.ResolveServiceURL(mf.Namespace) + "/" + strings.TrimLeft(mf.Spec.ServiceWorker.PrecacheFromJson, "/")
			entries, err := fetchPrecacheFromJsonURL(ctx, resolvedURL)
			if err != nil {
				pwa.logger.Error(err, "Failed to fetch precache from JSON",
					"url", resolvedURL, "microfrontend", mf.Name, "namespace", mf.Namespace)
				return
			}
			var result []v1alpha1.PreCacheEntry
			for _, p := range entries {
				resolved := buildPreCachePath(mf, p)
				result = append(result, v1alpha1.PreCacheEntry{URL: resolved})
			}
			precacheMap.Store(idx, result)
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		pwa.logger.Error(ctx.Err(), "Timeout waiting for precache JSON fetches")
	}

	var preCache []v1alpha1.PreCacheEntry
	precacheMap.Range(func(_, value any) bool {
		if entries, ok := value.([]v1alpha1.PreCacheEntry); ok {
			preCache = append(preCache, entries...)
		}
		return true
	})

	return preCache
}

func (pwa *ProgressiveWebApplication) collectCacheRoutes(sw *v1alpha1.ServiceWorker, microFrontends []*v1alpha1.MicroFrontend) []CacheRouteResponse {
	var routes []CacheRouteResponse

	for _, route := range sw.CacheRoutes {
		routes = append(routes, CacheRouteResponse{
			CacheRoute: route,
		})
	}

	for _, mf := range microFrontends {
		if mf.Spec.ServiceWorker == nil {
			continue
		}
		for _, route := range mf.Spec.ServiceWorker.CacheRoutes {
			routes = append(routes, CacheRouteResponse{
				CacheRoute: route,
				Prefix:     rebaseCacheRoute(mf, ""),
			})
		}
	}

	return routes
}

func (pwa *ProgressiveWebApplication) collectInterceptors(sw *v1alpha1.ServiceWorker, microFrontends []*v1alpha1.MicroFrontend) []v1alpha1.SWInterceptor {
	var interceptors []v1alpha1.SWInterceptor

	for _, interceptor := range sw.Interceptors {
		interceptors = append(interceptors, v1alpha1.SWInterceptor{
			Name:      interceptor.Name,
			ModuleUrl: interceptor.ModuleUrl,
			Priority:  interceptor.Priority,
		})
	}

	allowMfeInterceptors := true
	if sw.NoMicroFrontEndInterceptors != nil {
		allowMfeInterceptors = !*sw.NoMicroFrontEndInterceptors
	}
	if allowMfeInterceptors {
		for _, mf := range microFrontends {
			if mf.Spec.ServiceWorker != nil {
				for _, interceptor := range mf.Spec.ServiceWorker.Interceptors {
					resolvedModuleUrl := buildPreCachePath(mf, interceptor.ModuleUrl)
					interceptors = append(interceptors, v1alpha1.SWInterceptor{
						Name:      interceptor.Name,
						ModuleUrl: *resolvedModuleUrl,
						Priority:  interceptor.Priority,
					})
				}
			}
		}
	}

	slices.SortFunc(interceptors, func(a, b v1alpha1.SWInterceptor) int {
		var ap, bp int32
		if a.Priority != nil {
			ap = *a.Priority
		}
		if b.Priority != nil {
			bp = *b.Priority
		}
		if ap < bp {
			return 1
		} else if ap > bp {
			return -1
		}
		return 0
	})

	for i := range interceptors {
		interceptors[i].Priority = nil
	}

	return interceptors
}

func fetchPrecacheFromJsonURL(ctx context.Context, jsonUrl string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jsonUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var paths []string
	if err := json.NewDecoder(resp.Body).Decode(&paths); err != nil {
		return nil, err
	}

	return paths, nil
}

func buildPreCachePath(mf *v1alpha1.MicroFrontend, reference string) *string {
	path := "polyfea/proxy/" + mf.Namespace + "/" + mf.Name + "/" + hashOrDefault(mf.Spec.CacheBustingHash) + "/" + strings.TrimLeft(reference, "/")
	return &path
}

func rebaseCacheRoute(mf *v1alpha1.MicroFrontend, reference string) *string {
	base := &url.URL{Path: "polyfea/proxy/" + mf.Namespace + "/" + mf.Name + "/" + hashOrDefault(mf.Spec.CacheBustingHash) + "/"}
	resolved := strings.TrimPrefix(base.ResolveReference(&url.URL{Path: reference}).String(), "/") // if path then still relative to base href of mfe class
	return &resolved
}
