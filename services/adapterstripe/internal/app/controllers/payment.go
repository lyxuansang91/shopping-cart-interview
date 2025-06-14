package controllers

import (
	"context"
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/models"
	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
	"github.com/cinchprotocol/cinch-api/packages/proto/pkg/proto/assets/adapterstripe"
	"github.com/cinchprotocol/cinch-api/services/adapterstripe/internal/pkg/mappers"
	"github.com/cinchprotocol/cinch-api/services/adapterstripe/internal/pkg/services/interfaces"
)

// PaymentController implements the payment-related RPC methods
type PaymentController struct {
	adapterstripe.UnimplementedAdapterstripeServiceServer
	paymentService interfaces.IPaymentService
}

// NewPaymentController creates a new instance of PaymentController
func NewPaymentController(paymentService interfaces.IPaymentService) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
	}
}

// CreatePayment implements the CreatePayment RPC method
func (c *PaymentController) CreatePayment(ctx context.Context, req *adapterstripe.CreatePaymentRequest) (*adapterstripe.CreatePaymentResponse, error) {
	// Convert proto Payment to models Payment
	payment, err := mappers.MapProtoToDomainPayment(req.Payment)
	if err != nil {
		return nil, err
	}

	// Create payment using service
	createdPayment, err := c.paymentService.CreatePayment(ctx, payment)
	if err != nil {
		return nil, err
	}

	// Convert models Payment back to proto Payment
	protoPayment := mappers.MapDomainToProtoPayment(createdPayment)

	return &adapterstripe.CreatePaymentResponse{
		Payment: protoPayment,
	}, nil
}

// UpdatePayment implements the UpdatePayment RPC method
func (c *PaymentController) UpdatePayment(ctx context.Context, req *adapterstripe.UpdatePaymentRequest) (*adapterstripe.UpdatePaymentResponse, error) {
	// Convert proto Payment to models Payment
	payment, err := mappers.MapProtoToDomainPayment(req.Payment)
	if err != nil {
		return nil, err
	}

	// Convert proto Webhook to models Webhook
	webhook := models.Webhook{
		ID:               uuid.MustParse(req.Webhook.Id),
		Method:           req.Webhook.Method,
		URL:              req.Webhook.Url,
		Headers:          req.Webhook.Headers,
		Payload:          req.Webhook.Payload,
		PartnerWebhookID: req.Webhook.PartnerWebhookId,
		PartnerEventType: req.Webhook.PartnerEventType,
		PartnerPaymentID: req.Webhook.PartnerPaymentId,
		ReceivedAt:       time.Now(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Update payment using service
	updatedPayment, err := c.paymentService.UpdatePayment(ctx, payment, webhook)
	if err != nil {
		return nil, err
	}

	// Convert models Payment back to proto Payment
	protoPayment := mappers.MapDomainToProtoPayment(updatedPayment)

	return &adapterstripe.UpdatePaymentResponse{
		Payment: protoPayment,
	}, nil
}
