package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lyxuansang91/shoping-cart-interview/configs"
)

// ShortLink represents a shortened URL with metadata
type ShortLink struct {
	ID         string    `json:"id"`
	OriginalURL string   `json:"original_url"`
	CreatedAt  time.Time `json:"created_at"`
}

// URLShortenerService handles URL shortening operations
type URLShortenerService struct {
	links   map[string]*ShortLink // id -> ShortLink
	urlToID map[string]string     // originalURL -> id (for duplicate detection)
	mutex   sync.RWMutex
	config  *configs.Config
}

// ShortenRequest represents the request body for shortening a URL
type ShortenRequest struct {
	LongURL string `json:"long_url" validate:"required,url"`
}

// ShortenResponse represents the response for a shortened URL
type ShortenResponse struct {
	ShortURL string `json:"short_url"`
	ID       string `json:"id"`
}

// ShortLinkDetailResponse for GET /api/shortlinks/{id}
type ShortLinkDetailResponse struct {
	ID         string    `json:"id"`
	OriginalURL string   `json:"original_url"`
	CreatedAt  time.Time `json:"created_at"`
}

// NewURLShortenerService creates a new URL shortener service instance
func NewURLShortenerService(config *configs.Config) *URLShortenerService {
	return &URLShortenerService{
		links:   make(map[string]*ShortLink),
		urlToID: make(map[string]string),
		config:  config,
	}
}

// generateShortCode generates a random 6-character short code
func (s *URLShortenerService) generateShortCode() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:6]
}

// ShortenURL handles the POST /api/shortlinks request
func (s *URLShortenerService) ShortenURL(c echo.Context) error {
	var req ShortenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate URL
	if _, err := url.ParseRequestURI(req.LongURL); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid URL format",
		})
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if URL already exists
	if existingID, exists := s.urlToID[req.LongURL]; exists {
		shortURL := fmt.Sprintf("%s/shortlinks/%s", s.config.BaseURL, existingID)
		return c.JSON(http.StatusOK, ShortenResponse{
			ShortURL: shortURL,
			ID:       existingID,
		})
	}

	// Generate new short code
	id := s.generateShortCode()
	for _, exists := s.links[id]; exists; {
		id = s.generateShortCode()
	}

	shortLink := &ShortLink{
		ID:         id,
		OriginalURL: req.LongURL,
		CreatedAt:  time.Now().UTC(),
	}
	s.links[id] = shortLink
	s.urlToID[req.LongURL] = id

	shortURL := fmt.Sprintf("%s/shortlinks/%s", s.config.BaseURL, id)

	return c.JSON(http.StatusCreated, ShortenResponse{
		ShortURL: shortURL,
		ID:       id,
	})
}

// GetShortLinkDetail handles GET /api/shortlinks/{id}
func (s *URLShortenerService) GetShortLinkDetail(c echo.Context) error {
	id := c.Param("id")

	s.mutex.RLock()
	link, exists := s.links[id]
	s.mutex.RUnlock()

	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Short link not found",
		})
	}

	return c.JSON(http.StatusOK, ShortLinkDetailResponse{
		ID:         link.ID,
		OriginalURL: link.OriginalURL,
		CreatedAt:  link.CreatedAt,
	})
}

// RedirectToLongURL handles GET /shortlinks/{id} for public redirect
func (s *URLShortenerService) RedirectToLongURL(c echo.Context) error {
	id := c.Param("id")

	s.mutex.RLock()
	link, exists := s.links[id]
	s.mutex.RUnlock()

	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Short link not found",
		})
	}

	return c.Redirect(http.StatusFound, link.OriginalURL) // 302
} 