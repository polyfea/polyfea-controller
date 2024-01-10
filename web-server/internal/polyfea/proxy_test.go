package polyfea

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/jarcoal/httpmock"
	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
	"github.com/polyfea/polyfea-controller/web-server/api"
	"github.com/polyfea/polyfea-controller/web-server/internal/polyfea/generated"
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
	// Arrange
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	requestedMicroFrontend := createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true)
	requestedMicroFrontend.Spec.Service = &[]string{"http://test-service.default.svc.cluster.local"}[0]

	testMicroFrontendRepository.StoreItem(requestedMicroFrontend)
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("other-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	testMicrofrontendClassRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontendClass]()
	testMicrofrontendClassRepository.StoreItem(createTestMicroFrontendClass("test-frontend-class", "/"))
	testMicrofrontendClassRepository.StoreItem(createTestMicroFrontendClass("other-frontend-class", "other"))

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

	proxy := NewPolyfeaProxy(testMicrofrontendClassRepository, testMicroFrontendRepository, &http.Client{})

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

	testMicrofrontendClassRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontendClass]()
	testMicrofrontendClassRepository.StoreItem(createTestMicroFrontendClass("test-frontend-class", "/"))
	testMicrofrontendClassRepository.StoreItem(createTestMicroFrontendClass("other-frontend-class", "other"))

	proxy := NewPolyfeaProxy(testMicrofrontendClassRepository, testMicroFrontendRepository, &http.Client{})

	writer := httptest.NewRecorder()

	// Act
	proxy.HandleProxy(writer, createTestRequest("default", "test-microfrontend", "/test-module"))

	// Assert
	if writer.Code != 500 {
		t.Error("The proxy did not return the correct status code.")
	}
}

func TestPolyfeaProxyHandleProxyProxiesReturnsResultWithExtraHeaders(t *testing.T) {
	// Arrange
	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()
	requestedMicroFrontend := createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true)
	requestedMicroFrontend.Spec.Service = &[]string{"http://test-service.default.svc.cluster.local"}[0]

	testMicroFrontendRepository.StoreItem(requestedMicroFrontend)
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("other-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	testMicrofrontendClassRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontendClass]()
	expectedFrontendClass := createTestMicroFrontendClass("test-frontend-class", "/")
	expectedFrontendClass.Spec.ExtraHeaders = []v1alpha1.Header{
		{
			Name:  "test-header",
			Value: "test-value",
		},
	}

	testMicrofrontendClassRepository.StoreItem(expectedFrontendClass)
	testMicrofrontendClassRepository.StoreItem(createTestMicroFrontendClass("other-frontend-class", "other"))

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

	proxy := NewPolyfeaProxy(testMicrofrontendClassRepository, testMicroFrontendRepository, &http.Client{})

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

	myHandler.HandleFunc("/test-module", func(w http.ResponseWriter, r *http.Request) {
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
	response, err := http.Get(testServerUrl + "/polyfea/proxy/default/test-microfrontend/test-module")

	// Assert
	if err != nil {
		t.Error("Unexpected error occurred while calling the proxy.")
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		t.Error("The proxy did not return the correct status code.")
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}
	bodyString := string(bodyBytes)

	if bodyString != "test-module" {
		t.Error("The proxy did not return the correct body.")
	}

	if response.Header.Get("test-header") != "test-value" {
		t.Error("The proxy did not return the correct header.")
	}
}

func polyfeaProxyApiSetupRouter() http.Handler {
	testWebComponentRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.WebComponent]()

	testWebComponentRepository.StoreItem(createTestWebComponent(
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

	testWebComponentRepository.StoreItem(createTestWebComponent(
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

	testMicroFrontendRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontend]()

	mf := createTestMicroFrontend("test-microfrontend", []string{}, "test-module", "test-frontend-class", true)
	mf.Spec.Service = &[]string{"http://localhost:5003"}[0]
	testMicroFrontendRepository.StoreItem(mf)
	testMicroFrontendRepository.StoreItem(createTestMicroFrontend("other-microfrontend", []string{}, "test-module", "test-frontend-class", true))

	testMicroFrontendClassRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontendClass]()

	mfc := createTestMicroFrontendClass("test-frontend-class", "/")
	mfc.Spec.ExtraHeaders = []v1alpha1.Header{
		{
			Name:  "test-header",
			Value: "test-value",
		},
	}
	testMicroFrontendClassRepository.StoreItem(mfc)
	testMicroFrontendClassRepository.StoreItem(createTestMicroFrontendClass("other-frontend-class", "other"))

	polyfeaAPIService := NewPolyfeaAPIService(
		testWebComponentRepository,
		testMicroFrontendRepository,
		testMicroFrontendClassRepository)

	polyfeaAPIController := generated.NewPolyfeaAPIController(polyfeaAPIService)

	router := generated.NewRouter(polyfeaAPIController)

	router.HandleFunc("/openapi", api.HandleOpenApi)

	proxy := NewPolyfeaProxy(testMicroFrontendClassRepository, testMicroFrontendRepository, &http.Client{})

	router.HandleFunc("/polyfea/proxy/{"+NamespacePathParamName+"}/{"+MicrofrontendPathParamName+"}/{"+PathPathParamName+"}", proxy.HandleProxy)

	return router
}
