package polyfea

import (
	"testing"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/rs/zerolog"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestServeAppWebManifest(t *testing.T) {
	// Arrange
	pwa := NewProgressiveWebApplication(&zerolog.Logger{})

	mfc := &v1alpha1.MicroFrontendClass{
		ObjectMeta: v1.ObjectMeta{
			Name:      "polyfea",
			Namespace: "polyfea",
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			Title:     &[]string{"Test MicroFrontendClass"}[0],
			BaseUri:   &[]string{"/someother"}[0],
			CspHeader: "default-src 'self';",
			ExtraHeaders: []v1alpha1.Header{
				{
					Name:  "X-Frame-Options",
					Value: "DENY",
				},
			},
			UserRolesHeader: "X-User-Roles",
			UserHeader:      "X-User-Id",
			ProgressiveWebApp: &v1alpha1.ProgressiveWebApp{
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
			},
		},
	}

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
	actual := pwa.serveAppWebManifest(mfc)

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
