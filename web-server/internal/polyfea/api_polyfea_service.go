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
)

type PolyfeaApiService struct {
	webComponentRepository      repository.PolyfeaRepository[*v1alpha1.WebComponent]
	microFrontendRepository     repository.PolyfeaRepository[*v1alpha1.MicroFrontend]
	microFrontedClassRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontendClass]
}

func NewPolyfeaAPIService(
	webComponentRepository repository.PolyfeaRepository[*v1alpha1.WebComponent],
	microFrontendRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontend],
	microFrontedClassRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontendClass]) *PolyfeaApiService {

	return &PolyfeaApiService{
		webComponentRepository:      webComponentRepository,
		microFrontendRepository:     microFrontendRepository,
		microFrontedClassRepository: microFrontedClassRepository,
	}
}

func (s *PolyfeaApiService) GetContextArea(ctx context.Context, name string, path string, take int32, headers http.Header) (ImplResponse, error) {
	result := ContextArea{
		Elements:       []ElementSpec{},
		Microfrontends: map[string]MicrofrontendSpec{},
	}

	// Get base path from context or use default
	basePathValue := ctx.Value(PolyfeaContextKeyBasePath)
	var basePath string
	if basePathValue == nil {
		basePath = "/"
	} else {
		basePath = basePathValue.(string)
	}

	// Get frontend class for base path
	frontendClasses, err := s.microFrontedClassRepository.GetItems(func(mfc *v1alpha1.MicroFrontendClass) bool {
		frontendClassBasePath := *mfc.Spec.BaseUri
		if frontendClassBasePath[0] != '/' {
			frontendClassBasePath = "/" + frontendClassBasePath
		}
		if frontendClassBasePath[len(frontendClassBasePath)-1] != '/' {
			frontendClassBasePath += "/"
		}

		return basePath == frontendClassBasePath
	})

	if err != nil {
		return ImplResponse{Code: 500}, err
	}

	if len(frontendClasses) == 0 {
		return Response(404, "No frontend class found for base path "+basePath), nil
	}

	if len(frontendClasses) > 1 {
		return Response(400, "Multiple frontend classes found for base path "+basePath), nil
	}

	frontendClass := frontendClasses[0]

	// Read user roles from header
	userRoleHeaders := headers.Values(frontendClass.Spec.UserRolesHeader)
	userRoles := []string{}
	for _, userRoleHeader := range userRoleHeaders {
		for _, userRole := range strings.Split(userRoleHeader, ",") {
			userRoles = append(userRoles, strings.TrimSpace(userRole))
		}
	}

	// Get microfrontends for frontend class
	microFrontendsForClass, err := s.microFrontendRepository.GetItems(func(mf *v1alpha1.MicroFrontend) bool {
		return *mf.Spec.FrontendClass == frontendClass.Name
	})

	// Get names of microfrontends for frontend class for query
	microFrontendsNamesForClass := []string{}
	for _, mf := range microFrontendsForClass {
		microFrontendsNamesForClass = append(microFrontendsNamesForClass, mf.Name)
	}

	if err != nil {
		return ImplResponse{Code: 500}, err
	}

	if len(microFrontendsForClass) == 0 {
		return Response(404, "No microfrontends found for frontend class "+frontendClass.Name), nil
	}

	// Get webcomponents for given query and frontend class
	webComponents, err := s.webComponentRepository.GetItems(func(mf *v1alpha1.WebComponent) bool {
		return selectMatchingWebComponents(mf, name, path, userRoles) && slices.Contains(microFrontendsNamesForClass, *mf.Spec.MicroFrontend)
	})

	if err != nil {
		return ImplResponse{Code: 500}, err
	}

	if len(webComponents) == 0 {
		return ImplResponse{Code: 404, Body: "No webcomponents found based on query. Name: " + name + ", Path: " + path}, nil
	}

	sort.Slice(webComponents, func(i, j int) bool {
		return *webComponents[i].Spec.Priority > *webComponents[j].Spec.Priority
	})

	// Limit number of webcomponents by take parameter
	if take > 0 && int(take) < len(webComponents) {
		webComponents = webComponents[:take]
	}

	// Chek which microfrontends are needed for selected webcomponents and convert webcomponents to response
	microFrontendsToLoad := []string{}
	for _, webComponent := range webComponents {
		microFrontendsToLoad = append(microFrontendsToLoad, *webComponent.Spec.MicroFrontend)
		result.Elements = append(result.Elements, ElementSpec{
			Microfrontend: *webComponent.Spec.MicroFrontend,
			TagName:       *webComponent.Spec.Element,
			Attributes:    convertAttributes(webComponent.Spec.Attributes),
			Style:         convertStyles(webComponent.Spec.Style),
		})
	}

	// Load all needed microfrontends with dependencies
	allMicroFrontends, err := loadAllMicroFrontends(microFrontendsToLoad, s.microFrontendRepository, []string{})

	if err != nil {
		return ImplResponse{Code: 500}, err
	}

	if len(allMicroFrontends) == 0 {
		return ImplResponse{Code: 404, Body: "None of referenced microfrontends were found"}, nil
	}

	// Convert microfrontends to response
	for _, microFrontend := range allMicroFrontends {
		result.Microfrontends[microFrontend.Name] = MicrofrontendSpec{
			DependsOn: microFrontend.Spec.DependsOn,
			Module:    *microFrontend.Spec.ModulePath, // TODO: consider base path
			Resources: convertMicrofrontendResources(microFrontend.Spec.StaticResources),
		}
	}

	return Response(200, result), nil
}

