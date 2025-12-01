package polyfea

import (
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/internal/repository"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	defer func() {
		err := response.Body.Close()
		if err != nil {
			t.Errorf("Expected no error on closing response body, got %v", err)
		}
	}()

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
	defer func() {
		err := response.Body.Close()
		if err != nil {
			t.Errorf("Expected no error on closing response body, got %v", err)
		}
	}()

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
	defer func() {
		err := response.Body.Close()
		if err != nil {
			t.Errorf("Expected no error on closing response body, got %v", err)
		}
	}()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	if response.Header.Get("Content-Type") != "text/html; charset=utf-8" {
		t.Fatalf("expected content type %s, got %s", "text/html; charset=utf-8", response.Header.Get("Content-Type"))
	}

	nonceRegex := regexp.MustCompile(`'nonce-[^']*'`)

	expectedWithoutNonce := nonceRegex.ReplaceAllString("default-src 'self'; font-src 'self'; script-src 'strict-dynamic' 'nonce-"+nonce+"'; worker-src 'self'; manifest-src 'self'; style-src 'self' 'strict-dynamic' 'nonce-"+nonce+"'; style-src-attr 'self' 'unsafe-inline';", "'nonce-NONCE'")
	gotWithoutNonce := nonceRegex.ReplaceAllString(response.Header.Get("Content-Security-Policy"), "'nonce-NONCE'")

	if expectedWithoutNonce != gotWithoutNonce {
		t.Fatalf("expected content security policy %s, got %s", "default-src 'self'; font-src 'self'; script-src 'strict-dynamic' 'nonce-"+nonce+"'; worker-src 'self'; manifest-src 'self'; style-src 'self' 'strict-dynamic' 'nonce-"+nonce+"'; style-src-attr 'self' 'unsafe-inline';", response.Header.Get("Content-Security-Policy"))
	}

	if response.Header.Get("test-header") != TestHeaderValue {
		t.Fatalf("expected header %s, got %s", TestHeaderValue, response.Header.Get("test-header"))
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}
	bodyString := string(bodyBytes)

	// Check that template placeholders are replaced (not just any braces, since import maps contain JSON)
	if strings.Contains(bodyString, "{{") {
		t.Fatalf("expected body to not contain template placeholder %s", "{{")
	}

	if strings.Contains(bodyString, "}}") {
		t.Fatalf("expected body to not contain template placeholder %s", "}}")
	}

	if !strings.Contains(bodyString, "webmanifest") {
		t.Fatalf("expected body to  contain %s", "webmanifest")
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
	defer func() {
		err := response.Body.Close()
		if err != nil {
			t.Errorf("Expected no error on closing response body, got %v", err)
		}
	}()

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
	mfc := createTestMicroFrontendClass("test-frontend-class", "/")
	mfc.Spec.ExtraHeaders = []v1alpha1.Header{
		{
			Name:  "test-header",
			Value: TestHeaderValue,
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

	mux := http.NewServeMux()

	mfc.Spec.ProgressiveWebApp = &v1alpha1.ProgressiveWebApp{
		WebAppManifest: &v1alpha1.WebAppManifest{
			Name: &[]string{"Test"}[0],
			Icons: []v1alpha1.PWAIcon{
				{
					Type:  &[]string{"image/png"}[0],
					Sizes: &[]string{"192x192"}[0],
					Src:   &[]string{"icon.png"}[0],
				},
			},
			StartUrl: &[]string{"/"}[0],
			Display:  &[]string{"standalone"}[0],
		},
	}

	microFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	spa := NewSinglePageApplication(&logr.Logger{}, microFrontendRepository)

	mux.HandleFunc("/polyfea/simulate-known-route", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/polyfea/boot.mjs", spa.HandleBootJs)

	mux.HandleFunc("/polyfea/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	mux.HandleFunc("/", spa.HandleSinglePageApplication)

	return addDummyMiddleware(mux, "/", mfc)
}

func TestBuildImportMapWithNoMicrofrontends(t *testing.T) {
	// Arrange
	logger := logr.Logger{}
	microFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	spa := NewSinglePageApplication(&logger, microFrontendRepository)

	mfc := &v1alpha1.MicroFrontendClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-class",
			Namespace: "default",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			BaseUri: ptr("/"),
			Title:   ptr("Test"),
		},
	}

	req, _ := http.NewRequest("GET", "/", nil)

	// Act
	result, err := spa.buildImportMap(req, mfc, logger)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != "{}" {
		t.Fatalf("expected empty import map {}, got %s", result)
	}
}

func TestBuildImportMapWithSingleMicrofrontend(t *testing.T) {
	// Arrange
	logger := logr.Logger{}
	microFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()

	className := "test-class"
	mf := &v1alpha1.MicroFrontend{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "mf1",
			Namespace:         "default",
			CreationTimestamp: metav1.Now(),
		},
		Spec: v1alpha1.MicroFrontendSpec{
			FrontendClass: &className,
			Service:       &v1alpha1.ServiceReference{URI: ptr("https://example.com")},
			ModulePath:    ptr("app.js"),
			ImportMap: &v1alpha1.ImportMap{
				Imports: map[string]string{
					"react":     "./react.js",
					"react-dom": "./react-dom.js",
				},
			},
		},
		Status: v1alpha1.MicroFrontendStatus{
			FrontendClassRef: &v1alpha1.MicroFrontendClassReference{
				Name:     className,
				Accepted: true,
			},
		},
	}

	_ = microFrontendRepository.Store(mf)
	spa := NewSinglePageApplication(&logger, microFrontendRepository)

	mfc := &v1alpha1.MicroFrontendClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      className,
			Namespace: "default",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			BaseUri: ptr("/"),
			Title:   ptr("Test"),
		},
	}

	req, _ := http.NewRequest("GET", "/", nil)

	// Act
	result, err := spa.buildImportMap(req, mfc, logger)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(result, `"react"`) {
		t.Fatalf("expected import map to contain react, got %s", result)
	}
	if !strings.Contains(result, `"react-dom"`) {
		t.Fatalf("expected import map to contain react-dom, got %s", result)
	}
	if !strings.Contains(result, `"./react.js"`) {
		t.Fatalf("expected import map to contain ./react.js, got %s", result)
	}
}

