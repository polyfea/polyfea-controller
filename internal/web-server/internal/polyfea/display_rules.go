package polyfea

import (
	"encoding/json"
	"net/http"
	"regexp"
	"slices"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/internal/web-server/internal/polyfea/generated"
)

// selectMatchingWebComponents checks if a WebComponent matches the given display criteria.
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

// convertAttributes converts API attributes to a map for the response.
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

// convertStyles converts API styles to a map for the response.
func convertStyles(styles []v1alpha1.Style) *map[string]string {
	result := make(map[string]string)

	for _, style := range styles {
		result[style.Name] = style.Value
	}

	return &result
}

// convertMicrofrontendResources converts API resources to the response format.
func convertMicrofrontendResources(microFrontendNamespace string, microFrontendName string, resources []v1alpha1.StaticResources, service *v1alpha1.ServiceReference) *[]generated.MicrofrontendResource {
	result := make([]generated.MicrofrontendResource, 0, len(resources))

	for _, resource := range resources {
		kind := generated.MicrofrontendResourceKind(resource.Kind)
		result = append(result, generated.MicrofrontendResource{
			Kind:       &kind,
			Href:       buildModulePath(microFrontendNamespace, microFrontendName, resource.Path, *resource.Proxy, service),
			Attributes: convertAttributes(resource.Attributes),
			WaitOnLoad: &resource.WaitOnLoad,
		})
	}

	return arrToPtr(result)
}

// addExtraHeaders adds extra headers from the MicroFrontendClass spec to the response.
func addExtraHeaders(w http.ResponseWriter, extraHeaders []v1alpha1.Header) {
	for _, header := range extraHeaders {
		w.Header().Add(header.Name, header.Value)
	}
}
