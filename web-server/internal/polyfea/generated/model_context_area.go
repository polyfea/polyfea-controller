// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Polyfea Browser application
 *
 * This is the OpenAPI definition for the Polyfea endpoint, which serves context information to the  browser client. The client typically requests this context information from the backend when it  encounters the `<polyfea-context>` element. This context information is then used to render the  application's UI.
 *
 * API version: v1alpha1
 */

package generated

// ContextArea - Elements to be inserted into the microfrontend context area.  The context area refers to a section in the document flow, the content of which depends  on the system's configuration. For instance, the context area `top-level-application`  could be used to render the top-level application tiles.
type ContextArea struct {

	// The elements to be incorporated into the context area.  These elements will be rendered in the sequence they appear in the array.
	Elements []ElementSpec `json:"elements"`

	// The microfrontends referenced by any of the elements. The browser triggers the loading of microfrontend resources when the element is rendered.
	Microfrontends map[string]MicrofrontendSpec `json:"microfrontends,omitempty"`
}

// AssertContextAreaRequired checks if the required fields are not zero-ed
func AssertContextAreaRequired(obj ContextArea) error {
	elements := map[string]interface{}{
		"elements": obj.Elements,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.Elements {
		if err := AssertElementSpecRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertContextAreaConstraints checks if the values respects the defined constraints
func AssertContextAreaConstraints(obj ContextArea) error {
	for _, el := range obj.Elements {
		if err := AssertElementSpecConstraints(el); err != nil {
			return err
		}
	}
	return nil
}
