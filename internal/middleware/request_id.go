package middleware

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/PaingPhyoAungKhant/url-shortener/internal/logger"
)

const headerRequestID = "X-Request-ID"

// RequestID get/set requestID in header and set it as traceID in context
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(headerRequestID)
		if id == "" {
			id = uuid.NewString()
		}
		ctx := logger.WithTraceID(r.Context(), id)
		r.Header.Set(headerRequestID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
