package interfaces

import (
	"context"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/models"
	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
)

// IPaymentMethodService defines the interface for payment method operations
type IPaymentMethodService interface {
	// GetByCode gets a payment method by its code
	GetByCode(ctx context.Context, code string) (*models.PaymentMethod, error)

	// List gets all payment methods
	List(ctx context.Context) ([]*models.PaymentMethod, error)

	// Enable enables a payment method
	Enable(ctx context.Context, code string) error

	// Disable disables a payment method
	Disable(ctx context.Context, code string) error

	// Delete deletes a payment method
	Delete(ctx context.Context, code string) error
}

// IRefundService defines the interface for refund-related operations
type IRefundService interface {
	// CreateRefund creates a new refund
	CreateRefund(ctx context.Context, payment *models.Payment, refund *models.Refund) (*models.Refund, error)

	// UpdateRefund updates an existing refund
	UpdateRefund(
		ctx context.Context,
		refundID uuid.ID,
		partnerRefundID string,
		status models.RefundStatus,
		eventType string,
		eventID string,
		metadata map[string]string,
	) (*models.Refund, error)
}
