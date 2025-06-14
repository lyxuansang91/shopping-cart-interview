package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewField(t *testing.T) {
	t.Run("creates field with string value", func(t *testing.T) {
		field := NewField("key", "value")
		assert.Equal(t, "key", field.Key)
		assert.Equal(t, "value", field.Value)
	})

	t.Run("creates field with int value", func(t *testing.T) {
		field := NewField("count", 42)
		assert.Equal(t, "count", field.Key)
		assert.Equal(t, 42, field.Value)
	})

	t.Run("creates field with nil value", func(t *testing.T) {
		field := NewField("null_field", nil)
		assert.Equal(t, "null_field", field.Key)
		assert.Nil(t, field.Value)
	})

	t.Run("creates field with complex object", func(t *testing.T) {
		obj := map[string]interface{}{"nested": "value"}
		field := NewField("object", obj)
		assert.Equal(t, "object", field.Key)
		assert.Equal(t, obj, field.Value)
	})
}

func TestGetTraceID(t *testing.T) {
	t.Run("returns empty string for nil context", func(t *testing.T) {
		traceID := GetTraceID(nil)
		assert.Equal(t, "", traceID)
	})

	t.Run("returns empty string for context without trace ID", func(t *testing.T) {
		ctx := context.Background()
		traceID := GetTraceID(ctx)
		assert.Equal(t, "", traceID)
	})

	t.Run("returns trace ID from context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TraceIDKey, "test-trace-id")
		traceID := GetTraceID(ctx)
		assert.Equal(t, "test-trace-id", traceID)
	})

	t.Run("returns empty string for wrong type in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TraceIDKey, 123)
		traceID := GetTraceID(ctx)
		assert.Equal(t, "", traceID)
	})
}

func TestWithTraceID(t *testing.T) {
	t.Run("adds trace ID to nil context", func(t *testing.T) {
		ctx := WithTraceID(nil, "new-trace-id")
		traceID := GetTraceID(ctx)
		assert.Equal(t, "new-trace-id", traceID)
	})

	t.Run("adds trace ID to background context", func(t *testing.T) {
		ctx := WithTraceID(context.Background(), "new-trace-id")
		traceID := GetTraceID(ctx)
		assert.Equal(t, "new-trace-id", traceID)
	})

	t.Run("appends trace ID to existing trace ID", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TraceIDKey, "existing-id")
		ctx = WithTraceID(ctx, "new-id")
		traceID := GetTraceID(ctx)
		assert.Equal(t, "existing-id-new-id", traceID)
	})

	t.Run("chains multiple trace IDs", func(t *testing.T) {
		ctx := context.Background()
		ctx = WithTraceID(ctx, "first")
		ctx = WithTraceID(ctx, "second")
		ctx = WithTraceID(ctx, "third")
		traceID := GetTraceID(ctx)
		assert.Equal(t, "first-second-third", traceID)
	})

	t.Run("handles empty trace ID", func(t *testing.T) {
		ctx := WithTraceID(context.Background(), "")
		traceID := GetTraceID(ctx)
		assert.Equal(t, "", traceID)
	})

	t.Run("appends to existing trace ID even with empty new ID", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TraceIDKey, "existing")
		ctx = WithTraceID(ctx, "")
		traceID := GetTraceID(ctx)
		assert.Equal(t, "existing-", traceID)
	})
}
