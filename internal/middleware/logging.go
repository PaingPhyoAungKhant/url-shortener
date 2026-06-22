package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/PaingPhyoAungKhant/url-shortener/internal/apperr"
	"github.com/PaingPhyoAungKhant/url-shortener/internal/logger"
)

// StatusRecorder wraps the http.ResponseWriter to capture
// the HTTP Status Code written by downstream handler
type StatusRecorder struct {
	http.ResponseWriter
	status int
}

// WriteHeader wraps the statusCode before forwarding it.
func (sr *StatusRecorder) WriteHeader(statusCode int) {
	sr.status = statusCode
	sr.ResponseWriter.WriteHeader(statusCode)
}

// Logger returns middleware that one structured log line per request after downstream handler returns
func Logging(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &StatusRecorder{w, http.StatusOK}
			next.ServeHTTP(rec, r)
			ctx := r.Context()
			if appErrPtr, ok := ctx.Value(apperr.AppErrKey).(**apperr.AppErr); ok && appErrPtr != nil {
				appErr := *appErrPtr
				if rec.status >= 500 {
					fields := []slog.Attr{
						logger.String("path", r.URL.Path),
						logger.Int("status", rec.status),
						logger.String("latency", time.Since(start).String()),
					}
					if appErr != nil {
						fields = append(fields, logger.Err(appErr.LogError()))
					}
					log.Error(ctx, "internal error", fields...)
				} else {
					log.Info(
						ctx,
						"request",
						logger.String("path", r.URL.Path),
						logger.Int("status", rec.status),
						logger.String("latency", time.Since(start).String()),
					)
				}
			}
		})
	}
}
