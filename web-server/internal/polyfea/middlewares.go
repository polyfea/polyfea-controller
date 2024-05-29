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

func BasePathStrippingMiddleware(next http.Handler, microFrontendClassRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontendClass]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := telemetry().tracer.Start(r.Context(), "spa_d.basepath_stripping_middleware",
			trace.WithAttributes(
				attribute.String("path", r.URL.Path),
				attribute.String("method", r.Method),
			),
		)
		defer span.End()

		originalPath := r.URL.Path
		var basePath string

		// Check if the route contains /polyfea
		if polyfeaIndex := strings.Index(originalPath, "/polyfea"); polyfeaIndex != -1 {
			// Consider everything before /polyfea as basePath
			basePath = originalPath[:polyfeaIndex]
			// Adjust r.URL.Path to include everything after /polyfea
			r.URL.Path = "/polyfea" + originalPath[polyfeaIndex+len("/polyfea"):]
		} else {
			basePath = originalPath
		}

		// Retrieve the micro frontend class based on the adjusted path
		basePath, microFrontendClass, err := getMicrofrontendClassAndBase(basePath, microFrontendClassRepository)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update the context with the basePath and microFrontendClass
		ctx = context.WithValue(ctx, PolyfeaContextKeyBasePath, basePath)
		ctx = context.WithValue(ctx, PolyfeaContextKeyMicroFrontendClass, microFrontendClass)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func getMicrofrontendClassAndBase(requestPath string, microFrontendClassRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontendClass]) (string, *v1alpha1.MicroFrontendClass, error) {
	slash := func(in string) string {
		if len(in) == 0 || in[len(in)-1] != '/' {
			in += "/"
		}
		if in[0] != '/' {
			in = "/" + in
		}
		return in
	}

	requestPath = slash(requestPath) // needed for user's forgotten trailing slash

	microFrontendClasses, err := microFrontendClassRepository.GetItems(func(mfc *v1alpha1.MicroFrontendClass) bool {
		return strings.HasPrefix(requestPath, slash(*mfc.Spec.BaseUri))
	})

	if err != nil {
		return "", nil, err
	}

	if len(microFrontendClasses) == 0 {
		return "/", nil, nil
	}

	baseHref := "/"
	longestMfc := microFrontendClasses[0]
	// find longest match
	for _, mfc := range microFrontendClasses {
		mfcBase := slash(*mfc.Spec.BaseUri)
		if len(mfcBase) > len(baseHref) {
			baseHref = mfcBase
			longestMfc = mfc
		}
	}

	return baseHref, longestMfc, nil
}
