package polyfea

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/internal/repository"
	"github.com/polyfea/polyfea-controller/internal/web-server/api"
	"github.com/polyfea/polyfea-controller/internal/web-server/internal/polyfea/generated"
)

var proxyTestSuite = IntegrationTestSuite{
	TestRouter: polyfeaProxyApiSetupRouter(),
	TestSet: []Test{
		{
			Name: "PolyfeaProxyHandleProxyReturnsErrorIfServiceIsNotFound",
			Func: PolyfeaProxyHandleProxyReturnsErrorIfServiceIsNotFound,
		},
		{
			Name: "PolyfeaProxyHandleProxyProxiesReturnsResultWithExtraHeaders",
			Func: PolyfeaProxyHandleProxyProxiesReturnsResultWithExtraHeaders,
		},
	},
}

func TestPolyfeaProxyHandleProxyProxiesTheCallAndReturnsResult(t *testing.T) {
	// Test that the proxy correctly forwards the call and returns the result
	// Arrange
	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	requestedMicroFrontend := createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true)

	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test-module"))
	}))
	defer mockServer.Close()

	requestedMicroFrontend.Spec.Service = &mockServer.URL
	testMicroFrontendRepository.Store(requestedMicroFrontend)
	testMicroFrontendRepository.Store(createTestMicroFrontend("other-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	testMicrofrontendClassRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontendClass]()
	testMicrofrontendClassRepository.Store(createTestMicroFrontendClass("test-frontend-class", "/"))
	testMicrofrontendClassRepository.Store(createTestMicroFrontendClass("other-frontend-class", "other"))

	proxy := NewPolyfeaProxy(testMicrofrontendClassRepository, testMicroFrontendRepository, &http.Client{}, &logr.Logger{})

	writer := httptest.NewRecorder()

	// Act
	proxy.HandleProxy(writer, createTestRequest("default", "test-microfrontend", "/test-module"))

	// Assert
	if writer.Code != 200 {
		t.Error("The proxy did not return the correct status code.")
	}

	if writer.Body.String() != "test-module" {
		t.Error("The proxy did not return the correct body.")
	}
}

func TestPolyfeaProxyHandleProxyReturnsErrorIfServiceIsNotFound(t *testing.T) {
	// Test that the proxy returns an error if the service is not found
	// Arrange
	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	requestedMicroFrontend := createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true)
	requestedMicroFrontend.Spec.Service = &[]string{"http://test-service.default.svc.cluster.local"}[0]

	testMicroFrontendRepository.Store(requestedMicroFrontend)
	testMicroFrontendRepository.Store(createTestMicroFrontend("other-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	testMicrofrontendClassRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontendClass]()
	testMicrofrontendClassRepository.Store(createTestMicroFrontendClass("test-frontend-class", "/"))
	testMicrofrontendClassRepository.Store(createTestMicroFrontendClass("other-frontend-class", "other"))

	proxy := NewPolyfeaProxy(testMicrofrontendClassRepository, testMicroFrontendRepository, &http.Client{}, &logr.Logger{})

	writer := httptest.NewRecorder()

	// Act
	proxy.HandleProxy(writer, createTestRequest("default", "test-microfrontend", "/test-module"))

	// Assert
	if writer.Code != 500 {
		t.Error("The proxy did not return the correct status code.")
	}
}

func TestPolyfeaProxyHandleProxyProxiesReturnsResultWithExtraHeaders(t *testing.T) {
	// Test that the proxy correctly forwards the call and returns the result with extra headers
	// Arrange
	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	requestedMicroFrontend := createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true)

	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test-module"))
	}))
	defer mockServer.Close()

	requestedMicroFrontend.Spec.Service = &mockServer.URL
	testMicroFrontendRepository.Store(requestedMicroFrontend)
	testMicroFrontendRepository.Store(createTestMicroFrontend("other-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	testMicrofrontendClassRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontendClass]()
	expectedFrontendClass := createTestMicroFrontendClass("test-frontend-class", "/")
	expectedFrontendClass.Spec.ExtraHeaders = []v1alpha1.Header{
		{
			Name:  "test-header",
			Value: "test-value",
		},
	}

	testMicrofrontendClassRepository.Store(expectedFrontendClass)
	testMicrofrontendClassRepository.Store(createTestMicroFrontendClass("other-frontend-class", "other"))

	proxy := NewPolyfeaProxy(testMicrofrontendClassRepository, testMicroFrontendRepository, &http.Client{}, &logr.Logger{})

	writer := httptest.NewRecorder()

	// Act
	proxy.HandleProxy(writer, createTestRequest("default", "test-microfrontend", "/test-module"))

	// Assert
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

func PolyfeaProxyHandleProxyReturnsErrorIfServiceIsNotFound(t *testing.T) {
	// Arrange
	testServerUrl := os.Getenv(TestServerUrlName)

	// Act
	response, err := http.Get(testServerUrl + "/polyfea/proxy/default/other-microfrontend/test-module")

	// Assert
	if err != nil {
		t.Error("Unexpected error occurred while calling the proxy.")
	}
	defer response.Body.Close()

	if response.StatusCode != 500 {
		t.Error("The proxy did not return the correct status code.")
	}
}

func PolyfeaProxyHandleProxyProxiesReturnsResultWithExtraHeaders(t *testing.T) {
	// Arrange
	testServerUrl := os.Getenv(TestServerUrlName)

	myHandler := http.NewServeMux()

	myHandler.HandleFunc("/test-module.css", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("test-module"))
	})

	server := &http.Server{
		Addr:    "localhost:5003",
		Handler: myHandler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on '%s': %v\n", server.Addr, err)
		}
	}()

	for {
		resp, err := http.Get("http://" + server.Addr)
		if err == nil {
			resp.Body.Close()
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	defer server.Close()

	// Act
	response, err := http.Get(testServerUrl + "/polyfea/proxy/default/test-microfrontend/test-module.css")

	// Assert
	if err != nil {
		t.Error("Unexpected error occurred while calling the proxy.")
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		t.Error("The proxy did not return the correct status code.")
	}
	buffer := make([]byte, response.ContentLength)
	n, err := response.Body.Read(buffer)
	if err != nil && err != io.EOF {
		t.Error(err)
	}
	bodyString := string(buffer[:n])

	if bodyString != "test-module" {
		t.Error("The proxy did not return the correct body.")
	}

	if response.Header.Get("test-header") != "test-value" {
		t.Error("The proxy did not return the correct header.")
	}
}

func polyfeaProxyApiSetupRouter() http.Handler {
	testWebComponentRepository := repository.NewInMemoryRepository[*v1alpha1.WebComponent]()

	testWebComponentRepository.Store(createTestWebComponent(
		"test-name",
		"test-microfrontend",
		"test-tag-name",
		[]v1alpha1.DisplayRules{
			{
				NoneOf: []v1alpha1.Matcher{
					{
						Path: "tt-path",
					},
					{
						ContextName: "tt-name",
					},
				},
			},
		},
		&[]int32{1}[0]))

	testWebComponentRepository.Store(createTestWebComponent(
		"test-other-name",
		"other-microfrontend",
		"test-tag-name",
		[]v1alpha1.DisplayRules{
			{
				NoneOf: []v1alpha1.Matcher{
					{
						Path: "tes-path",
					},
					{
						ContextName: "tt-name",
					},
				},
				AllOf: []v1alpha1.Matcher{
					{
						Path: "test-path",
					},
					{
						ContextName: "test-name",
					},
				},
				AnyOf: []v1alpha1.Matcher{
					{
						Path: "t-path",
					},
					{
						Role: "test-role",
					},
				},
			},
		},
		&[]int32{10}[0]))

	testMicroFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()

	mf := createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true)
	mf.Spec.Service = &[]string{"http://localhost:5003"}[0]
	testMicroFrontendRepository.Store(mf)
	testMicroFrontendRepository.Store(createTestMicroFrontend("other-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	testMicroFrontendClassRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontendClass]()

	mfc := createTestMicroFrontendClass("test-frontend-class", "/")
	mfc.Spec.ExtraHeaders = []v1alpha1.Header{
		{
			Name:  "test-header",
			Value: "test-value",
		},
	}
	testMicroFrontendClassRepository.Store(mfc)
	testMicroFrontendClassRepository.Store(createTestMicroFrontendClass("other-frontend-class", "other"))

	polyfeaAPIService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		&logr.Logger{})

	polyfeaAPIController := generated.NewPolyfeaAPIController(polyfeaAPIService)

	router := generated.NewRouter(polyfeaAPIController)

	router.HandleFunc("/openapi", api.HandleOpenApi)

	proxy := NewPolyfeaProxy(testMicroFrontendClassRepository, testMicroFrontendRepository, &http.Client{}, &logr.Logger{})

	router.HandleFunc("/polyfea/proxy/{"+NamespacePathParamName+"}/{"+MicrofrontendPathParamName+"}/{"+PathPathParamName+"}", proxy.HandleProxy)

	return router
}
