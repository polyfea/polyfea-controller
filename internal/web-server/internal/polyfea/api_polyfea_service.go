package polyfea

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/go-logr/logr"
	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/internal/repository"
	"github.com/polyfea/polyfea-controller/internal/web-server/internal/polyfea/generated"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type PolyfeaApiService struct {
	webComponentRepository  repository.Repository[*v1alpha1.WebComponent]
	microFrontendRepository repository.Repository[*v1alpha1.MicroFrontend]
	logger                  *logr.Logger
}

func NewPolyfeaAPIService(
	webComponentRepository repository.Repository[*v1alpha1.WebComponent],
	microFrontendRepository repository.Repository[*v1alpha1.MicroFrontend],
	logger *logr.Logger,
) *PolyfeaApiService {

	l := logger.WithValues("component", "PolyfeaApiService")
	return &PolyfeaApiService{
		webComponentRepository:  webComponentRepository,
		microFrontendRepository: microFrontendRepository,
		logger:                  &l,
	}
}

func (s *PolyfeaApiService) GetContextArea(w http.ResponseWriter, r *http.Request, name string, params generated.GetContextAreaParams) {
	logger := s.prepareLogger("GetContextArea", name, params.Path, params.Take)
	ctx, span := s.startSpan(r.Context(), "get_context_area", name, params.Path, params.Take)
	defer span.End()

	result := s.initializeContextArea()
	basePath, frontendClass := s.extractContextValues(ctx)
	logger, span = s.updateLoggerAndSpan(logger, span, basePath)

	userRoles := s.extractUserRoles(r.Header, frontendClass)
	logger, span = s.updateLoggerAndSpanWithRoles(logger, span, frontendClass, userRoles)

	microFrontendsForClass, err := s.getMicroFrontendsForClass(frontendClass)
	if err != nil {
		s.handleRepositoryError(w, logger, span, "microfrontend_repository_error", frontendClass, err)
		return
	}

	webComponents, err := s.getWebComponents(name, params.Path, userRoles, microFrontendsForClass)
	if err != nil {
		s.handleRepositoryError(w, logger, span, "webcomponent_repository_error", frontendClass, err)
		return
	}

	if len(webComponents) == 0 {
		s.handleNoWebComponentsFound(w, logger, span, frontendClass, result, ctx)
		return
	}

	webComponents = s.limitWebComponents(webComponents, params.Take)
	// Initialize microFrontendsToLoad
	microFrontendsToLoad := []string{}
	result.Elements = s.convertWebComponentsToResponse(webComponents, &microFrontendsToLoad)

	allMicroFrontends, err := s.loadAllMicroFrontends(microFrontendsToLoad, s.microFrontendRepository, []string{})
	if err != nil {
		s.handleRepositoryError(w, logger, span, "microfrontend_load_error", frontendClass, err)
		return
	}

	result.Microfrontends = s.convertMicroFrontendsToResponse(allMicroFrontends)

	s.finalizeResponse(w, logger, span, frontendClass, result, ctx)
}

func (s *PolyfeaApiService) GetStaticConfig(w http.ResponseWriter, r *http.Request) {
	_, span := telemetry().tracer.Start(r.Context(), "get_static_config")
	defer span.End()

	s.logger.Error(nil, "Static config is not implemented by the controller", "function", "GetStaticConfig")
	span.SetStatus(codes.Error, "not_implemented")

	w.WriteHeader(http.StatusNotImplemented)
	_, err := w.Write([]byte("Not implemented"))
	if err != nil {
		s.logger.Error(err, "Error while writing response")
		span.SetStatus(codes.Error, "response_writing_error")
	}
}

// Helper methods for GetContextArea
func (s *PolyfeaApiService) prepareLogger(functionName, name, path string, take *int) logr.Logger {
	return s.logger.WithValues("function", functionName, "context-area", name, "path", path, "take", take)
}

func (s *PolyfeaApiService) startSpan(ctx context.Context, spanName, name, path string, take *int) (context.Context, trace.Span) {
	attrs := []attribute.KeyValue{
		attribute.String("context-area", name),
		attribute.String("path", path),
	}
	if take != nil {
		attrs = append(attrs, attribute.Int("take", *take))
	}
	return telemetry().tracer.Start(ctx, spanName, trace.WithAttributes(attrs...))
}

