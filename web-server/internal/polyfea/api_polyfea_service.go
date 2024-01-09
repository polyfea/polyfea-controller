package polyfea

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"slices"
	"sort"

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

	basePathValue := ctx.Value(PolyfeaContextKeyBasePath)
	var basePath string
	if basePathValue == nil {
		basePath = "/"
	} else {
		basePath = basePathValue.(string)
	}

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
	userRoleHeaders := headers[frontendClass.Spec.UserRolesHeader]

	microFrontendsForClass, err := s.microFrontendRepository.GetItems(func(mf *v1alpha1.MicroFrontend) bool {
		return *mf.Spec.FrontendClass == frontendClass.Name
	})

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

	webComponents, err := s.webComponentRepository.GetItems(func(mf *v1alpha1.WebComponent) bool {
		return selectMatchingWebComponents(mf, name, path, userRoleHeaders) && slices.Contains(microFrontendsNamesForClass, *mf.Spec.MicroFrontend)
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

	if take > 0 && int(take) < len(webComponents) {
		webComponents = webComponents[:take]
	}

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

	for _, microFrontendName := range microFrontendsToLoad {
		microFrontends, err := s.microFrontendRepository.GetItems(func(mf *v1alpha1.MicroFrontend) bool {
			return mf.Name == microFrontendName
		})

		if err != nil {
			return ImplResponse{Code: 500}, err
		}

		if len(microFrontends) == 0 {
			return ImplResponse{Code: 404, Body: "Referenced microfrontend " + microFrontendName + " not found"}, nil
		}

		result.Microfrontends[microFrontendName] = MicrofrontendSpec{
			DependsOn: microFrontends[0].Spec.DependsOn,
			Module:    *microFrontends[0].Spec.ModulePath, // TODO: consider base path
			Resources: convertMicrofrontendResources(microFrontends[0].Spec.StaticResources),
		}
	}

	return Response(200, result), nil
}

func (s *PolyfeaApiService) GetStaticConfig(ctx context.Context, headers http.Header) (ImplResponse, error) {
	return ImplResponse{Code: 501}, nil
}

func selectMatchingWebComponents(webComponent *v1alpha1.WebComponent, name string, path string, userRoleHeaders []string) bool {
	pathRegex := regexp.MustCompile(path)
	selectCurrent := true
	for _, displayRule := range webComponent.Spec.DisplayRules {

		for _, matcher := range displayRule.NoneOf {
			if len(matcher.ContextName) > 0 && matcher.ContextName == name ||
				len(matcher.Path) > 0 && pathRegex.MatchString(matcher.Path) ||
				len(matcher.Role) > 0 && slices.Contains(userRoleHeaders, matcher.Role) {

				selectCurrent = false
				break
			}
		}

		if !selectCurrent {
			continue
		}

		for _, matcher := range displayRule.AllOf {
			if len(matcher.ContextName) > 0 && matcher.ContextName != name ||
				len(matcher.Path) > 0 && !pathRegex.MatchString(matcher.Path) ||
				len(matcher.Role) > 0 && !slices.Contains(userRoleHeaders, matcher.Role) {

				selectCurrent = false
				break
			}
		}

		if !selectCurrent {
			continue
		}

		if len(displayRule.AnyOf) > 0 {
			selectCurrent = false
		}

		for _, matcher := range displayRule.AnyOf {
			if len(matcher.ContextName) > 0 && matcher.ContextName == name ||
				len(matcher.Path) > 0 && pathRegex.MatchString(matcher.Path) ||
				len(matcher.Role) > 0 && slices.Contains(userRoleHeaders, matcher.Role) {

				selectCurrent = true
				break
			}
		}

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
