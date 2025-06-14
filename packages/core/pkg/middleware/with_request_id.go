package middleware

import (
	"net/http"

	"github.com/cinchprotocol/cinch-api/packages/core"
	"github.com/google/uuid"
)

// RequestID adds a trace ID to the context
func WithRequestID(logger core.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := uuid.New().String()
			ctx := core.WithTraceID(r.Context(), requestID)
			logger.Info(ctx, "Request started",
				core.NewField("method", r.Method),
				core.NewField("path", r.URL.Path),
			)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
