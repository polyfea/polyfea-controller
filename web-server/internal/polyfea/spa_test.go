package polyfea

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/dlclark/regexp2"
	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
	"github.com/polyfea/polyfea-controller/web-server/internal/polyfea/generated"
)

var spaTestSuite = IntegrationTestSuite{
	TestRouter: polyfeaSPAApiSetupRouter(),
	TestSet: []Test{
		{
			Name: "PolyfeaSinglePageApplicationReturnsNotFoundIfUnknownPolyfeaPathIsRequested",
			Func: PolyfeaSinglePageApplicationReturnsNotFoundIfUnknownPolyfeaPathIsRequested,
		},
		{
			Name: "PolyfeaSinglePageApplicationReturnsSuccessIfKnownPolyfeaPathIsRequested",
			Func: PolyfeaSinglePageApplicationReturnsSuccessIfKnownPolyfeaPathIsRequested,
		},
		{
			Name: "PolyfeaSinglePageApplicationReturnsTemplatedHtmlIfAnythingBesidesPolyfeaIsRequested",
			Func: PolyfeaSinglePageApplicationReturnsTemplatedHtmlIfAnythingBesidesPolyfeaIsRequested,
		},
		{
			Name: "PolyfeaSinglePageApplicationReturnsBootJsWhenRequested",
			Func: PolyfeaSinglePageApplicationReturnsBootJsWhenRequested,
		},
		{
			Name: "PolyfeaSinglePageApplicationReturnsTemplatedHtmlWithCorrectBaseHref",
			Func: PolyfeaSinglePageApplicationReturnsTemplatedHtmlWithCorrectBaseHref,
		},
	},
}

func PolyfeaSinglePageApplicationReturnsNotFoundIfUnknownPolyfeaPathIsRequested(t *testing.T) {
	// Arrange
	testServerUrl := os.Getenv(TestServerUrlName)

	// Act
	response, err := http.Get(testServerUrl + "/polyfea/unknown-path")

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status code %d, got %d", http.StatusNotFound, response.StatusCode)
	}
}

func PolyfeaSinglePageApplicationReturnsSuccessIfKnownPolyfeaPathIsRequested(t *testing.T) {
	// Arrange
	testServerUrl := os.Getenv(TestServerUrlName)

	// Act
	response, err := http.Get(testServerUrl + "/polyfea/simulate-known-route")

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusNotFound, response.StatusCode)
	}
}

func PolyfeaSinglePageApplicationReturnsTemplatedHtmlIfAnythingBesidesPolyfeaIsRequested(t *testing.T) {
	// Arrange
	testServerUrl := os.Getenv(TestServerUrlName)

	nonce, err := generateNonce()

	if err != nil {
		log.Panic(err)
	}

	// Act
	response, err := http.Get(testServerUrl + "/qweqwesop")

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	if response.Header.Get("Content-Type") != "text/html; charset=utf-8" {
		t.Fatalf("expected content type %s, got %s", "text/html; charset=utf-8", response.Header.Get("Content-Type"))
	}

	nonceRegex := regexp2.MustCompile(`'nonce-(?!{NONCE_VALUE})[^']*'`, regexp2.None)

	expectedWithoutNonce, _ := nonceRegex.Replace("default-src 'self'; font-src 'self'; script-src 'strict-dynamic' 'nonce-"+nonce+"'; worker-src 'self'; manifest-src 'self'; style-src 'self' 'strict-dynamic' 'nonce-"+nonce+"'; style-src-attr 'self' 'unsafe-inline';", "'nonce-NONCE'", -1, -1)
	gotWithoutNonce, _ := nonceRegex.Replace(response.Header.Get("Content-Security-Policy"), "'nonce-NONCE'", -1, -1)

	if expectedWithoutNonce != gotWithoutNonce {
		t.Fatalf("expected content security policy %s, got %s", "default-src 'self'; font-src 'self'; script-src 'strict-dynamic' 'nonce-"+nonce+"'; worker-src 'self'; manifest-src 'self'; style-src 'self' 'strict-dynamic' 'nonce-"+nonce+"'; style-src-attr 'self' 'unsafe-inline';", response.Header.Get("Content-Security-Policy"))
	}

	if response.Header.Get("test-header") != "test-value" {
		t.Fatalf("expected header %s, got %s", "test-value", response.Header.Get("test-header"))
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}
	bodyString := string(bodyBytes)

	if strings.Contains(bodyString, "{") != false {
		t.Fatalf("expected body to not contain %s", "{")
	}

	if strings.Contains(bodyString, "}") != false {
		t.Fatalf("expected body to not contain %s", "}")
	}
}

func PolyfeaSinglePageApplicationReturnsTemplatedHtmlWithCorrectBaseHref(t *testing.T) {
	// Arrange
	testServerUrl := os.Getenv(TestServerUrlName)

	// Act
	response, err := http.Get(testServerUrl + "/other/qweqwesop")

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	if response.Header.Get("Content-Type") != "text/html; charset=utf-8" {
		t.Fatalf("expected content type %s, got %s", "text/html; charset=utf-8", response.Header.Get("Content-Type"))
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}
	bodyString := string(bodyBytes)

	if strings.Contains(bodyString, "<base href=\"/other/\"") == false {
		t.Fatalf("expected base href not found in rendered document")
	}
}

func PolyfeaSinglePageApplicationReturnsBootJsWhenRequested(t *testing.T) {
	// Arrange
	testServerUrl := os.Getenv(TestServerUrlName)

	// Act
	response, err := http.Get(testServerUrl + "/polyfea/boot.mjs")

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	if response.Header.Get("Content-Type") != "application/javascript;" {
		t.Fatalf("expected content type %s, got %s", "application/javascript;", response.Header.Get("Content-Type"))
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}
	bodyString := string(bodyBytes)

	expectedString := string(bootJs)

	if bodyString != expectedString {
		t.Fatalf("expected body %s, got %s", expectedString, bodyString)
	}
}

func polyfeaSPAApiSetupRouter() http.Handler {
	testMicroFrontendClassRepository := repository.NewInMemoryPolyfeaRepository[*v1alpha1.MicroFrontendClass]()

	mfc := createTestMicroFrontendClass("test-frontend-class", "/")
	mfc.Spec.ExtraHeaders = []v1alpha1.Header{
		{
			Name:  "test-header",
			Value: "test-value",
		},
	}

	mfc.Spec.ExtraMetaTags = []v1alpha1.MetaTag{
		{
			Name:    "test-meta-tag",
			Content: "test-content",
		},
	}

	mfc.Spec.Title = &[]string{"Polyfea"}[0]

	mfc.Spec.CspHeader = "default-src 'self'; font-src 'self'; script-src 'strict-dynamic' 'nonce-{NONCE_VALUE}'; worker-src 'self'; manifest-src 'self'; style-src 'self' 'strict-dynamic' 'nonce-{NONCE_VALUE}'; style-src-attr 'self' 'unsafe-inline';"

	testMicroFrontendClassRepository.StoreItem(mfc)
	testMicroFrontendClassRepository.StoreItem(createTestMicroFrontendClass("other-frontend-class", "other"))

	router := generated.NewRouter()

	spa := NewSinglePageApplication(testMicroFrontendClassRepository)

	router.HandleFunc("/polyfea/simulate-known-route", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("/polyfea/boot.mjs", spa.HandleBootJs)

	router.PathPrefix("/polyfea/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	router.PathPrefix("/").HandlerFunc(spa.HandleSinglePageApplication)

	return router
}
