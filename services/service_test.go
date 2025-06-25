package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lyxuansang91/shoping-cart-interview/configs"
	"github.com/stretchr/testify/assert"
)

func TestNewURLShortenerService(t *testing.T) {
	config := &configs.Config{
		Port:      "8080",
		BaseURL:   "http://localhost:8080",
		LogLevel:  "info",
		EnableCORS: true,
	}
	service := NewURLShortenerService(config)
	assert.NotNil(t, service)
	assert.NotNil(t, service.links)
	assert.Equal(t, config, service.config)
}

func TestGenerateShortCode(t *testing.T) {
	config := &configs.Config{BaseURL: "http://localhost:8080"}
	service := NewURLShortenerService(config)
	
	// Test that generated codes are 6 characters long
	code1 := service.generateShortCode()
	code2 := service.generateShortCode()
	
	assert.Len(t, code1, 6)
	assert.Len(t, code2, 6)
	
	// Test that codes are different (very unlikely to be the same)
	assert.NotEqual(t, code1, code2)
}

func TestShortenURL_Success(t *testing.T) {
	// Setup
	e := echo.New()
	config := &configs.Config{BaseURL: "http://localhost:8080"}
	service := NewURLShortenerService(config)
	
	reqBody := ShortenRequest{
		LongURL: "https://www.google.com",
	}
	jsonBody, _ := json.Marshal(reqBody)
	
	req := httptest.NewRequest(http.MethodPost, "/api/shortlinks", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Test
	err := service.ShortenURL(c)
	
	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	
	var response ShortenResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(response.ShortURL, "http://localhost:8080/shortlinks/"))
	assert.Len(t, response.ID, 6)
}

