package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/cinchprotocol/cinch-api/packages/core"
	"github.com/cinchprotocol/cinch-api/services/adapterstripe/internal/app"
	"github.com/cinchprotocol/cinch-api/services/adapterstripe/internal/app/jobs"
)

func main() {
	if err := app.Init(); err != nil {
		log.Fatalf("Failed to initialize app: %+v", err)
	}

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create WaitGroup to wait for all servers to shutdown
	var wg sync.WaitGroup

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// Log configuration for debugging
	app.Logger.Info(ctx, "Starting service with configuration",
		core.NewField("env", app.Config.Env),
		core.NewField("database_uri", app.Config.DatabaseURI),
	)

	// Start Temporal worker
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	if err := jobs.StartTemporalWorker(ctx); err != nil && ctx.Err() == nil {
	// 		app.Logger.Fatal(ctx, "Temporal worker failed",
	// 			core.NewField("error", err),
	// 		)
	// 	}
	// }()

	// Start gRPC server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := jobs.StartGRPCServer(ctx); err != nil && ctx.Err() == nil {
			app.Logger.Fatal(ctx, "gRPC server failed",
				core.NewField("error", err),
			)
		}
	}()

	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := jobs.StartHTTPServer(ctx); err != nil && ctx.Err() == nil {
			app.Logger.Fatal(ctx, "HTTP server failed",
				core.NewField("error", err),
			)
		}
	}()

	// Wait for termination signal
	<-sigChan
	app.Logger.Info(ctx, "Received termination signal, initiating shutdown...")

	// Cancel context to initiate shutdown
	cancel()

	// Wait for all servers to complete shutdown
	wg.Wait()
	app.Logger.Info(ctx, "All servers shutdown complete")
}
