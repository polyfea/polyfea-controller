// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Polyfea Browser application
 *
 * This is the OpenAPI definition for the Polyfea endpoint, which serves context information to the  browser client. The client typically requests this context information from the backend when it  encounters the `<polyfea-context>` element. This context information is then used to render the  application's UI.
 *
 * API version: v1alpha1
 */

package generated

// MicrofrontendResource - The resource that the microfrontend requires. This resource could be a script,  stylesheet, or any other `link` element. The browser loads this resource when the  microfrontend is requested. The loading process can occur either synchronously or asynchronously.
type MicrofrontendResource struct {

	// The type of the resource. This could be a script, stylesheet, or any other `link` element.
	Kind string `json:"kind,omitempty"`

	// The URL of the resource. This URL is usually relative to the application's base URL  and is typically served as a subpath  of `<base_href>/polyfea/webcomponent/<microfrontend-name>/<resource-path...>`.
	Href string `json:"href,omitempty"`

	// Additional attributes to be assigned to the `link` or `script` element,  alongside the `rel` and `href` attributes.
	Attributes map[string]string `json:"attributes,omitempty"`

	// If set to `true`, the browser will complete loading the resource before it finishes  loading the microfrontend. If set to `false`, the browser will load the resource  asynchronously, allowing for continued loading and rendering in the meantime.
	WaitOnLoad bool `json:"waitOnLoad,omitempty"`
}

// AssertMicrofrontendResourceRequired checks if the required fields are not zero-ed
func AssertMicrofrontendResourceRequired(obj MicrofrontendResource) error {
	return nil
}

// AssertMicrofrontendResourceConstraints checks if the values respects the defined constraints
func AssertMicrofrontendResourceConstraints(obj MicrofrontendResource) error {
	return nil
}