func (s *PolyfeaApiService) initializeContextArea() generated.ContextArea {
	microfrontends := make(map[string]generated.MicrofrontendSpec)
	return generated.ContextArea{
		Elements:       []generated.ElementSpec{},
		Microfrontends: &microfrontends,
	}
}

func (s *PolyfeaApiService) extractContextValues(ctx context.Context) (string, *v1alpha1.MicroFrontendClass) {
	basePath := ctx.Value(PolyfeaContextKeyBasePath).(string)
	frontendClass := ctx.Value(PolyfeaContextKeyMicroFrontendClass).(*v1alpha1.MicroFrontendClass)
	return basePath, frontendClass
}

func (s *PolyfeaApiService) updateLoggerAndSpan(logger logr.Logger, span trace.Span, basePath string) (logr.Logger, trace.Span) {
	logger = logger.WithValues("base-path", basePath)
	span.SetAttributes(attribute.String("base-path", basePath))
	return logger, span
}

func (s *PolyfeaApiService) extractUserRoles(headers http.Header, frontendClass *v1alpha1.MicroFrontendClass) []string {
	userRoleHeaders := headers.Values(frontendClass.Spec.UserRolesHeader)
	userRoles := []string{}
	for _, userRoleHeader := range userRoleHeaders {
		for _, userRole := range strings.Split(userRoleHeader, ",") {
			userRoles = append(userRoles, strings.TrimSpace(userRole))
		}
	}
	return userRoles
}

func (s *PolyfeaApiService) updateLoggerAndSpanWithRoles(logger logr.Logger, span trace.Span, frontendClass *v1alpha1.MicroFrontendClass, userRoles []string) (logr.Logger, trace.Span) {
	logger = logger.WithValues("frontend-class", frontendClass.Name, "user-roles", userRoles)
	span.SetAttributes(
		attribute.String("frontend-class", frontendClass.Name),
		attribute.String("user-roles", strings.Join(userRoles, ",")),
	)
	return logger, span
}

func (s *PolyfeaApiService) getMicroFrontendsForClass(frontendClass *v1alpha1.MicroFrontendClass) ([]*v1alpha1.MicroFrontend, error) {
	return s.microFrontendRepository.List(func(mf *v1alpha1.MicroFrontend) bool {
		return *mf.Spec.FrontendClass == frontendClass.Name
	})
}

func (s *PolyfeaApiService) getWebComponents(name, path string, userRoles []string, microFrontendsForClass []*v1alpha1.MicroFrontend) ([]*v1alpha1.WebComponent, error) {
	microFrontendsNamesForClass := []string{}
	for _, mf := range microFrontendsForClass {
		microFrontendsNamesForClass = append(microFrontendsNamesForClass, mf.Name)
	}
	return s.webComponentRepository.List(func(mf *v1alpha1.WebComponent) bool {
		return selectMatchingWebComponents(mf, name, path, userRoles) && (mf.Spec.MicroFrontend == nil || slices.Contains(microFrontendsNamesForClass, *mf.Spec.MicroFrontend))
	})
}

func (s *PolyfeaApiService) handleRepositoryError(w http.ResponseWriter, logger logr.Logger, span trace.Span, errorCode string, frontendClass *v1alpha1.MicroFrontendClass, err error) {
	logger.Error(err, "Error while processing repository")
	span.SetStatus(codes.Error, errorCode)
	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write([]byte("Internal server error"))
	if err != nil {
		logger.Error(err, "Error while writing response")
		span.SetStatus(codes.Error, "response_writing_error")
	}
	addExtraHeaders(w, frontendClass.Spec.ExtraHeaders)
}

func (s *PolyfeaApiService) handleNoWebComponentsFound(w http.ResponseWriter, logger logr.Logger, span trace.Span, frontendClass *v1alpha1.MicroFrontendClass, result generated.ContextArea, ctx context.Context) {
	logger.Info("No webcomponents found for query")
	telemetry().not_found.Add(ctx, 1)
	span.SetStatus(codes.Ok, "webcomponent_not_found")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		logger.Error(err, "Error while encoding response")
		span.SetStatus(codes.Error, "response_encoding_error")
	}
	addExtraHeaders(w, frontendClass.Spec.ExtraHeaders)
}

