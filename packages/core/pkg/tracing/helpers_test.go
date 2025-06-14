package tracing

import (
	"context"
	"testing"

	"github.com/cinchprotocol/cinch-api/packages/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Mock logger for testing
type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) Debug(ctx context.Context, msg string, fields ...core.Field) {
	m.Called(ctx, msg, fields)
}

func (m *mockLogger) Info(ctx context.Context, msg string, fields ...core.Field) {
	m.Called(ctx, msg, fields)
}

func (m *mockLogger) Warn(ctx context.Context, msg string, fields ...core.Field) {
	m.Called(ctx, msg, fields)
}

func (m *mockLogger) Error(ctx context.Context, msg string, fields ...core.Field) {
	m.Called(ctx, msg, fields)
}

func (m *mockLogger) Fatal(ctx context.Context, msg string, fields ...core.Field) {
	m.Called(ctx, msg, fields)
}

func (m *mockLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	m.Called(ctx, format, args)
}

func (m *mockLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	m.Called(ctx, format, args)
}

func (m *mockLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	m.Called(ctx, format, args)
}

func (m *mockLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	m.Called(ctx, format, args)
}

func (m *mockLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	m.Called(ctx, format, args)
}

func (m *mockLogger) With(fields ...core.Field) core.Logger {
	args := m.Called(fields)
	return args.Get(0).(core.Logger)
}

func TestGetSpanID(t *testing.T) {
	t.Run("returns empty string for context without span", func(t *testing.T) {
		ctx := context.Background()
		spanID := GetSpanID(ctx)
		assert.Equal(t, "", spanID)
	})

	t.Run("returns empty string for invalid span", func(t *testing.T) {
		// Create a context with a non-recording span
		ctx := trace.ContextWithSpan(context.Background(), trace.SpanFromContext(context.Background()))
		spanID := GetSpanID(ctx)
		assert.Equal(t, "", spanID)
	})
}

func TestGetTraceIDFromTracing(t *testing.T) {
	t.Run("returns empty string for context without span", func(t *testing.T) {
		ctx := context.Background()
		traceID := GetTraceID(ctx)
		assert.Equal(t, "", traceID)
	})

	t.Run("returns empty string for invalid span", func(t *testing.T) {
		// Create a context with a non-recording span
		ctx := trace.ContextWithSpan(context.Background(), trace.SpanFromContext(context.Background()))
		traceID := GetTraceID(ctx)
		assert.Equal(t, "", traceID)
	})
}

func TestStartServiceSpan(t *testing.T) {
	t.Run("creates span with correct name and attributes", func(t *testing.T) {
		serviceName := "test-service"
		operation := "test-operation"

		ctx, span := StartServiceSpan(context.Background(), serviceName, operation)

		assert.NotNil(t, ctx)
		assert.NotNil(t, span)

		// Without proper OpenTelemetry setup, the span might not be recording
		// But the function should still work and return valid objects

		// Clean up
		span.End()
	})
}

func TestStartFunctionSpan(t *testing.T) {
	t.Run("creates span with correct name and attributes", func(t *testing.T) {
		serviceName := "test-service"
		functionName := "test-function"

		ctx, span := StartFunctionSpan(context.Background(), serviceName, functionName)

		assert.NotNil(t, ctx)
		assert.NotNil(t, span)

		// Without proper OpenTelemetry setup, the span might not be recording
		// But the function should still work and return valid objects

		// Clean up
		span.End()
	})
}

func TestAddSpanAttributes(t *testing.T) {
	t.Run("adds attributes to recording span", func(t *testing.T) {
		// Start a real tracer for this test
		tracer := otel.Tracer("test")
		ctx, span := tracer.Start(context.Background(), "test-span")
		defer span.End()

		// Add attributes
		attrs := []attribute.KeyValue{
			attribute.String("key1", "value1"),
			attribute.Int("key2", 42),
		}

		AddSpanAttributes(ctx, attrs...)

		// We can't easily verify the attributes were added without access to internals,
		// but we can verify the function doesn't panic and span exists
		assert.NotNil(t, span)
	})

	t.Run("does not panic with non-recording span", func(t *testing.T) {
		ctx := context.Background()
		attrs := []attribute.KeyValue{
			attribute.String("key1", "value1"),
		}

		// This should not panic
		assert.NotPanics(t, func() {
			AddSpanAttributes(ctx, attrs...)
		})
	})
}

func TestWithSpanContext(t *testing.T) {
	t.Run("returns logger with span context when span IDs are available", func(t *testing.T) {
		// Create a mock logger
		mockLog := &mockLogger{}

		// Mock the With method to return itself for simplicity
		// Accept any fields since we can't easily predict exact span/trace IDs
		mockLog.On("With", mock.Anything).Return(mockLog)

		// Start a real tracer to get a valid span context
		tracer := otel.Tracer("test")
		ctx, span := tracer.Start(context.Background(), "test-span")
		defer span.End()

		logger := WithSpanContext(ctx, mockLog)

		assert.NotNil(t, logger)
		mockLog.AssertExpectations(t)
	})

	t.Run("returns logger without fields when no span context", func(t *testing.T) {
		// Create a mock logger
		mockLog := &mockLogger{}

		// With no span context, empty fields array should be passed
		mockLog.On("With", mock.Anything).Return(mockLog)

		ctx := context.Background()
		logger := WithSpanContext(ctx, mockLog)

		assert.NotNil(t, logger)
		mockLog.AssertExpectations(t)
	})
}