func TestBuildImportMapFirstRegisteredWins(t *testing.T) {
	// Arrange
	logger := logr.Logger{}
	microFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()

	className := "test-class"

	// First microfrontend (older)
	now := metav1.Now()
	mf1 := &v1alpha1.MicroFrontend{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "mf1",
			Namespace:         "default",
			CreationTimestamp: now,
		},
		Spec: v1alpha1.MicroFrontendSpec{
			FrontendClass: &className,
			Service:       &v1alpha1.ServiceReference{URI: ptr("https://example.com")},
			ModulePath:    ptr("app.js"),
			ImportMap: &v1alpha1.ImportMap{
				Imports: map[string]string{
					"react": "./react-v18.js",
				},
			},
		},
		Status: v1alpha1.MicroFrontendStatus{
			FrontendClassRef: &v1alpha1.MicroFrontendClassReference{
				Name:     className,
				Accepted: true,
			},
		},
	}

	// Second microfrontend (newer) - should not override
	later := metav1.NewTime(now.Add(time.Minute))
	mf2 := &v1alpha1.MicroFrontend{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "mf2",
			Namespace:         "default",
			CreationTimestamp: later,
		},
		Spec: v1alpha1.MicroFrontendSpec{
			FrontendClass: &className,
			Service:       &v1alpha1.ServiceReference{URI: ptr("https://example.com")},
			ModulePath:    ptr("app.js"),
			ImportMap: &v1alpha1.ImportMap{
				Imports: map[string]string{
					"react": "./react-v17.js", // Should be ignored
					"vue":   "./vue.js",       // Should be included
				},
			},
		},
		Status: v1alpha1.MicroFrontendStatus{
			FrontendClassRef: &v1alpha1.MicroFrontendClassReference{
				Name:     className,
				Accepted: true,
			},
		},
	}

	_ = microFrontendRepository.Store(mf1)
	_ = microFrontendRepository.Store(mf2)
	spa := NewSinglePageApplication(&logger, microFrontendRepository)

	mfc := &v1alpha1.MicroFrontendClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      className,
			Namespace: "default",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			BaseUri: ptr("/"),
			Title:   ptr("Test"),
		},
	}

	req, _ := http.NewRequest("GET", "/", nil)

	// Act
	result, err := spa.buildImportMap(req, mfc, logger)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should use mf1's version (first registered)
	if !strings.Contains(result, `"./react-v18.js"`) {
		t.Fatalf("expected import map to contain ./react-v18.js (first registered), got %s", result)
	}

	// Should NOT contain mf2's conflicting version
	if strings.Contains(result, `"./react-v17.js"`) {
		t.Fatalf("expected import map to NOT contain ./react-v17.js (conflicting), got %s", result)
	}

	// Should include non-conflicting entry from mf2
	if !strings.Contains(result, `"vue"`) {
		t.Fatalf("expected import map to contain vue from mf2, got %s", result)
	}
}

func TestBuildImportMapWithScopes(t *testing.T) {
	// Arrange
	logger := logr.Logger{}
	microFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()

	className := "test-class"
	mf := &v1alpha1.MicroFrontend{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "mf1",
			Namespace:         "default",
			CreationTimestamp: metav1.Now(),
		},
		Spec: v1alpha1.MicroFrontendSpec{
			FrontendClass: &className,
			Service:       &v1alpha1.ServiceReference{URI: ptr("https://example.com")},
			ModulePath:    ptr("app.js"),
			ImportMap: &v1alpha1.ImportMap{
				Imports: map[string]string{
					"react": "./react-v18.js",
				},
				Scopes: map[string]map[string]string{
					"/legacy/": {
						"react": "./react-v16.js",
					},
				},
			},
		},
		Status: v1alpha1.MicroFrontendStatus{
			FrontendClassRef: &v1alpha1.MicroFrontendClassReference{
				Name:     className,
				Accepted: true,
			},
		},
	}

	_ = microFrontendRepository.Store(mf)
	spa := NewSinglePageApplication(&logger, microFrontendRepository)

	mfc := &v1alpha1.MicroFrontendClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      className,
			Namespace: "default",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			BaseUri: ptr("/"),
			Title:   ptr("Test"),
		},
	}

	req, _ := http.NewRequest("GET", "/", nil)

	// Act
	result, err := spa.buildImportMap(req, mfc, logger)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(result, `"scopes"`) {
		t.Fatalf("expected import map to contain scopes, got %s", result)
	}
	if !strings.Contains(result, `"/legacy/"`) {
		t.Fatalf("expected import map to contain /legacy/ scope, got %s", result)
	}
	if !strings.Contains(result, `"./react-v16.js"`) {
		t.Fatalf("expected import map to contain ./react-v16.js in scope, got %s", result)
	}
}