func TestShortenURL_InvalidJSON(t *testing.T) {
	// Setup
	e := echo.New()
	config := &configs.Config{BaseURL: "http://localhost:8080"}
	service := NewURLShortenerService(config)
	
	req := httptest.NewRequest(http.MethodPost, "/api/shortlinks", strings.NewReader("invalid json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Test
	err := service.ShortenURL(c)
	
	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	
	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid request body", response["error"])
}

func TestShortenURL_InvalidURL(t *testing.T) {
	// Setup
	e := echo.New()
	config := &configs.Config{BaseURL: "http://localhost:8080"}
	service := NewURLShortenerService(config)
	
	reqBody := ShortenRequest{
		LongURL: "not-a-valid-url",
	}
	jsonBody, _ := json.Marshal(reqBody)
	
	req := httptest.NewRequest(http.MethodPost, "/api/shortlinks", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Test
	err := service.ShortenURL(c)
	
	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	
	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid URL format", response["error"])
}

func TestShortenURL_EmptyURL(t *testing.T) {
	// Setup
	e := echo.New()
	config := &configs.Config{BaseURL: "http://localhost:8080"}
	service := NewURLShortenerService(config)
	
	reqBody := ShortenRequest{
		LongURL: "",
	}
	jsonBody, _ := json.Marshal(reqBody)
	
	req := httptest.NewRequest(http.MethodPost, "/api/shortlinks", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	// Test
	err := service.ShortenURL(c)
	
	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	
	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid URL format", response["error"])
}

func TestGetShortLinkDetail_Success(t *testing.T) {
	e := echo.New()
	config := &configs.Config{BaseURL: "http://localhost:8080"}
	service := NewURLShortenerService(config)

	// Create a short link
	id := "abc123"
	createdAt := time.Now().UTC()
	service.links[id] = &ShortLink{
		ID:          id,
		OriginalURL: "https://example.com",
		CreatedAt:   createdAt,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/shortlinks/abc123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc123")

	err := service.GetShortLinkDetail(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response ShortLinkDetailResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, id, response.ID)
	assert.Equal(t, "https://example.com", response.OriginalURL)
	assert.WithinDuration(t, createdAt, response.CreatedAt, time.Second)
}

func TestRedirectToLongURL_Success(t *testing.T) {
	e := echo.New()
	config := &configs.Config{BaseURL: "http://localhost:8080"}
	service := NewURLShortenerService(config)

	id := "abc123"
	service.links[id] = &ShortLink{
		ID:          id,
		OriginalURL: "https://example.com",
		CreatedAt:   time.Now().UTC(),
	}

	req := httptest.NewRequest(http.MethodGet, "/shortlinks/abc123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc123")

	err := service.RedirectToLongURL(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "https://example.com", rec.Header().Get("Location"))
}

func TestGetShortLinkDetail_NotFound(t *testing.T) {
	e := echo.New()
	config := &configs.Config{BaseURL: "http://localhost:8080"}
	service := NewURLShortenerService(config)

	req := httptest.NewRequest(http.MethodGet, "/api/shortlinks/nonexistent", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("nonexistent")

	err := service.GetShortLinkDetail(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Short link not found", response["error"])
}

func TestRedirectToLongURL_NotFound(t *testing.T) {
	e := echo.New()
	config := &configs.Config{BaseURL: "http://localhost:8080"}
	service := NewURLShortenerService(config)

	req := httptest.NewRequest(http.MethodGet, "/shortlinks/nonexistent", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("nonexistent")

	err := service.RedirectToLongURL(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Short link not found", response["error"])
}

func TestURLShortenerService_Concurrency(t *testing.T) {
	// Setup
	config := &configs.Config{BaseURL: "http://localhost:8080"}
	service := NewURLShortenerService(config)
	
	// Test concurrent access
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			// Simulate concurrent shortening
			service.mutex.Lock()
			shortCode := service.generateShortCode()
			service.links[shortCode] = &ShortLink{
				ID:          shortCode,
				OriginalURL: "https://example.com",
				CreatedAt:   time.Now().UTC(),
			}
			service.mutex.Unlock()
			
			// Simulate concurrent reading
			service.mutex.RLock()
			_, exists := service.links[shortCode]
			service.mutex.RUnlock()
			
			assert.True(t, exists)
			done <- true
		}()
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Verify we have 10 entries
	assert.Len(t, service.links, 10)
}

func TestShortenURL_UniqueCodes(t *testing.T) {
	// Setup
	e := echo.New()
	config := &configs.Config{BaseURL: "http://localhost:8080"}
	service := NewURLShortenerService(config)
	
	// Create multiple short URLs
	urls := []string{
		"https://www.google.com",
		"https://www.github.com",
		"https://www.stackoverflow.com",
	}
	
	shortCodes := make(map[string]bool)
	
	for _, longURL := range urls {
		reqBody := ShortenRequest{LongURL: longURL}
		jsonBody, _ := json.Marshal(reqBody)
		
		req := httptest.NewRequest(http.MethodPost, "/api/shortlinks", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		err := service.ShortenURL(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		
		var response ShortenResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		shortCode := strings.TrimPrefix(response.ShortURL, "http://localhost:8080/shortlinks/")
		assert.False(t, shortCodes[shortCode], "Short code should be unique: %s", shortCode)
		shortCodes[shortCode] = true
	}
	
	assert.Len(t, shortCodes, 3)
}

func TestGetConfig(t *testing.T) {
	// Test default values
	config := configs.GetConfig()
	assert.Equal(t, "8080", config.Port)
	assert.Equal(t, "http://localhost:8080", config.BaseURL)
	assert.Equal(t, "info", config.LogLevel)
	assert.True(t, config.EnableCORS)

	// Test environment variable override
	os.Setenv("PORT", "9090")
	os.Setenv("BASE_URL", "https://example.com")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("ENABLE_CORS", "false")

	config = configs.GetConfig()
	assert.Equal(t, "9090", config.Port)
	assert.Equal(t, "https://example.com", config.BaseURL)
	assert.Equal(t, "debug", config.LogLevel)
	assert.False(t, config.EnableCORS)

	// Clean up
	os.Unsetenv("PORT")
	os.Unsetenv("BASE_URL")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("ENABLE_CORS")
}

func TestShortenURL_DuplicateURL(t *testing.T) {
	// Setup
	e := echo.New()
	config := &configs.Config{BaseURL: "http://localhost:8080"}
	service := NewURLShortenerService(config)
	
	// First request
	reqBody := ShortenRequest{
		LongURL: "https://www.google.com",
	}
	jsonBody, _ := json.Marshal(reqBody)
	
	req1 := httptest.NewRequest(http.MethodPost, "/api/shortlinks", bytes.NewReader(jsonBody))
	req1.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec1 := httptest.NewRecorder()
	c1 := e.NewContext(req1, rec1)
	
	err := service.ShortenURL(c1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec1.Code)
	
	var response1 ShortenResponse
	err = json.Unmarshal(rec1.Body.Bytes(), &response1)
	assert.NoError(t, err)
	firstID := response1.ID
	
	// Second request with same URL
	req2 := httptest.NewRequest(http.MethodPost, "/api/shortlinks", bytes.NewReader(jsonBody))
	req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)
	
	err = service.ShortenURL(c2)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec2.Code) // Should return 200 OK for existing URL
	
	var response2 ShortenResponse
	err = json.Unmarshal(rec2.Body.Bytes(), &response2)
	assert.NoError(t, err)
	
	// Should return the same ID
	assert.Equal(t, firstID, response2.ID)
	assert.Equal(t, response1.ShortURL, response2.ShortURL)
	
	// Verify only one entry exists in the service
	assert.Len(t, service.links, 1)
} 