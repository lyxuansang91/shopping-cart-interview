package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock cache implementation
type mockCache struct {
	mock.Mock
}

func (m *mockCache) CloseConnection() {
	m.Called()
}

func (m *mockCache) Exists(key string) (bool, error) {
	args := m.Called(key)
	return args.Bool(0), args.Error(1)
}

func (m *mockCache) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *mockCache) Set(key string, value string) (bool, error) {
	args := m.Called(key, value)
	return args.Bool(0), args.Error(1)
}

func TestRequestEventId(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		urlPath        string
		eventID        string
		mockSetup      func(*mockCache)
		expectedStatus int
	}{
		{
			name:    "Valid Event ID",
			urlPath: "/api/v1/resource1",
			eventID: "valid-event-id",
			mockSetup: func(mc *mockCache) {
				mc.On("Exists", "valid-event-id").Return(false, nil)
				mc.On("Set", "valid-event-id", "1").Return(true, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing Event ID Header",
			urlPath:        "/api/v1/resource1",
			eventID:        "",
			mockSetup:      func(mc *mockCache) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "Event ID Already Processed",
			urlPath: "/api/v1/resource1",
			eventID: "processed-event-id",
			mockSetup: func(mc *mockCache) {
				mc.On("Exists", "processed-event-id").Return(true, nil)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Event ID not checked for non-resource1 path",
			urlPath:        "/api/v1/not-resource1",
			eventID:        "",
			mockSetup:      func(mc *mockCache) {},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "Cache error during Exists check",
			urlPath: "/api/v1/resource1",
			eventID: "error-event-id",
			mockSetup: func(mc *mockCache) {
				mc.On("Exists", "error-event-id").Return(false, assert.AnError)
			},
			expectedStatus: http.StatusOK, // middleware continues on error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock cache
			mockCacheClient := &mockCache{}
			tt.mockSetup(mockCacheClient)

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			middleware := RequestEventId(mockCacheClient)(testHandler)

			req := httptest.NewRequest(http.MethodGet, tt.urlPath, nil)
			if tt.eventID != "" {
				req.Header.Set("X-Event-ID", tt.eventID)
			}

			rec := httptest.NewRecorder()

			middleware.ServeHTTP(rec, req)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockCacheClient.AssertExpectations(t)
		})
	}
}
