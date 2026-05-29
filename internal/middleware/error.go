package middleware

import (
	"context"
	"net/http"

	"github.com/PaingPhyoAungKhant/url-shortener/internal/apperr"
)

// Error is an middleware that inject an AppErr pointer-to-pointer value in ctx whith AppErrKey for downsteam handlers to mutate and log error
func Error(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appErrPtr := new(*apperr.AppErr)
		ctx := context.WithValue(r.Context(), apperr.AppErrKey, appErrPtr)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
