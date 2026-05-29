package middleware

import (
	"fmt"
	"net/http"

	"github.com/PaingPhyoAungKhant/url-shortener/internal/logger"
)

// Recovery returns middleware that catches panic from downstream handler
func Recovery(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			defer func() {
				if rec := recover(); rec != nil {
					log.Error(
						ctx,
						"panic",
						logger.Err(fmt.Errorf("%v", rec)),
					)

					http.Error(
						w,
						http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError,
					)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
