package tracing

import (
	"context"
	"fmt"

	"github.com/cinchprotocol/cinch-api/packages/core"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// GetSpanID returns the current span ID from the context
func GetSpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().HasSpanID() {
		return ""
	}
	return span.SpanContext().SpanID().String()
}

// GetTraceID returns the current trace ID from the context
func GetTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().HasTraceID() {
		return ""
	}
	return span.SpanContext().TraceID().String()
}

// StartServiceSpan starts a new span for a service entry point
func StartServiceSpan(ctx context.Context, serviceName, operation string) (context.Context, trace.Span) {
	tracer := otel.Tracer(serviceName)
	ctx, span := tracer.Start(ctx, fmt.Sprintf("%s.%s", serviceName, operation),
		trace.WithAttributes(
			attribute.String("service", serviceName),
			attribute.String("operation", operation),
		),
	)
	return ctx, span
}

// StartFunctionSpan starts a new span for a function call
func StartFunctionSpan(ctx context.Context, serviceName, functionName string) (context.Context, trace.Span) {
	tracer := otel.Tracer(serviceName)
	ctx, span := tracer.Start(ctx, functionName,
		trace.WithAttributes(
			attribute.String("service", serviceName),
			attribute.String("function", functionName),
		),
	)
	return ctx, span
}

// AddSpanAttributes adds attributes to the current span
func AddSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.SetAttributes(attrs...)
	}
}

// WithSpanContext adds span context to a logger
func WithSpanContext(ctx context.Context, logger core.Logger) core.Logger {
	spanID := GetSpanID(ctx)
	traceID := GetTraceID(ctx)

	fields := []core.Field{}
	if spanID != "" {
		fields = append(fields, core.NewField("span_id", spanID))
	}
	if traceID != "" {
		fields = append(fields, core.NewField("trace_id", traceID))
	}

	return logger.With(fields...)
}
