package polyfea

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jarcoal/httpmock"
	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
)

func TestPolyfeaProxyHandleProxyProxiesTheCallAndReturnsResult(t *testing.T) {
	// Arrange
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	requestedMicroFrontend := createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true)
	requestedMicroFrontend.Spec.Service = &[]string{"http://test-service.default.svc.cluster.local"}[0]

	testMicroFrontendRepository.StoreItem(requestedMicroFrontend)
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("other-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	testMicroFrontedClassRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontendClass]()
	testMicroFrontedClassRepository.StoreItem(createTestMicroFrontendClass("test-frontend-class", "/"))
	testMicroFrontedClassRepository.StoreItem(createTestMicroFrontendClass("other-frontend-class", "other"))

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	correctModuleRequested := false
	httpmock.RegisterResponder("GET", *requestedMicroFrontend.Spec.Service+"/test-module",
		func(req *http.Request) (*http.Response, error) {
			correctModuleRequested = true
			return &http.Response{
				Status:     strconv.Itoa(200),
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte("test-module"))),
				Header:     http.Header{},
			}, nil
		},
	)

	proxy := NewPolyfeaProxy(testMicroFrontedClassRepository, testMicroFrontendRepository, &http.Client{})

	writer := httptest.NewRecorder()

	// Act
	proxy.HandleProxy(writer, createTestRequest("default", "test-microfrontend", "/test-module"))

	// Assert
	if correctModuleRequested == false {
		t.Error("The proxy did not request the correct module.")
	}

	if writer.Code != 200 {
		t.Error("The proxy did not return the correct status code.")
	}

	if writer.Body.String() != "test-module" {
		t.Error("The proxy did not return the correct body.")
	}
}

func TestPolyfeaProxyHandleProxyReturnsErrorIfServiceIsNotFound(t *testing.T) {
	// Arrange
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	requestedMicroFrontend := createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true)
	requestedMicroFrontend.Spec.Service = &[]string{"http://test-service.default.svc.cluster.local"}[0]

	testMicroFrontendRepository.StoreItem(requestedMicroFrontend)
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("other-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	testMicroFrontedClassRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontendClass]()
	testMicroFrontedClassRepository.StoreItem(createTestMicroFrontendClass("test-frontend-class", "/"))
	testMicroFrontedClassRepository.StoreItem(createTestMicroFrontendClass("other-frontend-class", "other"))

	proxy := NewPolyfeaProxy(testMicroFrontedClassRepository, testMicroFrontendRepository, &http.Client{})

	writer := httptest.NewRecorder()

	// Act
	proxy.HandleProxy(writer, createTestRequest("default", "test-microfrontend", "/test-module"))

	// Assert
	if writer.Code != 500 {
		t.Error("The proxy did not return the correct status code.")
	}

	if !strings.Contains(writer.Body.String(), "no such host") {
		t.Error("The proxy did not return the correct body.")
	}
}

func TestPolyfeaProxyHandleProxyProxiesReturnsResultWithExtraHeaders(t *testing.T) {
	// Arrange
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	requestedMicroFrontend := createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true)
	requestedMicroFrontend.Spec.Service = &[]string{"http://test-service.default.svc.cluster.local"}[0]

	testMicroFrontendRepository.StoreItem(requestedMicroFrontend)
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("other-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	testMicroFrontedClassRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontendClass]()
	expectedFrontendClass := createTestMicroFrontendClass("test-frontend-class", "/")
	expectedFrontendClass.Spec.ExtraHeaders = []v1alpha1.Header{
		{
			Name:  "test-header",
			Value: "test-value",
		},
	}

	testMicroFrontedClassRepository.StoreItem(expectedFrontendClass)
	testMicroFrontedClassRepository.StoreItem(createTestMicroFrontendClass("other-frontend-class", "other"))

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	correctModuleRequested := false
	httpmock.RegisterResponder("GET", *requestedMicroFrontend.Spec.Service+"/test-module",
		func(req *http.Request) (*http.Response, error) {
			correctModuleRequested = true
			return &http.Response{
				Status:     strconv.Itoa(200),
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte("test-module"))),
				Header:     http.Header{},
			}, nil
		},
	)

	proxy := NewPolyfeaProxy(testMicroFrontedClassRepository, testMicroFrontendRepository, &http.Client{})

	writer := httptest.NewRecorder()

	// Act
	proxy.HandleProxy(writer, createTestRequest("default", "test-microfrontend", "/test-module"))

	// Assert
	if correctModuleRequested == false {
		t.Error("The proxy did not request the correct module.")
	}

	if writer.Code != 200 {
		t.Error("The proxy did not return the correct status code.")
	}

	if writer.Body.String() != "test-module" {
		t.Error("The proxy did not return the correct body.")
	}

	if writer.Header().Get("test-header") != "test-value" {
		t.Error("The proxy did not return the correct header.")
	}
}

func createTestRequest(namespace string, microfrontend string, path string) *http.Request {
	req, _ := http.NewRequest("GET", "/polyfea/proxy/"+namespace+"/"+microfrontend+"/"+path, io.Reader(nil))

	vars := map[string]string{
		NamespacePathParamName:     namespace,
		MicrofrontendPathParamName: microfrontend,
		PathPathParamName:          path,
	}

	req = mux.SetURLVars(req, vars)
	return req
}
