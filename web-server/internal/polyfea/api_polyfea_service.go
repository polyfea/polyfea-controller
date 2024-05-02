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

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
	"github.com/polyfea/polyfea-controller/web-server/internal/polyfea/generated"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type PolyfeaApiService struct {
	webComponentRepository  repository.PolyfeaRepository[*v1alpha1.WebComponent]
	microFrontendRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontend]
	logger                  *zerolog.Logger
}

func NewPolyfeaAPIService(
	webComponentRepository repository.PolyfeaRepository[*v1alpha1.WebComponent],
	microFrontendRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontend],
	logger *zerolog.Logger,
) *PolyfeaApiService {

	l := logger.With().Str("component", "PolyfeaApiService").Logger()
	return &PolyfeaApiService{
		webComponentRepository:  webComponentRepository,
		microFrontendRepository: microFrontendRepository,
		logger:                  &l,
	}
}

func (s *PolyfeaApiService) GetContextArea(ctx context.Context, name string, path string, take int32, headers http.Header) (generated.ImplResponse, error) {
	logger := s.logger.With().
		Str("function", "GetContextArea").
		Str("context-area", name).
		Str("path", path).
		Int32("take", take).
		Logger()

	ctx, span := telemetry().tracer.Start(ctx, "get_context_area", trace.WithAttributes(
		attribute.String("context-area", name),
		attribute.String("path", path),
		attribute.Int("take", int(take)),
	))
	defer span.End()

	result := generated.ContextArea{
		Elements:       []generated.ElementSpec{},
		Microfrontends: map[string]generated.MicrofrontendSpec{},
	}

	// Get base path from context or use default
	basePath := ctx.Value(PolyfeaContextKeyBasePath).(string)
	frontendClass := ctx.Value(PolyfeaContextKeyMicroFrontendClass).(*v1alpha1.MicroFrontendClass)
	logger = logger.With().Str("base-path", basePath).Logger()
	span.SetAttributes(attribute.String("base-path", basePath))

	// Get frontend class for base path

	// Read user roles from header
	userRoleHeaders := headers.Values(frontendClass.Spec.UserRolesHeader)
	userRoles := []string{}
	for _, userRoleHeader := range userRoleHeaders {
		for _, userRole := range strings.Split(userRoleHeader, ",") {
			userRoles = append(userRoles, strings.TrimSpace(userRole))
		}
	}

	logger = logger.With().
		Str("frontend-class", frontendClass.Name).
		Strs("user-roles", userRoles).
		Logger()
	span.SetAttributes(
		attribute.String("frontend-class", frontendClass.Name),
		attribute.String("user-roles", strings.Join(userRoles, ",")),
	)

	// Get microfrontends for frontend class
	microFrontendsForClass, err := s.microFrontendRepository.GetItems(func(mf *v1alpha1.MicroFrontend) bool {
		return *mf.Spec.FrontendClass == frontendClass.Name
	})

	if err != nil {
		logger.Err(err).Msg("Error while getting microfrontends for frontend class")
		span.SetStatus(codes.Error, "microfrontend_repository_error")
		return addExtraHeaders(generated.Response(http.StatusInternalServerError, "Internal Server Error"), frontendClass.Spec.ExtraHeaders), err
	}

	// Get names of microfrontends for frontend class for query
	microFrontendsNamesForClass := []string{}
	for _, mf := range microFrontendsForClass {
		microFrontendsNamesForClass = append(microFrontendsNamesForClass, mf.Name)
	}

	// Get webcomponents for given query and frontend class
	webComponents, err := s.webComponentRepository.GetItems(func(mf *v1alpha1.WebComponent) bool {
		return selectMatchingWebComponents(mf, name, path, userRoles) && (len(microFrontendsNamesForClass) == 0 || mf.Spec.MicroFrontend == nil || slices.Contains(microFrontendsNamesForClass, *mf.Spec.MicroFrontend))
	})

	if err != nil {
		logger.Err(err).Msg("Error while getting webcomponents for query")
		span.SetStatus(codes.Error, "webcomponent_repository_error")
		return addExtraHeaders(generated.Response(http.StatusInternalServerError, "Internal Server Error"), frontendClass.Spec.ExtraHeaders), err
	}

	if len(webComponents) == 0 {
		logger.Info().Msg("No webcomponents found for query")
		telemetry().not_found.Add(ctx, 1)
		span.SetStatus(codes.Ok, "webcomponent_not_found") // this is possible situation when context is simply not filled
		return addExtraHeaders(generated.Response(http.StatusOK, result), frontendClass.Spec.ExtraHeaders), nil
	}

	sort.Slice(webComponents, func(i, j int) bool {
		return *webComponents[i].Spec.Priority > *webComponents[j].Spec.Priority
	})

	// Limit number of webcomponents by take parameter
	if take > 0 && int(take) < len(webComponents) {
		logger.Info().Msg("Limiting number of webcomponents by take parameter")
		webComponents = webComponents[:take]
	}

	// Chek which microfrontends are needed for selected webcomponents and convert webcomponents to response
	microFrontendsToLoad := []string{}
	for _, webComponent := range webComponents {
		var microfrontendName string

		if webComponent.Spec.MicroFrontend != nil {
			microFrontendsToLoad = append(microFrontendsToLoad, *webComponent.Spec.MicroFrontend)
			microfrontendName = *webComponent.Spec.MicroFrontend
		}

		result.Elements = append(result.Elements, generated.ElementSpec{
			Microfrontend: microfrontendName,
			TagName:       *webComponent.Spec.Element,
			Attributes:    convertAttributes(webComponent.Spec.Attributes),
			Style:         convertStyles(webComponent.Spec.Style),
		})
	}

	// Load all needed microfrontends with dependencies
	allMicroFrontends, err := s.loadAllMicroFrontends(microFrontendsToLoad, s.microFrontendRepository, []string{})

	if err != nil {
		logger.Err(err).Msg("Error while loading microfrontends")
		span.SetStatus(codes.Error, "microfrontend_load_error")
		return addExtraHeaders(generated.Response(http.StatusInternalServerError, "Internal Server Error"+err.Error()), frontendClass.Spec.ExtraHeaders), err
	}

	// Convert microfrontends to response
	for _, microFrontend := range allMicroFrontends {
		result.Microfrontends[microFrontend.Name] = generated.MicrofrontendSpec{
			DependsOn: microFrontend.Spec.DependsOn,
			Module:    buildModulePath(microFrontend.Namespace, microFrontend.Name, *microFrontend.Spec.ModulePath, *microFrontend.Spec.Proxy),
			Resources: convertMicrofrontendResources(microFrontend.Namespace, microFrontend.Name, microFrontend.Spec.StaticResources),
		}
	}

	logger.Info().Msg("Context area successfully generated")
	span.SetStatus(codes.Ok, "ok")
	telemetry().context_areas.Add(ctx, 1)
	return addExtraHeaders(generated.Response(http.StatusOK, result), frontendClass.Spec.ExtraHeaders), nil

}

