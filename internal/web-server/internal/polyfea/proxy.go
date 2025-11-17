package polyfea

import (
	"context"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/internal/repository"
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
	logger := p.prepareLogger("HandleProxy", r.Method, r.URL.Path)
	ctx, span := p.startSpan(r.Context(), "polyfea_d.serve_asset", r.Method, r.URL.Path)
	defer span.End()

	params := mux.Vars(r)
	nameSpace := params[NamespacePathParamName]
	nameMicroFrontend := params[MicrofrontendPathParamName]
	path := params[PathPathParamName]

	microfrontend, err := p.getMicrofrontend(nameSpace, nameMicroFrontend, logger, span, w, ctx)
	if err != nil {
		return
	}

	microfrontendClass, err := p.getMicrofrontendClass(microfrontend, logger, span, w)
	if err != nil {
		return
	}

	proxyUrl := p.buildProxyUrl(microfrontend.Spec.Service, path)
	resp, err := p.proxyRequest(ctx, proxyUrl, r, logger, span, w)
	if err != nil {
		logger.Err(err).Msg("Error while proxying request.")
		span.SetStatus(codes.Error, "proxying_request: "+err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	p.finalizeResponse(w, resp, microfrontendClass.Spec.ExtraHeaders, ctx, span)
}

// Helper methods for HandleProxy
func (p *PolyfeaProxy) prepareLogger(functionName, method, path string) zerolog.Logger {
	return p.logger.With().Str("function", functionName).Str("method", method).Str("path", path).Logger()
}

func (p *PolyfeaProxy) startSpan(ctx context.Context, spanName, method, path string) (context.Context, trace.Span) {
	return telemetry().tracer.Start(ctx, spanName, trace.WithAttributes(
		attribute.String("path", path),
		attribute.String("method", method),
	))
}

func (p *PolyfeaProxy) getMicrofrontend(nameSpace, nameMicroFrontend string, logger zerolog.Logger, span trace.Span, w http.ResponseWriter, ctx context.Context) (*v1alpha1.MicroFrontend, error) {
	microfrontends, err := p.microfrontendRepository.List(func(mf *v1alpha1.MicroFrontend) bool {
		return mf.Namespace == nameSpace && mf.Name == nameMicroFrontend
	})

	if err != nil {
		logger.Err(err).Msg("microfrontend_repository_error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	if len(microfrontends) == 0 {
		logger.Warn().Msg("No microfrontend found for the given namespace and name.")
		span.SetStatus(codes.Error, "microfrontend_not_found")
		http.Error(w, "No microfrontend found for the given namespace and name.", http.StatusNotFound)
		telemetry().not_found.Add(ctx, 1)
		return nil, nil
	}

	return microfrontends[0], nil
}

func (p *PolyfeaProxy) getMicrofrontendClass(microfrontend *v1alpha1.MicroFrontend, logger zerolog.Logger, span trace.Span, w http.ResponseWriter) (*v1alpha1.MicroFrontendClass, error) {
	microfrontendClasses, err := p.microfrontendClassRepository.List(func(mfc *v1alpha1.MicroFrontendClass) bool {
		return mfc.Name == *microfrontend.Spec.FrontendClass
	})

	if err != nil {
		logger.Err(err).Msg("Error while getting microfrontend class from repository.")
		span.SetStatus(codes.Error, "microfrontend_class_repository_error: "+err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	if len(microfrontendClasses) == 0 {
		logger.Warn().Msg("No microfrontend class found for the given namespace and name.")
		http.Error(w, "No microfrontend class found for the given namespace and name.", http.StatusNotFound)
		span.SetStatus(codes.Error, "microfrontend_class_not_found")
		return nil, nil
	}

	return microfrontendClasses[0], nil
}

func (p *PolyfeaProxy) buildProxyUrl(service *string, path string) string {
	if (*service)[len(*service)-1] != '/' && path[0] != '/' {
		return *service + "/" + path
	}
	return *service + path
}

func (p *PolyfeaProxy) proxyRequest(ctx context.Context, proxyUrl string, r *http.Request, logger zerolog.Logger, span trace.Span, w http.ResponseWriter) (*http.Response, error) {
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
}

func (p *PolyfeaProxy) finalizeResponse(w http.ResponseWriter, resp *http.Response, extraHeaders []v1alpha1.Header, ctx context.Context, span trace.Span) {
	copyHeaders(w.Header(), resp.Header)
	copyExtraHeaders(w.Header(), extraHeaders)

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
