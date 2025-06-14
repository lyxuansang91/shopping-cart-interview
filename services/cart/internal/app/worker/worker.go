package worker

import (
	"fmt"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/cinchprotocol/cinch-api/services/cart/internal/app/workflows"
	"github.com/cinchprotocol/cinch-api/services/cart/internal/pkg/interfaces"
)

type Worker struct {
	client         client.Client
	paymentService interfaces.PaymentService
}

func NewWorker(temporalClient client.Client, paymentService interfaces.PaymentService) *Worker {
	return &Worker{
		client:         temporalClient,
		paymentService: paymentService,
	}
}

func (w *Worker) Start() error {
	// Create worker options
	workerOptions := worker.Options{
		MaxConcurrentActivityExecutionSize: 100, // Process up to 100 activities concurrently
	}

	// Create the worker
	temporalWorker := worker.New(w.client, "batch-cart", workerOptions)

	// Register workflow and activities
	batchWorkflow := workflows.NewBatchPaymentWorkflow(w.paymentService)
	temporalWorker.RegisterWorkflow(batchWorkflow.Execute)
	temporalWorker.RegisterActivity(batchWorkflow.ProcessTransaction)

	// Start the worker
	err := temporalWorker.Start()
	if err != nil {
		return fmt.Errorf("failed to start worker: %w", err)
	}

	return nil
}
