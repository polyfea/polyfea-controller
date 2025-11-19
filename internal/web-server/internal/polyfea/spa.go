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

	"github.com/go-logr/logr"
	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type SingePageApplication struct {
	logger *logr.Logger
}

type TemplateData struct {
	BaseUri           string
	Title             string
	Nonce             string
	ExtraMeta         template.HTML
	EnablePWA         bool
	ReconcileInterval int32
}

//go:embed .resources/index.html
var html string

//go:embed .resources/boot.mjs
var bootJs []byte

func NewSinglePageApplication(
	logger *logr.Logger,
) *SingePageApplication {
	return &SingePageApplication{
		logger: logger,
	}
}

func (s *SingePageApplication) HandleSinglePageApplication(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.WithValues("function", "HandleSinglePageApplication", "method", r.Method, "path", r.URL.Path)

	_, span := telemetry().tracer.Start(
		r.Context(), "spa_d.serve_asset",
		trace.WithAttributes(
			attribute.String("path", r.URL.Path),
			attribute.String("method", r.Method),
		))
	defer span.End()

	basePath := r.Context().Value(PolyfeaContextKeyBasePath).(string)
	microFrontendClass := r.Context().Value(PolyfeaContextKeyMicroFrontendClass).(*v1alpha1.MicroFrontendClass)

	if microFrontendClass == nil {
		logger.Error(nil, "Microfrontend class not found")
		_, err := w.Write([]byte("Microfrontend class not found"))
		if err != nil {
			logger.Error(err, "Error while writing response")
		}
		w.WriteHeader(http.StatusNotFound)
		telemetry().not_found.Add(r.Context(), 1)
		span.SetStatus(codes.Error, "microfrontend_class_not_found")
		return
	}

	logger = logger.WithValues("base", basePath, "frontendClass", microFrontendClass.Name)
	span.SetAttributes(
		attribute.String("base", basePath),
		attribute.String("frontendClass", microFrontendClass.Name),
	)

	nonce, err := generateNonce()

	if err != nil {
		logger.Error(err, "Error while generating nonce")
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
		EnablePWA: microFrontendClass.Spec.ProgressiveWebApp != nil,
	}

	if microFrontendClass.Spec.ProgressiveWebApp != nil && microFrontendClass.Spec.ProgressiveWebApp.PolyfeaSWReconcileInterval != nil {
		templateVars.ReconcileInterval = *microFrontendClass.Spec.ProgressiveWebApp.PolyfeaSWReconcileInterval
	}

	templatedHtml, err := templateHtml(html, &templateVars)

	if err != nil {
		logger.Error(err, "Error while templating html")
		span.SetStatus(codes.Error, "template_error: "+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	cspHeader := strings.ReplaceAll(microFrontendClass.Spec.CspHeader, "{NONCE_VALUE}", nonce)

	for _, header := range microFrontendClass.Spec.ExtraHeaders {
		w.Header().Set(header.Name, header.Value)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", cspHeader)

	_, err = w.Write([]byte(templatedHtml))
	if err != nil {
		logger.Error(err, "Error while writing response")
	}

	logger.Info("Served single page application")
	telemetry().spa_served.Add(r.Context(), 1)
	span.SetStatus(codes.Ok, "served")
}

func (s *SingePageApplication) HandleBootJs(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.WithValues("function", "HandleBootJs", "method", r.Method, "path", r.URL.Path)

	_, span := telemetry().tracer.Start(
		r.Context(), "spa_d.serve_boot_js",
		trace.WithAttributes(
			attribute.String("path", r.URL.Path),
			attribute.String("method", r.Method),
		))
	defer span.End()

	microFrontendClass := r.Context().Value(PolyfeaContextKeyMicroFrontendClass).(*v1alpha1.MicroFrontendClass)

	if microFrontendClass == nil {
		logger.Info("Microfrontend class not found")
		_, err := w.Write([]byte("Microfrontend class not found"))
		if err != nil {
			logger.Error(err, "Error while writing response")
		}

		w.WriteHeader(http.StatusNotFound)
		telemetry().not_found.Add(r.Context(), 1)
		return
	}

	for _, header := range microFrontendClass.Spec.ExtraHeaders {
		w.Header().Set(header.Name, header.Value)
	}

	w.Header().Set("Content-Type", "application/javascript;")

	_, err := w.Write(bootJs)
	if err != nil {
		logger.Error(err, "Error while writing response")
	}

	logger.Info("Served boot js")
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
	randomCodes := make([]byte, 128)
	_, err := rand.Read(randomCodes)
	if err != nil {
		return "", err
	}

	text := string(randomCodes)
	nonce := base64.StdEncoding.EncodeToString([]byte(text))

	return nonce, nil
}
