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
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrTypeAssertionError is thrown when type an interface does not match the asserted type
	ErrTypeAssertionError = errors.New("unable to assert type")
)

// ParsingError indicates that an error has occurred when parsing request parameters
type ParsingError struct {
	Err error
}

func (e *ParsingError) Unwrap() error {
	return e.Err
}

func (e *ParsingError) Error() string {
	return e.Err.Error()
}

// RequiredError indicates that an error has occurred when parsing request parameters
type RequiredError struct {
	Field string
}

func (e *RequiredError) Error() string {
	return fmt.Sprintf("required field '%s' is zero value.", e.Field)
}

// ErrorHandler defines the required method for handling error. You may implement it and inject this into a controller if
// you would like errors to be handled differently from the DefaultErrorHandler
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error, result *ImplResponse)

// DefaultErrorHandler defines the default logic on how to handle errors from the controller. Any errors from parsing
// request params will return a StatusBadRequest. Otherwise, the error code originating from the servicer will be used.
func DefaultErrorHandler(w http.ResponseWriter, r *http.Request, err error, result *ImplResponse) {
	if _, ok := err.(*ParsingError); ok {
		// Handle parsing errors
		EncodeJSONResponse(err.Error(), func(i int) *int { return &i }(http.StatusBadRequest), map[string][]string{}, w)
	} else if _, ok := err.(*RequiredError); ok {
		// Handle missing required errors
		EncodeJSONResponse(err.Error(), func(i int) *int { return &i }(http.StatusUnprocessableEntity), map[string][]string{}, w)
	} else {
		// Handle all other errors
		EncodeJSONResponse(err.Error(), &result.Code, result.Headers, w)
	}
}
