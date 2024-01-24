package polyfea

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type SingePageApplication struct {
	microFrontendClassRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontendClass]
	logger                       *zerolog.Logger
}

type TemplateData struct {
	BaseUri   string
	Title     string
	Nonce     string
	ExtraMeta template.HTML
}

//go:embed .resources/index.html
var html string

//go:embed .resources/boot.mjs
var bootJs []byte

func NewSinglePageApplication(
	microFrontendClassRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontendClass],
	logger *zerolog.Logger,
) *SingePageApplication {
	return &SingePageApplication{
		microFrontendClassRepository: microFrontendClassRepository,
		logger:                       logger,
	}
}

func (s *SingePageApplication) HandleSinglePageApplication(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.With().
		Str("function", "HandleSinglePageApplication").
		Str("method", r.Method).
		Str("path", r.URL.Path).Logger()

	_, span := telemetry().tracer.Start(
		r.Context(), "spa_d.serve_asset",
		trace.WithAttributes(
			attribute.String("path", r.URL.Path),
			attribute.String("method", r.Method),
		))
	defer span.End()

	basePath, microFrontendClass, err := s.getMicrofrontendClassAndBase(r.URL.Path)

	if err != nil {
		logger.Warn().Err(err).Msg("Error while getting microfrontend and base")
		span.SetStatus(codes.Error, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

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

	nonce, err := generateNonce()

	if err != nil {
		logger.Err(err).Msg("Error while generating nonce")
		span.SetStatus(codes.Error, "nonce_error: "+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	extraMeta := ""

	for _, metaTag := range microFrontendClass.Spec.ExtraMetaTags {
		extraMeta += "<meta name=\"" + metaTag.Name + "\" content=\"" + metaTag.Content + "\" >"
	}

	templateVars := TemplateData{
		BaseUri:   basePath,
		Title:     *microFrontendClass.Spec.Title,
		Nonce:     nonce,
		ExtraMeta: template.HTML(extraMeta),
	}

	templatedHtml, err := templateHtml(html, &templateVars)

	if err != nil {
		logger.Err(err).Msg("Error while templating html")
		span.SetStatus(codes.Error, "template_error: "+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	cspHeader := strings.ReplaceAll(microFrontendClass.Spec.CspHeader, "{NONCE_VALUE}", nonce)

	for _, header := range microFrontendClass.Spec.ExtraHeaders {
		w.Header().Set(header.Name, header.Value)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", cspHeader)

	w.Write([]byte(templatedHtml))
	logger.Info().Msg("Served single page application")
	telemetry().spa_served.Add(r.Context(), 1)
	span.SetStatus(codes.Ok, "served")
}

func (s *SingePageApplication) HandleBootJs(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.With().
		Str("function", "HandleBootJs").
		Str("method", r.Method).
		Str("path", r.URL.Path).Logger()

	_, span := telemetry().tracer.Start(
		r.Context(), "spa_d.serve_boot_js",
		trace.WithAttributes(
			attribute.String("path", r.URL.Path),
			attribute.String("method", r.Method),
		))
	defer span.End()

	_, microFrontendClass, err := s.getMicrofrontendClassAndBase(r.URL.Path)

	if err != nil {
		logger.Warn().Err(err).Msg("Error while getting microfrontend and base")
		w.WriteHeader(http.StatusInternalServerError)
	}

	if microFrontendClass == nil {
		logger.Warn().Msg("Microfrontend class not found")
		w.Write([]byte("Microfrontend class not found"))
		w.WriteHeader(http.StatusNotFound)
		telemetry().not_found.Add(r.Context(), 1)
		return
	}

	for _, header := range microFrontendClass.Spec.ExtraHeaders {
		w.Header().Set(header.Name, header.Value)
	}

	w.Header().Set("Content-Type", "application/javascript;")

	w.Write(bootJs)
	logger.Info().Msg("Served boot js")
	telemetry().boot_served.Add(r.Context(), 1)
}

func templateHtml(content string, templateVars *TemplateData) (string, error) {

	tmpl, err := template.New("index.html").Parse(content)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var tmplBytes bytes.Buffer
	if err := tmpl.Execute(&tmplBytes, templateVars); err != nil {
		log.Println(err)
		return "", err
	}

	return tmplBytes.String(), nil
}

func generateNonce() (string, error) {
	codes := make([]byte, 128)
	_, err := rand.Read(codes)
	if err != nil {
		return "", err
	}

	text := string(codes)
	nonce := base64.StdEncoding.EncodeToString([]byte(text))

	return nonce, nil
}

func (s *SingePageApplication) getMicrofrontendClassAndBase(requestPath string) (string, *v1alpha1.MicroFrontendClass, error) {

	slash := func(in string) string {
		if in[len(in)-1] != '/' {
			in += "/"
		}
		if in[0] != '/' {
			in = "/" + in
		}
		return in
	}

	requestPath = slash(requestPath) // needed for user's forgotten trailing slash

	microFrontendClasses, err := s.microFrontendClassRepository.GetItems(func(mfc *v1alpha1.MicroFrontendClass) bool {
		return strings.HasPrefix(requestPath, slash(*mfc.Spec.BaseUri))
	})

	if err != nil {
		return "", nil, err
	}

	if len(microFrontendClasses) == 0 {
		return "/", nil, nil
	}

	baseHref := "/"
	longestMfc := microFrontendClasses[0]
	// find longest match
	for _, mfc := range microFrontendClasses {
		mfcBase := slash(*mfc.Spec.BaseUri)
		if len(mfcBase) > len(baseHref) {
			baseHref = mfcBase
			longestMfc = mfc
		}
	}

	return baseHref, longestMfc, nil
}
