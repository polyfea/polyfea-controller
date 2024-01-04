/*
 * Polyfea Browser application
 *
 * This is the OpenAPI definition for the Polyfea endpoint serving the context information to the browser client. The client is requesting context information from the backend typically  when approaching the `<polyfea-context>` element. The context information is then used to render the UI of the application.
 *
 * API version: v1alpha1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package polyfea

// ElementSpec - Specification of the element Elements are the building blocks of the application. Each element shall be a web component that is rendered by the browser. When rendered in context - using e.g. `polyfea-context` element, the element  attribute ˙context` is set to the context area's name.
type ElementSpec struct {

	// The name of the microfrontend that the element belongs to. The microfrontend is loaded by the browser before the element is rendered. If not provided, then it is assumed that all resources needed by the element are already loaded by the browser before the element is rendered.
	Microfrontend string `json:"microfrontend,omitempty"`

	// The name of the element - its tag name to be put into document flow.
	TagName string `json:"tagName"`

	// Attributes of the element to be set when the element is rendered.
	Attributes map[string]string `json:"attributes,omitempty"`

	// The styles of the element. Intended primary as a fallback for specific  cases, e.g. setting CSS variables.
	Style map[string]string `json:"style,omitempty"`
}

// AssertElementSpecRequired checks if the required fields are not zero-ed
func AssertElementSpecRequired(obj ElementSpec) error {
	elements := map[string]interface{}{
		"tagName": obj.TagName,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertElementSpecConstraints checks if the values respects the defined constraints
func AssertElementSpecConstraints(obj ElementSpec) error {
	return nil
}