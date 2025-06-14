package workflows

import (
	"context"
	"fmt"

	"github.com/cinchprotocol/cinch-api/services/adapterstripe/internal/pkg/interfaces"
)

type ProcessTransactionParams struct {
	Transaction Transaction
	BatchID     string
}

type ProcessTransactionResult struct {
	Success bool
	Error   string
}

// ProcessTransaction is an activity that processes a transaction
func (w *BatchPaymentWorkflow) ProcessTransaction(ctx context.Context, params ProcessTransactionParams) (ProcessTransactionResult, error) {
	// Call your payment service to process the transaction
	err := w.paymentService.ProcessPayment(ctx, interfaces.PaymentParams{
		CardToken: params.Transaction.CardToken,
		Amount:    params.Transaction.Amount,
		Currency:  params.Transaction.Currency,
		Reference: fmt.Sprintf("%s-%s", params.BatchID, params.Transaction.Reference),
	})

	if err != nil {
		return ProcessTransactionResult{
			Success: false,
			Error:   err.Error(),
		}, fmt.Errorf("failed to process transaction: %w", err)
	}

	return ProcessTransactionResult{
		Success: true,
	}, nil
}
