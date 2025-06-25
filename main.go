package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lyxuansang91/shoping-cart-interview/configs"
	"github.com/lyxuansang91/shoping-cart-interview/services"
)

func main() {
	// Load configuration from environment variables
	config := configs.GetConfig()

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	// Conditional CORS middleware
	if config.EnableCORS {
		e.Use(middleware.CORS())
	}

	// Initialize URL shortener service with config
	urlService := services.NewURLShortenerService(config)

	// API routes
	e.POST("/api/shortlinks", urlService.ShortenURL)
	e.GET("/api/shortlinks/:id", urlService.GetShortLinkDetail)

	// Public redirect endpoint
	e.GET("/shortlinks/:id", urlService.RedirectToLongURL)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})

	// Start server with configured port
	serverAddr := fmt.Sprintf(":%s", config.Port)
	log.Printf("Starting server on port %s", config.Port)
	log.Printf("Base URL: %s", config.BaseURL)
	log.Fatal(e.Start(serverAddr))
} 