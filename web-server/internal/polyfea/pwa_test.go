package polyfea

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/web-server/internal/polyfea/generated"
	"github.com/rs/zerolog"
)

var pwaTestSuite = IntegrationTestSuite{
	TestRouter: polyfeaPWAApiSetupRouter(),
	TestSet: []Test{
		{
			Name: "ServeAppWebManifestReturnsExpectedManifest",
			Func: ServeAppWebManifestReturnsExpectedManifest,
		},
		{
			Name: "ServeRegisterReturnsExpectedFile",
			Func: ServeRegisterReturnsExpectedFile,
		},
	},
}

func ServeAppWebManifestReturnsExpectedManifest(t *testing.T) {
	// Arrange
	testServerUrl := os.Getenv(TestServerUrlName)

	expected := &v1alpha1.WebAppManifest{
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
	}

	// Act
	response, err := http.Get(testServerUrl + "/polyfea/app.webmanifest")

	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	actual := &v1alpha1.WebAppManifest{}
	err = json.Unmarshal(body, &actual)
	if err != nil {
		log.Fatal(err)
	}

	// Assert
	if *actual.Name != *expected.Name {
		t.Errorf("Expected %s, got %s", *expected.Name, *actual.Name)
	}
	if *actual.Icons[0].Src != *expected.Icons[0].Src {
		t.Errorf("Expected %s, got %s", *expected.Icons[0].Src, *actual.Icons[0].Src)
	}
	if *actual.Icons[0].Sizes != *expected.Icons[0].Sizes {
		t.Errorf("Expected %s, got %s", *expected.Icons[0].Sizes, *actual.Icons[0].Sizes)
	}
	if *actual.Icons[0].Type != *expected.Icons[0].Type {
		t.Errorf("Expected %s, got %s", *expected.Icons[0].Type, *actual.Icons[0].Type)
	}
	if *actual.StartUrl != *expected.StartUrl {
		t.Errorf("Expected %s, got %s", *expected.StartUrl, *actual.StartUrl)
	}
	if *actual.Display != *expected.Display {
		t.Errorf("Expected %s, got %s", *expected.Display, *actual.Display)
	}
}

func ServeRegisterReturnsExpectedFile(t *testing.T) {
	// Arrange
	testServerUrl := os.Getenv(TestServerUrlName)
	file, err := os.ReadFile(".resources/register.mjs")
	if err != nil {
		t.Fatal(err)
	}
	expected := string(file)

	// Act
	response, err := http.Get(testServerUrl + "/polyfea/register.mjs")

	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Assert
	if string(body) != expected {
		t.Errorf("Expected %s, got %s", expected, string(body))
	}
}

func ServeServiceWorkerReturnsExpectedFile(t *testing.T) {
	// Arrange
	testServerUrl := os.Getenv(TestServerUrlName)
	file, err := os.ReadFile(".resources/sw.mjs")
	if err != nil {
		t.Fatal(err)
	}
	expected := string(file)

	// Act
	response, err := http.Get(testServerUrl + "/sw.mjs")

	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Assert
	if string(body) != expected {
		t.Errorf("Expected %s, got %s", expected, string(body))
	}
}

func polyfeaPWAApiSetupRouter() http.Handler {
	mfc := createTestMicroFrontendClass("test-frontend-class", "/")

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

	router := generated.NewRouter()

	spa := NewProgressiveWebApplication(&zerolog.Logger{})

	router.HandleFunc("/polyfea/app.webmanifest", spa.ServeAppWebManifest)
	router.HandleFunc("/polyfea/register.mjs", spa.ServeRegister)

	return addDummyMiddleware(router, "/", mfc)
}
