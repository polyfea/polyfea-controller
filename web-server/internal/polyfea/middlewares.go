package polyfea

import (
	"context"
	"net/http"
	"strings"
)

type PolyfeaContextKey string

const PolyfeaContextKeyUserRoles PolyfeaContextKey = "userRoles"
const PolyfeaContextKeyBasePath PolyfeaContextKey = "basePath"

func BasePathStrippingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		parts := strings.SplitN(r.URL.Path, "/polyfea", 2)
		if len(parts) > 1 {
			r.URL.Path = "/polyfea" + parts[1]
			basePath := parts[0] + "/"
			if basePath[0] != '/' {
				basePath = "/" + basePath
			}
			r = r.WithContext(context.WithValue(r.Context(), PolyfeaContextKeyBasePath, basePath))
		}

		next.ServeHTTP(w, r)
	})
}
