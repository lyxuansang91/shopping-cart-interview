package jobs

import (
	"context"
	"time"

	"net/http"

	corehttp "github.com/cinchprotocol/cinch-api/packages/core/pkg/http"
	coremiddleware "github.com/cinchprotocol/cinch-api/packages/core/pkg/middleware"
	"github.com/cinchprotocol/cinch-api/services/cart/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// StartHTTPServer initializes and starts the HTTP server
func StartHTTPServer(ctx context.Context) error {
	defaultShutdownTimeout := 30 * time.Second

	config := corehttp.ServerConfig{
		Port:            app.Config.HttpPort,
		Logger:          app.Logger,
		ShutdownTimeout: defaultShutdownTimeout,
		ServiceName:     "cart",
	}

	registerFn := func(r *chi.Mux) {
		// Middleware
		r.Use(middleware.Recoverer)
		r.Use(coremiddleware.WithRequestID(app.Logger))

		// Routes
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			app.Logger.Info(ctx, "Handling root request")
			w.Write([]byte("Hello World!"))
		})
	}

	return corehttp.StartServer(ctx, config, registerFn)
}