func (s *PolyfeaApiService) limitWebComponents(webComponents []*v1alpha1.WebComponent, take *int) []*v1alpha1.WebComponent {
	sort.Slice(webComponents, func(i, j int) bool {
		return *webComponents[i].Spec.Priority > *webComponents[j].Spec.Priority
	})
	if take != nil && *take > 0 && *take < len(webComponents) {
		return webComponents[:*take]
	}
	return webComponents
}

func (s *PolyfeaApiService) convertWebComponentsToResponse(webComponents []*v1alpha1.WebComponent, microFrontendsToLoad *[]string) []generated.ElementSpec {
	result := []generated.ElementSpec{}
	for _, webComponent := range webComponents {
		var microfrontendName string
		if webComponent.Spec.MicroFrontend != nil {
			*microFrontendsToLoad = append(*microFrontendsToLoad, *webComponent.Spec.MicroFrontend)
			microfrontendName = *webComponent.Spec.MicroFrontend
		}
		result = append(result, generated.ElementSpec{
			Microfrontend: strToPtr(microfrontendName),
			TagName:       *webComponent.Spec.Element,
			Attributes:    convertAttributes(webComponent.Spec.Attributes),
			Style:         convertStyles(webComponent.Spec.Style),
		})
	}
	return result
}

func (s *PolyfeaApiService) convertMicroFrontendsToResponse(allMicroFrontends []*v1alpha1.MicroFrontend) *map[string]generated.MicrofrontendSpec {
	result := make(map[string]generated.MicrofrontendSpec)
	for _, microFrontend := range allMicroFrontends {
		result[microFrontend.Name] = generated.MicrofrontendSpec{

			DependsOn: arrToPtr(microFrontend.Spec.DependsOn),
			Module:    buildModulePath(microFrontend.Namespace, microFrontend.Name, *microFrontend.Spec.ModulePath, *microFrontend.Spec.Proxy),
			Resources: convertMicrofrontendResources(microFrontend.Namespace, microFrontend.Name, microFrontend.Spec.StaticResources),
		}
	}
	return &result
}

func (s *PolyfeaApiService) finalizeResponse(w http.ResponseWriter, logger logr.Logger, span trace.Span, frontendClass *v1alpha1.MicroFrontendClass, result generated.ContextArea, ctx context.Context) {
	logger.Info("Context area successfully generated")
	span.SetStatus(codes.Ok, "ok")
	telemetry().context_areas.Add(ctx, 1)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(result)

	if err != nil {
		logger.Error(err, "Error while encoding response")
		span.SetStatus(codes.Error, "response_encoding_error")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Internal server error"))
		if err != nil {
			logger.Error(err, "Error while writing response")
		}
	}

	addExtraHeaders(w, frontendClass.Spec.ExtraHeaders)
}

func selectMatchingWebComponents(webComponent *v1alpha1.WebComponent, name string, path string, userRoles []string) bool {
	// Check if any of display rules matches
	for _, displayRule := range webComponent.Spec.DisplayRules {
		var pathRegex *regexp.Regexp
		selectCurrent := true

		// If any of noneOf rules matches, we can evaluate to false
		for _, matcher := range displayRule.NoneOf {
			if len(matcher.Path) != 0 {
				pathRegex = regexp.MustCompile(matcher.Path)
			}

			if len(matcher.ContextName) > 0 && matcher.ContextName == name ||
				len(matcher.Path) > 0 && pathRegex.MatchString(path) ||
				len(matcher.Role) > 0 && slices.Contains(userRoles, matcher.Role) {

				selectCurrent = false
				break
			}
		}

		if !selectCurrent {
			continue
		}

		// If any of allOf rules does not match, we can evaluate to false
		for _, matcher := range displayRule.AllOf {
			if len(matcher.Path) != 0 {
				pathRegex = regexp.MustCompile(matcher.Path)
			}

			if len(matcher.ContextName) > 0 && matcher.ContextName != name ||
				len(matcher.Path) > 0 && !pathRegex.MatchString(path) ||
				len(matcher.Role) > 0 && !slices.Contains(userRoles, matcher.Role) {

				selectCurrent = false
				break
			}
		}

		if !selectCurrent {
			continue
		}

		// If any of anyOf rules matches, we can evaluate to true therfore we need to set to false first
		if len(displayRule.AnyOf) > 0 {
			selectCurrent = false
		}

		// If any of anyOf rules matches, we can evaluate to true
		for _, matcher := range displayRule.AnyOf {
			if len(matcher.Path) != 0 {
				pathRegex = regexp.MustCompile(matcher.Path)
			}

			if len(matcher.ContextName) > 0 && matcher.ContextName == name ||
				len(matcher.Path) > 0 && pathRegex.MatchString(path) ||
				len(matcher.Role) > 0 && slices.Contains(userRoles, matcher.Role) {

				selectCurrent = true
				break
			}
		}

		// If any of display rules matches, we can evaluate to true
		if selectCurrent {
			return true
		}
	}

	return false
}