func (s *PolyfeaApiService) GetStaticConfig(ctx context.Context, headers http.Header) (ImplResponse, error) {
	return ImplResponse{Code: 501}, nil
}

func selectMatchingWebComponents(webComponent *v1alpha1.WebComponent, name string, path string, userRoles []string) bool {
	pathRegex := regexp.MustCompile(path)

	// Check if any of display rules matches
	for _, displayRule := range webComponent.Spec.DisplayRules {
		selectCurrent := true

		// If any of noneOf rules matches, we can evaluate to false
		for _, matcher := range displayRule.NoneOf {
			if len(matcher.ContextName) > 0 && matcher.ContextName == name ||
				len(matcher.Path) > 0 && pathRegex.MatchString(matcher.Path) ||
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
			if len(matcher.ContextName) > 0 && matcher.ContextName != name ||
				len(matcher.Path) > 0 && !pathRegex.MatchString(matcher.Path) ||
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
			if len(matcher.ContextName) > 0 && matcher.ContextName == name ||
				len(matcher.Path) > 0 && pathRegex.MatchString(matcher.Path) ||
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

func convertMicrofrontendResources(resources []v1alpha1.StaticResources) []MicrofrontendResource {
	result := []MicrofrontendResource{}

	for _, resource := range resources {
		result = append(result, MicrofrontendResource{
			Kind:       resource.Kind,
			Href:       resource.Path, // TODO: consider base path
			Attributes: convertAttributes(resource.Attributes),
			WaitOnLoad: resource.WaitOnLoad,
		})
	}

	return result
}

func loadAllMicroFrontends(microFrontendsToLoad []string, microFrontendRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontend], loadPath []string) ([]*v1alpha1.MicroFrontend, error) {
	result := []*v1alpha1.MicroFrontend{}

	for _, microFrontendName := range microFrontendsToLoad {
		if slices.Contains(loadPath, microFrontendName) {
			return nil, errors.New("Circular dependency detected: " + strings.Join(append(loadPath, microFrontendName), " -> "))
		}

		microFrontend, err := microFrontendRepository.GetItems(func(mf *v1alpha1.MicroFrontend) bool {
			return mf.Name == microFrontendName
		})

		if err != nil {
			return nil, err
		}

		if len(microFrontend) == 0 {
			return nil, errors.New("Microfrontend " + microFrontendName + " not found")
		}

		result = append(result, microFrontend[0])

		if len(microFrontend[0].Spec.DependsOn) > 0 {
			dependsOn, err := loadAllMicroFrontends(microFrontend[0].Spec.DependsOn, microFrontendRepository, append(loadPath, microFrontendName))

			if err != nil {
				return nil, err
			}

			result = append(result, dependsOn...)
		}
	}

	return result, nil
}