func TestBuildImportMapExcludesConflictedMicrofrontends(t *testing.T) {
	// Arrange
	logger := logr.Logger{}
	microFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()

	className := "test-class"

	// First microfrontend (no conflicts)
	mf1 := &v1alpha1.MicroFrontend{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "mf1",
			Namespace:         "default",
			CreationTimestamp: metav1.Now(),
		},
		Spec: v1alpha1.MicroFrontendSpec{
			FrontendClass: &className,
			Service:       &v1alpha1.ServiceReference{URI: ptr("https://example.com")},
			ModulePath:    ptr("app.js"),
			ImportMap: &v1alpha1.ImportMap{
				Imports: map[string]string{
					"react": "./react.js",
				},
			},
		},
		Status: v1alpha1.MicroFrontendStatus{
			FrontendClassRef: &v1alpha1.MicroFrontendClassReference{
				Name:     className,
				Accepted: true,
			},
		},
	}

	// Second microfrontend (has conflicts - should be excluded)
	mf2 := &v1alpha1.MicroFrontend{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "mf2",
			Namespace:         "default",
			CreationTimestamp: metav1.Now(),
		},
		Spec: v1alpha1.MicroFrontendSpec{
			FrontendClass: &className,
			Service:       &v1alpha1.ServiceReference{URI: ptr("https://example.com")},
			ModulePath:    ptr("app.js"),
			ImportMap: &v1alpha1.ImportMap{
				Imports: map[string]string{
					"vue": "./vue.js",
				},
			},
		},
		Status: v1alpha1.MicroFrontendStatus{
			FrontendClassRef: &v1alpha1.MicroFrontendClassReference{
				Name:     className,
				Accepted: true,
			},
			ImportMapConflicts: []v1alpha1.ImportMapConflict{
				{
					Specifier:      "vue",
					RequestedPath:  "./vue.js",
					RegisteredPath: "./vue-other.js",
					RegisteredBy:   "default/other-mf",
				},
			},
		},
	}

	_ = microFrontendRepository.Store(mf1)
	_ = microFrontendRepository.Store(mf2)
	spa := NewSinglePageApplication(&logger, microFrontendRepository)

	mfc := &v1alpha1.MicroFrontendClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      className,
			Namespace: "default",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			BaseUri: ptr("/"),
			Title:   ptr("Test"),
		},
	}

	req, _ := http.NewRequest("GET", "/", nil)

	// Act
	result, err := spa.buildImportMap(req, mfc, logger)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should include mf1
	if !strings.Contains(result, `"react"`) {
		t.Fatalf("expected import map to contain react from mf1, got %s", result)
	}

	// Should NOT include mf2 (has conflicts)
	if strings.Contains(result, `"vue"`) {
		t.Fatalf("expected import map to NOT contain vue from conflicted mf2, got %s", result)
	}
}

func TestBuildImportMapExcludesRejectedMicrofrontends(t *testing.T) {
	// Arrange
	logger := logr.Logger{}
	microFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()

	className := "test-class"

	// Rejected microfrontend
	mf := &v1alpha1.MicroFrontend{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "mf1",
			Namespace:         "default",
			CreationTimestamp: metav1.Now(),
		},
		Spec: v1alpha1.MicroFrontendSpec{
			FrontendClass: &className,
			Service:       &v1alpha1.ServiceReference{URI: ptr("https://example.com")},
			ModulePath:    ptr("app.js"),
			ImportMap: &v1alpha1.ImportMap{
				Imports: map[string]string{
					"react": "./react.js",
				},
			},
		},
		Status: v1alpha1.MicroFrontendStatus{
			FrontendClassRef: &v1alpha1.MicroFrontendClassReference{
				Name:     className,
				Accepted: false, // Not accepted
			},
		},
	}

	_ = microFrontendRepository.Store(mf)
	spa := NewSinglePageApplication(&logger, microFrontendRepository)

	mfc := &v1alpha1.MicroFrontendClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      className,
			Namespace: "default",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			BaseUri: ptr("/"),
			Title:   ptr("Test"),
		},
	}

	req, _ := http.NewRequest("GET", "/", nil)

	// Act
	result, err := spa.buildImportMap(req, mfc, logger)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should be empty since the only MF is not accepted
	if result != "{}" {
		t.Fatalf("expected empty import map for rejected MF, got %s", result)
	}
}
