package polyfea

import (
	"context"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type PolyfeaContextKey string

const PolyfeaContextKeyBasePath PolyfeaContextKey = "basePath"

func BasePathStrippingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := telemetry().tracer.Start(r.Context(), "spa_d.basepath_stripping_middleware",
			trace.WithAttributes(
				attribute.String("path", r.URL.Path),
				attribute.String("method", r.Method),
			),
		)
		defer span.End()
		parts := strings.SplitN(r.URL.Path, "/polyfea", 2)
		if len(parts) > 1 {
			r.URL.Path = "/polyfea" + parts[1]
			basePath := parts[0] + "/"
			if basePath[0] != '/' {
				basePath = "/" + basePath
			}
			r = r.WithContext(context.WithValue(ctx, PolyfeaContextKeyBasePath, basePath))
		}

		next.ServeHTTP(w, r)
	})
}
