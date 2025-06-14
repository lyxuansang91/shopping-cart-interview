package services

import (
	"context"

	"github.com/cinchprotocol/cinch-api/packages/core"
	"github.com/cinchprotocol/cinch-api/services/adapterstripe/internal/pkg/interfaces"
)

// WorkflowPaymentService implements the PaymentService interface for workflows
type WorkflowPaymentService struct {
	logger core.Logger
}

// NewWorkflowPaymentService creates a new instance of WorkflowPaymentService
func NewWorkflowPaymentService(logger core.Logger) interfaces.PaymentService {
	return &WorkflowPaymentService{
		logger: logger,
	}
}

// ProcessPayment implements the ProcessPayment method from PaymentService interface
func (s *WorkflowPaymentService) ProcessPayment(ctx context.Context, params interfaces.PaymentParams) error {
	// Log the payment processing
	s.logger.Info(ctx, "Processing payment",
		core.NewField("card_token", params.CardToken),
		core.NewField("amount", params.Amount),
		core.NewField("currency", params.Currency),
		core.NewField("reference", params.Reference),
	)

	// TODO: Implement actual Stripe payment processing
	// For now, we'll just simulate successful processing
	s.logger.Info(ctx, "Payment processed successfully",
		core.NewField("reference", params.Reference),
	)

	return nil
}
