package polyfea

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/go-logr/logr"
	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/internal/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type SingePageApplication struct {
	logger                  *logr.Logger
	microFrontendRepository repository.Repository[*v1alpha1.MicroFrontend]
}

type TemplateData struct {
	BaseUri           string
	Title             string
	Nonce             string
	ExtraMeta         template.HTML
	EnablePWA         bool
	ReconcileInterval int32
	ImportMapJson     template.HTML
}

//go:embed .resources/index.html
var html string

//go:embed .resources/boot.mjs
var bootJs []byte

func NewSinglePageApplication(
	logger *logr.Logger,
	microFrontendRepository repository.Repository[*v1alpha1.MicroFrontend],
) *SingePageApplication {
	return &SingePageApplication{
		logger:                  logger,
		microFrontendRepository: microFrontendRepository,
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

	// Build merged import map from all accepted microfrontends
	importMapJson, err := s.buildImportMap(microFrontendClass, logger)
	if err != nil {
		logger.Error(err, "Error while building import map")
		span.SetStatus(codes.Error, "import_map_error: "+err.Error())
		// Continue with empty import map rather than failing
		importMapJson = "{}"
	}

	templateVars := TemplateData{
		BaseUri:       basePath,
		Title:         *microFrontendClass.Spec.Title,
		Nonce:         nonce,
		ExtraMeta:     template.HTML(extraMeta),
		EnablePWA:     microFrontendClass.Spec.ProgressiveWebApp != nil,
		ImportMapJson: template.HTML(importMapJson),
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

// buildImportMap creates a merged import map from all accepted microfrontends in the class
// First-registered wins for conflicts (based on creation timestamp)
func (s *SingePageApplication) buildImportMap(microFrontendClass *v1alpha1.MicroFrontendClass, logger logr.Logger) (string, error) {
	// Get all eligible microfrontends for this class
	microfrontends, err := s.getEligibleMicrofrontends(microFrontendClass)
	if err != nil {
		return "{}", err
	}

	// Sort by creation timestamp (oldest first) for consistent ordering
	sortedMfs := s.sortMicrofrontendsByTimestamp(microfrontends)

	// Build merged import map with first-registered-wins policy
	imports, scopes := s.mergeImportMaps(sortedMfs)

	// Convert to JSON
	return s.buildImportMapJSON(imports, scopes, len(sortedMfs), logger)
}

// getEligibleMicrofrontends returns all accepted, conflict-free microfrontends for the class
func (s *SingePageApplication) getEligibleMicrofrontends(microFrontendClass *v1alpha1.MicroFrontendClass) ([]*v1alpha1.MicroFrontend, error) {
	return s.microFrontendRepository.List(func(mf *v1alpha1.MicroFrontend) bool {
		// Check if the MicroFrontend references this class
		if mf.Spec.FrontendClass == nil || *mf.Spec.FrontendClass != microFrontendClass.Name {
			return false
		}
		// Check if the MicroFrontend is accepted by the namespace policy
		if mf.Status.FrontendClassRef == nil || !mf.Status.FrontendClassRef.Accepted {
			return false
		}
		// Only include if no import map conflicts
		return len(mf.Status.ImportMapConflicts) == 0
	})
}

// mfWithTimestamp is a helper struct for sorting microfrontends by creation time
type mfWithTimestamp struct {
	mf        *v1alpha1.MicroFrontend
	timestamp int64
}

// sortMicrofrontendsByTimestamp sorts microfrontends by creation timestamp (oldest first)
func (s *SingePageApplication) sortMicrofrontendsByTimestamp(microfrontends []*v1alpha1.MicroFrontend) []mfWithTimestamp {
	mfList := make([]mfWithTimestamp, 0, len(microfrontends))
	for _, mf := range microfrontends {
		mfList = append(mfList, mfWithTimestamp{
			mf:        mf,
			timestamp: mf.CreationTimestamp.Unix(),
		})
	}

	// Simple bubble sort (sufficient for small lists)
	for i := 0; i < len(mfList)-1; i++ {
		for j := i + 1; j < len(mfList); j++ {
			if mfList[i].timestamp > mfList[j].timestamp {
				mfList[i], mfList[j] = mfList[j], mfList[i]
			}
		}
	}

	return mfList
}

// mergeImportMaps merges import maps from sorted microfrontends with first-registered-wins
func (s *SingePageApplication) mergeImportMaps(sortedMfs []mfWithTimestamp) (map[string]string, map[string]map[string]string) {
	imports := make(map[string]string)
	scopes := make(map[string]map[string]string)

	for _, item := range sortedMfs {
		mf := item.mf
		if mf.Spec.ImportMap == nil {
			continue
		}

		s.mergeTopLevelImports(mf, mf.Spec.ImportMap.Imports, imports)
		s.mergeScopedImports(mf, mf.Spec.ImportMap.Scopes, scopes)
	}

	return imports, scopes
}

// mergeTopLevelImports merges top-level imports with first-registered-wins policy
func (s *SingePageApplication) mergeTopLevelImports(mf *v1alpha1.MicroFrontend, newImports map[string]string, imports map[string]string) {
	for specifier, path := range newImports {
		if _, exists := imports[specifier]; !exists {
			// Transform relative paths to proxy paths
			resolvedPath := s.resolveImportMapPath(mf, path)
			imports[specifier] = resolvedPath
		}
	}
}

// mergeScopedImports merges scoped imports with first-registered-wins policy
func (s *SingePageApplication) mergeScopedImports(mf *v1alpha1.MicroFrontend, newScopes map[string]map[string]string, scopes map[string]map[string]string) {
	for scope, scopeImports := range newScopes {
		if scopes[scope] == nil {
			scopes[scope] = make(map[string]string)
		}
		for specifier, path := range scopeImports {
			if _, exists := scopes[scope][specifier]; !exists {
				// Transform relative paths to proxy paths
				resolvedPath := s.resolveImportMapPath(mf, path)
				scopes[scope][specifier] = resolvedPath
			}
		}
	}
}

// resolveImportMapPath resolves an import map path, converting relative paths to proxy paths
func (s *SingePageApplication) resolveImportMapPath(mf *v1alpha1.MicroFrontend, path string) string {
	// If path is already absolute (http:// or https://), return as-is
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}

	// Check if proxy is enabled (default is true)
	proxy := true
	if mf.Spec.Proxy != nil {
		proxy = *mf.Spec.Proxy
	}

	if proxy {
		// Build proxied path
		return "./polyfea/proxy/" + mf.Namespace + "/" + mf.Name + "/" + path
	}

	// For non-proxied services, combine service URL with path
	if mf.Spec.Service != nil {
		baseUrl := mf.Spec.Service.ResolveServiceURL(mf.Namespace)
		if baseUrl != "" {
			// Handle URL joining properly
			if len(baseUrl) > 0 && baseUrl[len(baseUrl)-1] != '/' && len(path) > 0 && path[0] != '/' {
				return baseUrl + "/" + path
			}
			return baseUrl + path
		}
	}

	// Fallback to original path
	return path
}

// buildImportMapJSON builds the final JSON representation of the import map
func (s *SingePageApplication) buildImportMapJSON(
	imports map[string]string,
	scopes map[string]map[string]string,
	microfrontendCount int,
	logger logr.Logger,
) (string, error) {
	importMap := make(map[string]interface{})

	if len(imports) > 0 {
		importMap["imports"] = imports
	}
	if len(scopes) > 0 {
		importMap["scopes"] = scopes
	}

	jsonBytes, err := json.Marshal(importMap)
	if err != nil {
		return "{}", err
	}

	logger.Info("Built import map",
		"microfrontendCount", microfrontendCount,
		"importsCount", len(imports),
		"scopesCount", len(scopes))

	return string(jsonBytes), nil
}
