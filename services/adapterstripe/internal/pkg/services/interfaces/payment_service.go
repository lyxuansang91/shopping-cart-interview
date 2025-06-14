package interfaces

import (
	"context"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/models"
)

// IPaymentService defines the interface for payment operations
type IPaymentService interface {
	// CreatePayment creates a new payment
	CreatePayment(ctx context.Context, payment *models.Payment) (*models.Payment, error)

	// UpdatePayment updates an existing payment with webhook data
	UpdatePayment(ctx context.Context, payment *models.Payment, webhook models.Webhook) (*models.Payment, error)
}
