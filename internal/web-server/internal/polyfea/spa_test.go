package polyfea

import (
	"html/template"
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

const TestClassName = "test-class"

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

	expectedWithoutNonce := nonceRegex.ReplaceAllString("default-src 'self'; font-src 'self'; script-src 'strict-dynamic' 'nonce-"+nonce+"'; worker-src 'self'; manifest-src 'self'; style-src 'self' 'nonce-"+nonce+"'; style-src-attr 'self' 'unsafe-inline';", "'nonce-NONCE'")
	gotWithoutNonce := nonceRegex.ReplaceAllString(response.Header.Get("Content-Security-Policy"), "'nonce-NONCE'")

	if expectedWithoutNonce != gotWithoutNonce {
		t.Fatalf("expected content security policy %s, got %s", "default-src 'self'; font-src 'self'; script-src 'strict-dynamic' 'nonce-"+nonce+"'; worker-src 'self'; manifest-src 'self'; style-src 'self' 'nonce-"+nonce+"'; style-src-attr 'self' 'unsafe-inline';", response.Header.Get("Content-Security-Policy"))
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
	mfc.Spec.CspHeader = "default-src 'self'; font-src 'self'; script-src 'strict-dynamic' 'nonce-{NONCE_VALUE}'; worker-src 'self'; manifest-src 'self'; style-src 'self' 'nonce-{NONCE_VALUE}'; style-src-attr 'self' 'unsafe-inline';"

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
			Name:      TestClassName,
			Namespace: "default",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			BaseUri: ptr("/"),
			Title:   ptr("Test"),
		},
	}

	// Act
	result, err := spa.buildImportMap(mfc, logger)

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

	className := TestClassName
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

	// Act
	result, err := spa.buildImportMap(mfc, logger)

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
	// Paths should be proxied through the controller
	if !strings.Contains(result, `"./polyfea/proxy/default/mf1/./react.js"`) {
		t.Fatalf("expected import map to contain proxied react.js path, got %s", result)
	}
}

