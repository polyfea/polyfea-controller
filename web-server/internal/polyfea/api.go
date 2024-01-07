/*
 * Polyfea Browser application
 *
 * This is the OpenAPI definition for the Polyfea endpoint, which serves context information to the  browser client. The client typically requests this context information from the backend when it  encounters the `<polyfea-context>` element. This context information is then used to render the  application's UI.
 *
 * API version: v1alpha1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package polyfea

import (
	"context"
	"net/http"
)

// PolyfeaAPIRouter defines the required methods for binding the api requests to a responses for the PolyfeaAPI
// The PolyfeaAPIRouter implementation should parse necessary information from the http request,
// pass the data to a PolyfeaAPIServicer to perform the required actions, then write the service results to the http response.
type PolyfeaAPIRouter interface {
	GetContextArea(http.ResponseWriter, *http.Request)
	GetStaticConfig(http.ResponseWriter, *http.Request)
}

// PolyfeaAPIServicer defines the api actions for the PolyfeaAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type PolyfeaAPIServicer interface {
	GetContextArea(context.Context, string, string, float32) (ImplResponse, error)
	GetStaticConfig(context.Context) (ImplResponse, error)
}
