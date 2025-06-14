package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core"
	"github.com/cinchprotocol/cinch-api/packages/core/pkg/middleware"
	"github.com/go-chi/chi/v5"
)

// RegisterFunction is a type for the router registration function
type RegisterFunction func(*chi.Mux)

// ServerConfig contains the configuration for the HTTP server
type ServerConfig struct {
	Port            string
	Logger          core.Logger
	ShutdownTimeout time.Duration
	ServiceName     string
}

// StartServer initializes and starts an HTTP server with the given configuration
func StartServer(ctx context.Context, config ServerConfig, register RegisterFunction) error {
	r := chi.NewRouter()

	// Add tracing middleware with service name
	r.Use(middleware.WithTracing(config.ServiceName))

	// Register routes using the provided function
	register(r)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: r,
	}

	// Handle context cancellation
	go func() {
		<-ctx.Done()
		config.Logger.Info(ctx, "Shutting down HTTP server gracefully...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			config.Logger.Error(ctx, "HTTP server shutdown error",
				core.NewField("error", err),
			)
		}
	}()

	config.Logger.Info(ctx, "HTTP Server starting",
		core.NewField("addr", httpServer.Addr),
	)

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("failed to serve HTTP: %v", err)
	}

	return nil
}
