package workflows

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/cinchprotocol/cinch-api/services/adapterstripe/internal/pkg/interfaces"
)

type BatchPaymentWorkflow struct {
	paymentService interfaces.PaymentService
}

type BatchPaymentParams struct {
	Transactions []Transaction
	BatchID      string
}

type Transaction struct {
	CardToken string
	Amount    int64
	Currency  string
	Reference string
}

type BatchProgress struct {
	TotalTransactions  int
	ProcessedCount     int
	FailedTransactions []FailedTransaction
}

type FailedTransaction struct {
	Transaction Transaction
	Error       string
}

func NewBatchPaymentWorkflow(paymentService interfaces.PaymentService) *BatchPaymentWorkflow {
	return &BatchPaymentWorkflow{
		paymentService: paymentService,
	}
}

func (w *BatchPaymentWorkflow) Execute(ctx workflow.Context, params BatchPaymentParams) (BatchProgress, error) {
	logger := workflow.GetLogger(ctx)
	progress := BatchProgress{
		TotalTransactions: len(params.Transactions),
	}

	// Configure batch processing options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 1 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Process transactions in batches of 100
	batchSize := 100
	for i := 0; i < len(params.Transactions); i += batchSize {
		end := i + batchSize
		if end > len(params.Transactions) {
			end = len(params.Transactions)
		}

		batch := params.Transactions[i:end]
		futures := make([]workflow.Future, len(batch))

		// Launch activities in parallel for this batch
		for j, tx := range batch {
			activityParams := interfaces.PaymentParams{
				CardToken: tx.CardToken,
				Amount:    tx.Amount,
				Currency:  tx.Currency,
				Reference: tx.Reference,
			}
			future := workflow.ExecuteActivity(ctx, w.ProcessTransaction, activityParams)
			futures[j] = future
		}

		// Wait for all activities in this batch to complete
		for j, future := range futures {
			var result ProcessTransactionResult
			err := future.Get(ctx, &result)
			progress.ProcessedCount++

			if err != nil {
				logger.Error("Transaction failed", "error", err, "transaction", batch[j])
				progress.FailedTransactions = append(progress.FailedTransactions, FailedTransaction{
					Transaction: batch[j],
					Error:       err.Error(),
				})
			}

			// Log progress
			logger.Info("Processing progress",
				"completed", fmt.Sprintf("%d/%d", progress.ProcessedCount, progress.TotalTransactions),
				"failures", len(progress.FailedTransactions),
			)

			// Update search attributes
			err = workflow.UpsertSearchAttributes(ctx, map[string]interface{}{
				"CustomIntField":     progress.ProcessedCount,                                       // Built-in Int type
				"CustomKeywordField": fmt.Sprintf("Failures: %d", len(progress.FailedTransactions)), // Built-in Keyword type
			})
			if err != nil {
				logger.Error("Failed to update search attributes", "error", err)
			}
		}
	}

	return progress, nil
}