func TestBuildImportMapFirstRegisteredWins(t *testing.T) {
	// Arrange
	logger := logr.Logger{}
	microFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()

	className := TestClassName

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

	// Act
	result, err := spa.buildImportMap(mfc, logger)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should use mf1's version (first registered) with proxy path
	if !strings.Contains(result, `"./polyfea/proxy/default/mf1/./react-v18.js"`) {
		t.Fatalf("expected import map to contain proxied react-v18.js (first registered), got %s", result)
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

	className := TestClassName
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

	// Act
	result, err := spa.buildImportMap(mfc, logger)

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
	// Scoped imports should also be proxied
	if !strings.Contains(result, `"./polyfea/proxy/default/mf1/./react-v16.js"`) {
		t.Fatalf("expected import map to contain proxied react-v16.js in scope, got %s", result)
	}
}

func TestBuildImportMapExcludesConflictedMicrofrontends(t *testing.T) {
	// Arrange
	logger := logr.Logger{}
	microFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()

	className := TestClassName

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

	// Act
	result, err := spa.buildImportMap(mfc, logger)

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

	className := TestClassName

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

	// Act
	result, err := spa.buildImportMap(mfc, logger)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should be empty since the only MF is not accepted
	if result != "{}" {
		t.Fatalf("expected empty import map for rejected MF, got %s", result)
	}
}

func TestImportMapJSONFormatting(t *testing.T) {
	// Arrange
	logger := logr.Logger{}
	microFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()

	className := TestClassName
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
					"lodash":    "./lodash.js",
				},
				Scopes: map[string]map[string]string{
					"/legacy/": {
						"react": "./react-legacy.js",
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

	// Act
	result, err := spa.buildImportMap(mfc, logger)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify JSON is valid and not HTML-escaped
	if strings.Contains(result, "&#34;") {
		t.Fatalf("expected JSON without HTML entities, but found &#34; (escaped quote) in: %s", result)
	}
	if strings.Contains(result, "&quot;") {
		t.Fatalf("expected JSON without HTML entities, but found &quot; (escaped quote) in: %s", result)
	}
	if strings.Contains(result, "&#39;") {
		t.Fatalf("expected JSON without HTML entities, but found &#39; (escaped apostrophe) in: %s", result)
	}

	// Verify JSON contains proper double quotes
	if !strings.Contains(result, `"react"`) {
		t.Fatalf("expected JSON with proper quotes around 'react', got: %s", result)
	}
	// Paths should be proxied
	if !strings.Contains(result, `"./polyfea/proxy/default/mf1/./react.js"`) {
		t.Fatalf("expected JSON with proxied react.js path, got: %s", result)
	}

	// Verify structure
	if !strings.Contains(result, `"imports"`) {
		t.Fatalf("expected JSON to contain 'imports' key, got: %s", result)
	}
	if !strings.Contains(result, `"scopes"`) {
		t.Fatalf("expected JSON to contain 'scopes' key, got: %s", result)
	}
	if !strings.Contains(result, `"/legacy/"`) {
		t.Fatalf("expected JSON to contain '/legacy/' scope, got: %s", result)
	}

	// Verify it's valid JSON by checking structure
	if !strings.HasPrefix(result, "{") || !strings.HasSuffix(result, "}") {
		t.Fatalf("expected JSON to be wrapped in curly braces, got: %s", result)
	}
}

func TestSinglePageApplicationImportMapRendering(t *testing.T) {
	// Arrange
	logger := logr.Logger{}
	microFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()
	microFrontendClassRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontendClass]()

	className := TestClassName
	mfc := &v1alpha1.MicroFrontendClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      className,
			Namespace: "default",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			BaseUri:   ptr("/"),
			Title:     ptr("Test App"),
			CspHeader: "default-src 'self'; script-src 'nonce-{NONCE_VALUE}';",
		},
	}
	_ = microFrontendClassRepository.Store(mfc)

	mf := &v1alpha1.MicroFrontend{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "test-mf",
			Namespace:         "default",
			CreationTimestamp: metav1.Now(),
		},
		Spec: v1alpha1.MicroFrontendSpec{
			FrontendClass: &className,
			Service:       &v1alpha1.ServiceReference{URI: ptr("https://example.com")},
			ModulePath:    ptr("app.js"),
			ImportMap: &v1alpha1.ImportMap{
				Imports: map[string]string{
					"lit":      "./lit.js",
					"lit-html": "./lit-html.js",
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

	// Act
	templateVars := TemplateData{
		BaseUri:       "/",
		Title:         "Test App",
		Nonce:         "test-nonce-12345",
		ExtraMeta:     "",
		EnablePWA:     false,
		ImportMapJson: "",
	}

	// Build import map
	importMapJson, err := spa.buildImportMap(mfc, logger)
	if err != nil {
		t.Fatalf("failed to build import map: %v", err)
	}
	templateVars.ImportMapJson = template.HTML(importMapJson)

	// Template the HTML
	htmlOutput, err := templateHtml(html, &templateVars)
	if err != nil {
		t.Fatalf("failed to template HTML: %v", err)
	}

	// Assert
	// Verify no HTML entity encoding in the rendered HTML
	if strings.Contains(htmlOutput, "&#34;") {
		t.Fatalf("HTML contains &#34; (escaped quote), JSON should not be HTML-escaped in: %s", htmlOutput)
	}
	if strings.Contains(htmlOutput, "&quot;") {
		t.Fatalf("HTML contains &quot; (escaped quote), JSON should not be HTML-escaped in: %s", htmlOutput)
	}

	// Verify the import map script tag exists with proper JSON
	if !strings.Contains(htmlOutput, `<script type="importmap"`) {
		t.Fatalf("expected HTML to contain importmap script tag")
	}

	// Verify actual JSON structure with proper quotes and proxied paths
	if !strings.Contains(htmlOutput, `"lit"`) {
		t.Fatalf("expected HTML to contain properly quoted 'lit' in import map")
	}
	// Paths should be proxied through the controller
	if !strings.Contains(htmlOutput, `"./polyfea/proxy/default/test-mf/./lit.js"`) {
		t.Fatalf("expected HTML to contain proxied lit.js path in import map")
	}
	if !strings.Contains(htmlOutput, `"imports"`) {
		t.Fatalf("expected HTML to contain 'imports' key in import map")
	}

	// Verify the JSON structure is intact
	importMapStart := strings.Index(htmlOutput, `<script type="importmap"`)
	importMapEnd := strings.Index(htmlOutput[importMapStart:], `</script>`)
	if importMapStart == -1 || importMapEnd == -1 {
		t.Fatalf("could not find importmap script tag boundaries")
	}

	importMapSection := htmlOutput[importMapStart : importMapStart+importMapEnd]

	// Extract just the JSON part (after the closing >)
	jsonStart := strings.Index(importMapSection, ">") + 1
	importMapJSON := importMapSection[jsonStart:]

	// Verify it looks like valid JSON
	if !strings.HasPrefix(strings.TrimSpace(importMapJSON), "{") {
		t.Fatalf("import map JSON should start with {, got: %s", importMapJSON)
	}
}

func TestImportMapAbsoluteURLsNotProxied(t *testing.T) {
	// Arrange
	logger := logr.Logger{}
	microFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()

	className := TestClassName
	mf := &v1alpha1.MicroFrontend{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "mf-with-cdn",
			Namespace:         "default",
			CreationTimestamp: metav1.Now(),
		},
		Spec: v1alpha1.MicroFrontendSpec{
			FrontendClass: &className,
			Service:       &v1alpha1.ServiceReference{URI: ptr("https://example.com")},
			ModulePath:    ptr("app.js"),
			ImportMap: &v1alpha1.ImportMap{
				Imports: map[string]string{
					"react":     "https://cdn.example.com/react.js",     // Absolute URL - should NOT be proxied
					"react-dom": "https://cdn.example.com/react-dom.js", // Absolute URL - should NOT be proxied
					"lodash":    "./lodash.js",                          // Relative path - SHOULD be proxied
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

	// Act
	result, err := spa.buildImportMap(mfc, logger)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Absolute URLs should remain unchanged
	if !strings.Contains(result, `"https://cdn.example.com/react.js"`) {
		t.Fatalf("expected absolute URL to remain unchanged, got: %s", result)
	}
	if !strings.Contains(result, `"https://cdn.example.com/react-dom.js"`) {
		t.Fatalf("expected absolute URL to remain unchanged, got: %s", result)
	}

	// Relative path should be proxied
	if !strings.Contains(result, `"./polyfea/proxy/default/mf-with-cdn/./lodash.js"`) {
		t.Fatalf("expected relative path to be proxied, got: %s", result)
	}

	// Should NOT have proxied the absolute URLs
	if strings.Contains(result, `"./polyfea/proxy/default/mf-with-cdn/https://`) {
		t.Fatalf("absolute URLs should not be proxied, got: %s", result)
	}
}

func TestImportMapNonProxiedService(t *testing.T) {
	// Arrange
	logger := logr.Logger{}
	microFrontendRepository := repository.NewInMemoryRepository[*v1alpha1.MicroFrontend]()

	className := TestClassName
	proxy := false // Explicitly disable proxy
	mf := &v1alpha1.MicroFrontend{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "mf-no-proxy",
			Namespace:         "default",
			CreationTimestamp: metav1.Now(),
		},
		Spec: v1alpha1.MicroFrontendSpec{
			FrontendClass: &className,
			Service:       &v1alpha1.ServiceReference{URI: ptr("https://external-service.com")},
			ModulePath:    ptr("app.js"),
			Proxy:         &proxy,
			ImportMap: &v1alpha1.ImportMap{
				Imports: map[string]string{
					"react": "./react.js", // Should be resolved against service URL
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

	// Act
	result, err := spa.buildImportMap(mfc, logger)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Non-proxied service should have paths resolved against service URL
	if !strings.Contains(result, `"https://external-service.com/./react.js"`) {
		t.Fatalf("expected path to be resolved against service URL, got: %s", result)
	}

	// Should NOT be proxied
	if strings.Contains(result, `"./polyfea/proxy/`) {
		t.Fatalf("non-proxied service should not use proxy path, got: %s", result)
	}
}
