package middleware

import (
	"context"
	"net/http"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
)

// WithTracing adds OpenTelemetry tracing to HTTP requests
func WithTracing(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract trace context from request headers
			ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			// Start a new service span
			ctx, span := tracing.StartServiceSpan(ctx, serviceName, r.URL.Path)
			defer span.End()

			// Add HTTP-specific attributes
			tracing.AddSpanAttributes(ctx,
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.host", r.Host),
				attribute.String("http.user_agent", r.UserAgent()),
			)

			// Add trace context to response headers
			otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(w.Header()))

			// Call next handler with the new context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GRPCTracingInterceptor adds OpenTelemetry tracing to gRPC requests
func GRPCTracingInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Start a new service span
		ctx, span := tracing.StartServiceSpan(ctx, serviceName, info.FullMethod)
		defer span.End()

		// Add gRPC-specific attributes
		tracing.AddSpanAttributes(ctx,
			attribute.String("rpc.system", "grpc"),
			attribute.String("rpc.method", info.FullMethod),
		)

		// Call the handler with the new context
		return handler(ctx, req)
	}
}