func convertAttributes(attributes []v1alpha1.Attribute) *map[string]string {
	result := make(map[string]string)
	for _, webcomponentAttribute := range attributes {
		var value string
		err := json.Unmarshal(webcomponentAttribute.Value.Raw, &value)
		if err != nil {
			value = string(webcomponentAttribute.Value.Raw)
		}
		result[webcomponentAttribute.Name] = value
	}

	return &result
}

func convertStyles(styles []v1alpha1.Style) *map[string]string {
	result := make(map[string]string)

	for _, style := range styles {
		result[style.Name] = style.Value
	}

	return &result
}

func convertMicrofrontendResources(microFrontendNamespace string, microFrontendName string, resources []v1alpha1.StaticResources) *[]generated.MicrofrontendResource {
	result := []generated.MicrofrontendResource{}

	for _, resource := range resources {
		result = append(result, generated.MicrofrontendResource{
			Kind:       (*generated.MicrofrontendResourceKind)(strToPtr(string(resource.Kind))),
			Href:       buildModulePath(microFrontendNamespace, microFrontendName, resource.Path, *resource.Proxy),
			Attributes: convertAttributes(resource.Attributes),
			WaitOnLoad: &resource.WaitOnLoad,
		})
	}

	return arrToPtr(result)
}

func (s *PolyfeaApiService) loadAllMicroFrontends(microFrontendsToLoad []string, microFrontendRepository repository.Repository[*v1alpha1.MicroFrontend], loadPath []string) ([]*v1alpha1.MicroFrontend, error) {
	logger := s.logger.WithValues("function", "loadAllMicroFrontends")

	result := []*v1alpha1.MicroFrontend{}

	for _, microFrontendName := range microFrontendsToLoad {
		if slices.Contains(loadPath, microFrontendName) {
			dependencyPath := strings.Join(append(loadPath, microFrontendName), " -> ")
			logger.Error(nil, "Circular dependency detected", "dependency-path", dependencyPath, "microfrontend", microFrontendName)
			return nil, errors.New("Circular dependency detected: " + dependencyPath)
		}

		microFrontend, err := microFrontendRepository.List(func(mf *v1alpha1.MicroFrontend) bool {
			return mf.Name == microFrontendName
		})

		if err != nil {
			return nil, err
		}

		if len(microFrontend) == 0 {
			logger.Error(nil, "Microfrontend not found", "microfrontend", microFrontendName)
			return nil, errors.New("Microfrontend " + microFrontendName + " not found")
		}

		result = append(result, microFrontend[0])

		if len(microFrontend[0].Spec.DependsOn) > 0 {
			dependsOn, err := s.loadAllMicroFrontends(microFrontend[0].Spec.DependsOn, microFrontendRepository, append(loadPath, microFrontendName))

			if err != nil {
				return nil, err
			}

			result = append(result, dependsOn...)
		}
	}

	return result, nil
}

func buildModulePath(microFrontendNamespace string, microFrontendName string, path string, proxy bool) *string {
	if proxy {
		return strToPtr("./polyfea/proxy/" + microFrontendNamespace + "/" + microFrontendName + "/" + path)
	} else {
		return strToPtr(path)
	}
}

func addExtraHeaders(w http.ResponseWriter, extraHeaders []v1alpha1.Header) {
	for _, header := range extraHeaders {
		w.Header().Add(header.Name, header.Value)
	}
}
