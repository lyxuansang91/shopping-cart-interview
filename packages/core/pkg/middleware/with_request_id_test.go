package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cinchprotocol/cinch-api/packages/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock logger for testing with_request_id middleware
type mockLoggerWithRequestID struct {
	mock.Mock
}

func (m *mockLoggerWithRequestID) Debug(ctx context.Context, msg string, fields ...core.Field) {
	m.Called(ctx, msg, fields)
}

func (m *mockLoggerWithRequestID) Info(ctx context.Context, msg string, fields ...core.Field) {
	m.Called(ctx, msg, fields)
}

func (m *mockLoggerWithRequestID) Warn(ctx context.Context, msg string, fields ...core.Field) {
	m.Called(ctx, msg, fields)
}

func (m *mockLoggerWithRequestID) Error(ctx context.Context, msg string, fields ...core.Field) {
	m.Called(ctx, msg, fields)
}

func (m *mockLoggerWithRequestID) Fatal(ctx context.Context, msg string, fields ...core.Field) {
	m.Called(ctx, msg, fields)
}

func (m *mockLoggerWithRequestID) Debugf(ctx context.Context, format string, args ...interface{}) {
	m.Called(ctx, format, args)
}

func (m *mockLoggerWithRequestID) Infof(ctx context.Context, format string, args ...interface{}) {
	m.Called(ctx, format, args)
}

func (m *mockLoggerWithRequestID) Warnf(ctx context.Context, format string, args ...interface{}) {
	m.Called(ctx, format, args)
}

func (m *mockLoggerWithRequestID) Errorf(ctx context.Context, format string, args ...interface{}) {
	m.Called(ctx, format, args)
}

func (m *mockLoggerWithRequestID) Fatalf(ctx context.Context, format string, args ...interface{}) {
	m.Called(ctx, format, args)
}

func (m *mockLoggerWithRequestID) With(fields ...core.Field) core.Logger {
	args := m.Called(fields)
	return args.Get(0).(core.Logger)
}

func TestWithRequestID(t *testing.T) {
	t.Run("adds request ID to context and logs request", func(t *testing.T) {
		// Create mock logger
		mockLog := &mockLoggerWithRequestID{}

		// Set up expectations - the Info method should be called with specific fields
		mockLog.On("Info", mock.MatchedBy(func(ctx context.Context) bool {
			// Verify that trace ID was added to context
			traceID := core.GetTraceID(ctx)
			return traceID != ""
		}), "Request started", mock.MatchedBy(func(fields []core.Field) bool {
			// Verify that method and path fields are present
			if len(fields) != 2 {
				return false
			}

			hasMethod := false
			hasPath := false

			for _, field := range fields {
				if field.Key == "method" && field.Value == "GET" {
					hasMethod = true
				}
				if field.Key == "path" && field.Value == "/test" {
					hasPath = true
				}
			}

			return hasMethod && hasPath
		})).Once()

		// Create test handler that verifies context has trace ID
		var contextTraceID string
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contextTraceID = core.GetTraceID(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		// Create middleware
		middleware := WithRequestID(mockLog)(testHandler)

		// Create test request
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		// Execute middleware
		middleware.ServeHTTP(rec, req)

		// Verify response
		assert.Equal(t, http.StatusOK, rec.Code)

		// Verify that trace ID was added to context
		assert.NotEmpty(t, contextTraceID)

		// Verify logger was called
		mockLog.AssertExpectations(t)
	})

	t.Run("handles POST request correctly", func(t *testing.T) {
		// Create mock logger
		mockLog := &mockLoggerWithRequestID{}

		// Set up expectations for POST request
		mockLog.On("Info", mock.AnythingOfType("*context.valueCtx"), "Request started", mock.MatchedBy(func(fields []core.Field) bool {
			// Verify that method and path fields are present with correct values
			if len(fields) != 2 {
				return false
			}

			hasMethod := false
			hasPath := false

			for _, field := range fields {
				if field.Key == "method" && field.Value == "POST" {
					hasMethod = true
				}
				if field.Key == "path" && field.Value == "/api/create" {
					hasPath = true
				}
			}

			return hasMethod && hasPath
		})).Once()

		// Create test handler
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		})

		// Create middleware
		middleware := WithRequestID(mockLog)(testHandler)

		// Create test request
		req := httptest.NewRequest(http.MethodPost, "/api/create", nil)
		rec := httptest.NewRecorder()

		// Execute middleware
		middleware.ServeHTTP(rec, req)

		// Verify response
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Verify logger was called
		mockLog.AssertExpectations(t)
	})

	t.Run("generates unique request IDs", func(t *testing.T) {
		// Create mock logger
		mockLog := &mockLoggerWithRequestID{}

		// Allow any calls to Info method
		mockLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()

		var traceID1, traceID2 string

		// Create test handler that captures trace ID
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if traceID1 == "" {
				traceID1 = core.GetTraceID(r.Context())
			} else {
				traceID2 = core.GetTraceID(r.Context())
			}
			w.WriteHeader(http.StatusOK)
		})

		// Create middleware
		middleware := WithRequestID(mockLog)(testHandler)

		// Execute first request
		req1 := httptest.NewRequest(http.MethodGet, "/test1", nil)
		rec1 := httptest.NewRecorder()
		middleware.ServeHTTP(rec1, req1)

		// Execute second request
		req2 := httptest.NewRequest(http.MethodGet, "/test2", nil)
		rec2 := httptest.NewRecorder()
		middleware.ServeHTTP(rec2, req2)

		// Verify unique trace IDs were generated
		assert.NotEmpty(t, traceID1)
		assert.NotEmpty(t, traceID2)
		assert.NotEqual(t, traceID1, traceID2)
	})
}
