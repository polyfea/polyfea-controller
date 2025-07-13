package polyfea

import (
	"context"
	"net/http"
	"strings"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type PolyfeaContextKey string

const PolyfeaContextKeyBasePath PolyfeaContextKey = "basePath"
const PolyfeaContextKeyMicroFrontendClass PolyfeaContextKey = "microFrontendClass"

func BasePathStrippingMiddleware(next http.Handler, microFrontendClassRepository repository.Repository[*v1alpha1.MicroFrontendClass]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := telemetry().tracer.Start(r.Context(), "spa_d.basepath_stripping_middleware",
			trace.WithAttributes(
				attribute.String("path", r.URL.Path),
				attribute.String("method", r.Method),
			),
		)
		defer span.End()

		originalPath := r.URL.Path
		basePath := extractBasePath(originalPath)

		// Retrieve the micro frontend class based on the adjusted path
		basePath, microFrontendClass, err := getMicrofrontendClassAndBase(basePath, microFrontendClassRepository)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Adjust specific paths
		adjustSpecialPaths(r, originalPath)

		// Update the context with the basePath and microFrontendClass
		ctx = context.WithValue(ctx, PolyfeaContextKeyBasePath, basePath)
		ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, microFrontendClass)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// Helper function to extract base path
func extractBasePath(originalPath string) string {
	if polyfeaIndex := strings.Index(originalPath, "/polyfea"); polyfeaIndex != -1 {
		return originalPath[:polyfeaIndex]
	}
	return originalPath
}

// Helper function to adjust specific paths
func adjustSpecialPaths(r *http.Request, originalPath string) {
	if strings.Contains(originalPath, "/sw.mjs") {
		r.URL.Path = "/sw.mjs"
	} else if strings.Contains(originalPath, "/polyfea-caching.json") {
		r.URL.Path = "/polyfea-caching.json"
	}
}

func getMicrofrontendClassAndBase(requestPath string, microFrontendClassRepository repository.Repository[*v1alpha1.MicroFrontendClass]) (string, *v1alpha1.MicroFrontendClass, error) {
	slash := func(in string) string {
		if len(in) == 0 || in[len(in)-1] != '/' {
			in += "/"
		}
		if in[0] != '/' {
			in = "/" + in
		}
		return in
	}

	requestPath = slash(requestPath) // Ensure trailing slash

	microFrontendClasses, err := microFrontendClassRepository.List(func(mfc *v1alpha1.MicroFrontendClass) bool {
		return strings.HasPrefix(requestPath, slash(*mfc.Spec.BaseUri))
	})

	if err != nil {
		return "", nil, err
	}

	if len(microFrontendClasses) == 0 {
		return "/", nil, nil
	}

	// Find the longest matching base URI
	longestMfc := microFrontendClasses[0]
	for _, mfc := range microFrontendClasses {
		if len(slash(*mfc.Spec.BaseUri)) > len(slash(*longestMfc.Spec.BaseUri)) {
			longestMfc = mfc
		}
	}

	return slash(*longestMfc.Spec.BaseUri), longestMfc, nil
}
