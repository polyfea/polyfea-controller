package polyfea

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
	"github.com/polyfea/polyfea-controller/web-server/internal/polyfea/generated"
)

type PolyfeaApiService struct {
	webComponentRepository       repository.PolyfeaRepository[*v1alpha1.WebComponent]
	microFrontendRepository      repository.PolyfeaRepository[*v1alpha1.MicroFrontend]
	microFrontendClassRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontendClass]
}

func NewPolyfeaAPIService(
	webComponentRepository repository.PolyfeaRepository[*v1alpha1.WebComponent],
	microFrontendRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontend],
	microFrontendClassRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontendClass]) *PolyfeaApiService {

	return &PolyfeaApiService{
		webComponentRepository:       webComponentRepository,
		microFrontendRepository:      microFrontendRepository,
		microFrontendClassRepository: microFrontendClassRepository,
	}
}

func (s *PolyfeaApiService) GetContextArea(ctx context.Context, name string, path string, take int32, headers http.Header) (generated.ImplResponse, error) {
	log.Println("GetContextArea called with name: " + name + ", path: " + path + ", take: " + string(take))
	result := generated.ContextArea{
		Elements:       []generated.ElementSpec{},
		Microfrontends: map[string]generated.MicrofrontendSpec{},
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
	frontendClasses, err := s.microFrontendClassRepository.GetItems(func(mfc *v1alpha1.MicroFrontendClass) bool {
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
		return generated.Response(http.StatusInternalServerError, "Internal Server Error"), err
	}

	if len(frontendClasses) == 0 {
		return generated.Response(http.StatusNotFound, "No frontend class found for base path "+basePath), nil
	}

	if len(frontendClasses) > 1 {
		return generated.Response(http.StatusBadRequest, "Multiple frontend classes found for base path "+basePath), nil
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
		return addExtraHeaders(generated.Response(http.StatusInternalServerError, "Internal Server Error"), frontendClass.Spec.ExtraHeaders), err
	}

	// Get webcomponents for given query and frontend class
	webComponents, err := s.webComponentRepository.GetItems(func(mf *v1alpha1.WebComponent) bool {
		return selectMatchingWebComponents(mf, name, path, userRoles) && (len(microFrontendsNamesForClass) == 0 || mf.Spec.MicroFrontend == nil || slices.Contains(microFrontendsNamesForClass, *mf.Spec.MicroFrontend))
	})

	if err != nil {
		return addExtraHeaders(generated.Response(http.StatusInternalServerError, "Internal Server Error"), frontendClass.Spec.ExtraHeaders), err
	}

	if len(webComponents) == 0 {
		log.Println("No webcomponents found for query")
		return addExtraHeaders(generated.Response(http.StatusOK, result), frontendClass.Spec.ExtraHeaders), nil
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
	allMicroFrontends, err := loadAllMicroFrontends(microFrontendsToLoad, s.microFrontendRepository, []string{})

	if err != nil {
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

	return addExtraHeaders(generated.Response(http.StatusOK, result), frontendClass.Spec.ExtraHeaders), nil
}

func (s *PolyfeaApiService) GetStaticConfig(ctx context.Context, headers http.Header) (generated.ImplResponse, error) {
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
