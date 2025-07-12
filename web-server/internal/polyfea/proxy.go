package polyfea

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	NamespacePathParamName     = "namespace"
	MicrofrontendPathParamName = "microfrontend"
	PathPathParamName          = "path"
)

type PolyfeaProxy struct {
	microfrontendClassRepository repository.Repository[*v1alpha1.MicroFrontendClass]
	microfrontendRepository      repository.Repository[*v1alpha1.MicroFrontend]
	client                       *http.Client
	logger                       *zerolog.Logger
}

func NewPolyfeaProxy(
	microfrontendClassRepository repository.Repository[*v1alpha1.MicroFrontendClass],
	microfrontendRepository repository.Repository[*v1alpha1.MicroFrontend],
	httpClient *http.Client,
	logger *zerolog.Logger,
) *PolyfeaProxy {

	l := logger.With().Str("component", "polyfea-proxy").Logger()
	return &PolyfeaProxy{
		microfrontendClassRepository: microfrontendClassRepository,
		microfrontendRepository:      microfrontendRepository,
		client:                       httpClient,
		logger:                       &l,
	}
}

func (p *PolyfeaProxy) HandleProxy(w http.ResponseWriter, r *http.Request) {
	logger := p.logger.With().
		Str("function", "HandleProxy").
		Str("method", r.Method).
		Str("path", r.URL.Path).Logger()

	ctx, span := telemetry().tracer.Start(
		r.Context(), "polyfea_d.serve_asset",
		trace.WithAttributes(
			attribute.String("path", r.URL.Path),
			attribute.String("method", r.Method),
		))
	defer span.End()
	params := mux.Vars(r)

	nameSpace := params[NamespacePathParamName]
	nameMicroFrontend := params[MicrofrontendPathParamName]
	path := params[PathPathParamName]

	microfrontends, err := p.microfrontendRepository.GetItems(func(mf *v1alpha1.MicroFrontend) bool {
		return mf.Namespace == nameSpace && mf.Name == nameMicroFrontend
	})

	if err != nil {
		logger.Err(err).Msg("microfrontend_repository_error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(microfrontends) == 0 {
		logger.Warn().Err(err).Msg("No microfrontend found for the given namespace and name.")
		span.SetStatus(codes.Error, "microfrontend_not_found")
		http.Error(w, "No microfrontend found for the given namespace and name.", http.StatusNotFound)
		telemetry().not_found.Add(ctx, 1)
		return
	}

	microfrontend := microfrontends[0]
	logger = logger.With().
		Str("microfrontend", microfrontend.Name).
		Str("microfrontend_namespace", microfrontend.Namespace).
		Logger()
	span.SetAttributes(
		attribute.String("microfrontend", microfrontend.Name),
		attribute.String("microfrontend_namespace", microfrontend.Namespace),
	)

	microfrontendClasses, err := p.microfrontendClassRepository.GetItems(func(mfc *v1alpha1.MicroFrontendClass) bool {
		return mfc.Name == *microfrontend.Spec.FrontendClass
	})

	if err != nil {
		logger.Err(err).Msg("Error while getting microfrontend class from repository.")
		span.SetStatus(codes.Error, "microfrontend_class_repository_error: "+err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(microfrontendClasses) == 0 {
		logger.Warn().Msg("No microfrontend class found for the given namespace and name.")
		http.Error(w, "No microfrontend class found for the given namespace and name.", http.StatusNotFound)
		span.SetStatus(codes.Error, "microfrontend_class_not_found")
		return
	}

	microfrontendClass := microfrontendClasses[0]

	logger = logger.With().
		Str("microfrontend_class", microfrontendClass.Name).
		Str("microfrontend_class_namespace", microfrontendClass.Namespace).
		Logger()

	span.SetAttributes(
		attribute.String("microfrontend_class", microfrontendClass.Name),
		attribute.String("microfrontend_class_namespace", microfrontendClass.Namespace),
	)

	proxyUrl := *microfrontend.Spec.Service + path

	if (*microfrontend.Spec.Service)[len(*microfrontend.Spec.Service)-1] != '/' && path[0] != '/' {
		proxyUrl = *microfrontend.Spec.Service + "/" + path
	}

	resp, err := func() (*http.Response, error) {

		subctx, subspan := telemetry().tracer.Start(ctx, "polyfea_d.proxy_request", trace.WithAttributes(
			attribute.String("proxy_url", proxyUrl),
		))
		defer subspan.End()
		req, err := http.NewRequestWithContext(subctx, "GET", proxyUrl, r.Body)
		copyHeaders(req.Header, r.Header)

		if err != nil {
			logger.Err(err).Msg("Error while creating request.")
			subspan.SetStatus(codes.Error, "creating_request")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil, err
		}

		logger.Info().Str("proxy-url", proxyUrl).Msg("Proxying request to the module.")
		subspan.SetStatus(codes.Ok, "proxying_request")
		return p.client.Do(req)
	}()
	if err != nil {
		logger.Err(err).Msg("Error while proxying request.")
		span.SetStatus(codes.Error, "proxying_request: "+err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header)

	copyExtraHeaders(w.Header(), microfrontendClass.Spec.ExtraHeaders)

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	telemetry().proxied_resource.Add(ctx, 1)
	span.SetStatus(codes.Ok, "proxied_request")
}

func copyHeaders(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func copyExtraHeaders(dst http.Header, extraHeaders []v1alpha1.Header) {
	for _, extraHeader := range extraHeaders {
		dst.Add(extraHeader.Name, extraHeader.Value)
	}
}
