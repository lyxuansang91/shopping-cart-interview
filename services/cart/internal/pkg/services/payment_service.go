package services

import (
	"context"
	"fmt"
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core"
	"github.com/cinchprotocol/cinch-api/packages/core/pkg/models"
	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
	"github.com/cinchprotocol/cinch-api/services/cart/internal/pkg/services/interfaces"
)

// PaymentService implements the IPaymentService interface
type PaymentService struct {
	logger core.Logger
}

// NewPaymentService creates a new instance of PaymentService
func NewPaymentService(
	logger core.Logger,
) interfaces.IPaymentService {
	return &PaymentService{
		logger: logger,
	}
}

// CreatePayment implements the CreatePayment method from IPaymentService interface
func (s *PaymentService) CreatePayment(ctx context.Context, payment *models.Payment) (*models.Payment, error) {
	// Generate new ID
	id, err := uuid.New()
	if err != nil {
		return nil, err
	}
	payment.ID = id
	now := time.Now()
	payment.CreatedAt = &now
	payment.UpdatedAt = &now
	payment.Status = models.PaymentStatusPending

	// Log the payment
	s.logger.Info(ctx, "Creating payment",
		core.NewField("payment_id", payment.ID),
		core.NewField("status", payment.Status),
		core.NewField("created_at", payment.CreatedAt))

	// Create a new payment attempt
	attemptID, err := uuid.New()
	if err != nil {
		return nil, fmt.Errorf("failed to generate attempt ID: %w", err)
	}

	// Generate a random UUID for the partner payment ID
	partnerPaymentID, err := uuid.New()
	if err != nil {
		return nil, fmt.Errorf("failed to generate partner payment ID: %w", err)
	}

	attempt := &models.PaymentAttempt{
		ID:               attemptID,
		PaymentID:        payment.ID,
		PaymentMethodID:  payment.PaymentMethodID,
		PartnerCode:      "stripe", // Hardcoded for now since this is the Stripe adapter
		PartnerPaymentID: string(partnerPaymentID),
		Status:           models.PaymentAttemptStatusPaid,
		CreatedAt:        &now,
	}

	// Add the attempt to the payment
	payment.PaymentAttempts = []*models.PaymentAttempt{attempt}

	// TODO: Implement Stripe payment creation
	// 1. Create payment in Stripe
	// 2. Update payment with Stripe payment ID
	// 3. Return updated payment

	s.logger.Info(ctx, "Created payment attempt",
		core.NewField("payment_id", payment.ID),
		core.NewField("attempt_id", attempt.ID),
		core.NewField("status", attempt.Status),
		core.NewField("partner_payment_id", attempt.PartnerPaymentID))

	return payment, nil
}

// UpdatePayment implements the UpdatePayment method from IPaymentService interface
func (s *PaymentService) UpdatePayment(ctx context.Context, payment *models.Payment, webhook models.Webhook) (*models.Payment, error) {
	// Log the payment update
	s.logger.Info(ctx, "Updating payment",
		core.NewField("payment_id", payment.ID),
		core.NewField("webhook_id", webhook.ID),
		core.NewField("partner_webhook_id", webhook.PartnerWebhookID),
		core.NewField("partner_event_type", webhook.PartnerEventType),
	)

	// TODO: Implement Stripe payment update
	// 1. Update payment in Stripe
	// 2. Update payment with Stripe payment ID
	// 3. Return updated payment

	return payment, nil
}
