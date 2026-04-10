package polyfea

import (
	"encoding/json"
	"net/http"
	"regexp"
	"slices"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/internal/web-server/internal/polyfea/generated"
)

// matchesMatcher evaluates a single Matcher against the given request criteria.
// All conditions present on the matcher are combined with AND semantics.
// Scalar fields (ContextName, Path, Role) and nested operators (AllOf, AnyOf, NoneOf)
// may appear together; every condition must hold for the matcher to return true.
func matchesMatcher(matcher v1alpha1.Matcher, name string, path string, userRoles []string) bool {
	// Scalar field checks
	if len(matcher.ContextName) > 0 && matcher.ContextName != name {
		return false
	}
	if len(matcher.Path) > 0 {
		pathRegex := regexp.MustCompile(matcher.Path)
		if !pathRegex.MatchString(path) {
			return false
		}
	}
	if len(matcher.Role) > 0 && !slices.Contains(userRoles, matcher.Role) {
		return false
	}

	// Nested NoneOf: if any sub-matcher matches, this matcher fails
	for _, m := range matcher.NoneOf {
		if matchesMatcher(m, name, path, userRoles) {
			return false
		}
	}

	// Nested AllOf: every sub-matcher must match
	for _, m := range matcher.AllOf {
		if !matchesMatcher(m, name, path, userRoles) {
			return false
		}
	}

	// Nested AnyOf: at least one sub-matcher must match (vacuously true when empty)
	if len(matcher.AnyOf) > 0 {
		anyMatch := false
		for _, m := range matcher.AnyOf {
			if matchesMatcher(m, name, path, userRoles) {
				anyMatch = true
				break
			}
		}
		if !anyMatch {
			return false
		}
	}

	return true
}

// selectMatchingWebComponents checks if a WebComponent matches the given display criteria.
func selectMatchingWebComponents(webComponent *v1alpha1.WebComponent, name string, path string, userRoles []string) bool {
	for _, displayRule := range webComponent.Spec.DisplayRules {
		selectCurrent := true

		// If any of noneOf rules matches, we can evaluate to false
		for _, matcher := range displayRule.NoneOf {
			if matchesMatcher(matcher, name, path, userRoles) {
				selectCurrent = false
				break
			}
		}

		if !selectCurrent {
			continue
		}

		// If any of allOf rules does not match, we can evaluate to false
		for _, matcher := range displayRule.AllOf {
			if !matchesMatcher(matcher, name, path, userRoles) {
				selectCurrent = false
				break
			}
		}

		if !selectCurrent {
			continue
		}

		// AnyOf: at least one must match (vacuously true when empty)
		if len(displayRule.AnyOf) > 0 {
			selectCurrent = false
			for _, matcher := range displayRule.AnyOf {
				if matchesMatcher(matcher, name, path, userRoles) {
					selectCurrent = true
					break
				}
			}
		}

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
func convertMicrofrontendResources(microFrontendNamespace string, microFrontendName string, resources []v1alpha1.StaticResources, service *v1alpha1.ServiceReference, cacheBustingHash string) *[]generated.MicrofrontendResource {
	result := make([]generated.MicrofrontendResource, 0, len(resources))

	for _, resource := range resources {
		kind := generated.MicrofrontendResourceKind(resource.Kind)
		result = append(result, generated.MicrofrontendResource{
			Kind:       &kind,
			Href:       buildModulePath(microFrontendNamespace, microFrontendName, resource.Path, cacheBustingHash, *resource.Proxy, service),
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
