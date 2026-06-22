package middleware

import (
	"net/http"

	"github.com/PaingPhyoAungKhant/url-shortener/internal/apperr"
)

// Error is a middleware that injects a mutable AppErr slot into ctx for downstream handlers to set.
func Error(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := apperr.WithAppErr(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
