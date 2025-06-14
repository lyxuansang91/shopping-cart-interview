package jobs

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/client"

	"github.com/cinchprotocol/cinch-api/packages/core"
	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
	"github.com/cinchprotocol/cinch-api/services/cart/internal/app"
	"github.com/cinchprotocol/cinch-api/services/cart/internal/app/worker"
	"github.com/cinchprotocol/cinch-api/services/cart/internal/app/workflows"
)

func StartTemporalWorker(ctx context.Context, temporalHost string) error {
	// Use configured Temporal host
	host := temporalHost
	app.Logger.Info(ctx, "Connecting to Temporal server",
		core.NewField("host", host),
	)

	// Add retry logic for Temporal client connection
	var temporalClient client.Client
	var err error
	for i := 0; i < 60; i++ { // Increase to 60 seconds
		temporalClient, err = client.Dial(client.Options{
			HostPort:  host,      // Use the docker service name
			Namespace: "default", // Add explicit namespace
		})
		if err == nil {
			break
		}
		app.Logger.Info(ctx, "Waiting for Temporal server to be ready...",
			core.NewField("attempt", i+1),
			core.NewField("error", err),
		)
		time.Sleep(time.Second)
	}
	if err != nil {
		return fmt.Errorf("failed to create temporal client after 60 attempts: %w", err)
	}
	defer temporalClient.Close()

	// Create and start the worker
	w := worker.NewWorker(temporalClient, app.Services.WorkflowPayment)
	if err := w.Start(); err != nil {
		return fmt.Errorf("failed to start worker: %w", err)
	}

	// Create 10 test transactions
	transactions := make([]workflows.Transaction, 10)
	for i := 0; i < 10; i++ {
		transactions[i] = workflows.Transaction{
			CardToken: fmt.Sprintf("tok_%d", i),
			Amount:    1000 + int64(i), // Different amount for each transaction
			Currency:  "USD",
			Reference: fmt.Sprintf("TX%04d", i),
		}
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        "batch-payment-" + string(uuid.MustNew()),
		TaskQueue: "batch-cart",
	}

	params := workflows.BatchPaymentParams{
		BatchID:      "BATCH-001",
		Transactions: transactions,
	}

	batchWorkflow := workflows.NewBatchPaymentWorkflow(app.Services.WorkflowPayment)
	_, err = temporalClient.ExecuteWorkflow(ctx, workflowOptions, batchWorkflow.Execute, params)
	if err != nil {
		app.Logger.Error(ctx, "Failed to start batch payment workflow",
			core.NewField("error", err),
		)
		// Don't return error here as it's not critical for service startup
	}

	// Block until context is cancelled
	<-ctx.Done()
	return nil
}