func (s *PolyfeaApiService) GetStaticConfig(ctx context.Context, headers http.Header) (generated.ImplResponse, error) {
	_, span := telemetry().tracer.Start(ctx, "get_static_config")
	defer span.End()

	s.logger.Error().Str("function", "GetStaticConfig").Msg("Static config is not implemented by the controller")
	span.SetStatus(codes.Error, "not_implemented")
	return generated.Response(http.StatusNotImplemented, "Not implemented"), nil
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

func convertAttributes(attributes []v1alpha1.Attribute) map[string]string {
	result := map[string]string{}

	for _, attribute := range attributes {
		var value string
		json.Unmarshal(attribute.Value.Raw, &value)
		result[attribute.Name] = value
	}

	return result
}

func convertStyles(styles []v1alpha1.Style) map[string]string {
	result := map[string]string{}

	for _, style := range styles {
		result[style.Name] = style.Value
	}

	return result
}

func convertMicrofrontendResources(microFrontendNamespace string, microFrontendName string, resources []v1alpha1.StaticResources) []generated.MicrofrontendResource {
	result := []generated.MicrofrontendResource{}

	for _, resource := range resources {
		result = append(result, generated.MicrofrontendResource{
			Kind:       resource.Kind,
			Href:       buildModulePath(microFrontendNamespace, microFrontendName, resource.Path, *resource.Proxy),
			Attributes: convertAttributes(resource.Attributes),
			WaitOnLoad: resource.WaitOnLoad,
		})
	}

	return result
}

func (s *PolyfeaApiService) loadAllMicroFrontends(microFrontendsToLoad []string, microFrontendRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontend], loadPath []string) ([]*v1alpha1.MicroFrontend, error) {
	logger := s.logger.With().Str("function", "loadAllMicroFrontends").Logger()

	result := []*v1alpha1.MicroFrontend{}

	for _, microFrontendName := range microFrontendsToLoad {
		if slices.Contains(loadPath, microFrontendName) {
			dependencyPath := strings.Join(append(loadPath, microFrontendName), " -> ")
			logger.Error().Str("dependency-path", dependencyPath).Str("microfrontend", microFrontendName).Msg("Circular dependency detected")
			return nil, errors.New("Circular dependency detected: " + dependencyPath)
		}

		microFrontend, err := microFrontendRepository.GetItems(func(mf *v1alpha1.MicroFrontend) bool {
			return mf.Name == microFrontendName
		})

		if err != nil {
			return nil, err
		}

		if len(microFrontend) == 0 {
			logger.Error().Str("microfrontend", microFrontendName).Msg("Microfrontend not found")
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

func buildModulePath(microFrontendNamespace string, microFrontendName string, path string, proxy bool) string {
	if proxy {
		return "./polyfea/proxy/" + microFrontendNamespace + "/" + microFrontendName + "/" + path
	} else {
		return path
	}
}

func addExtraHeaders(response generated.ImplResponse, extraHeaders []v1alpha1.Header) generated.ImplResponse {
	if extraHeaders == nil {
		return response
	}

	if response.Headers == nil {
		response.Headers = map[string][]string{}
	}

	for _, header := range extraHeaders {
		response.Headers[header.Name] = []string{header.Value}
	}

	return response

}
