package controllers

import (
	"context"
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/models"
	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
	"github.com/cinchprotocol/cinch-api/packages/proto/pkg/proto/assets/adapterstripe"
	"github.com/cinchprotocol/cinch-api/services/adapterstripe/internal/pkg/mappers"
	pkginterfaces "github.com/cinchprotocol/cinch-api/services/adapterstripe/internal/pkg/services/interfaces"
)

// AdapterstripeController implements all RPC methods
type AdapterstripeController struct {
	adapterstripe.UnimplementedAdapterstripeServiceServer
	paymentMethodService pkginterfaces.IPaymentMethodService
	paymentService       pkginterfaces.IPaymentService
	refundService        pkginterfaces.IRefundService
}

// NewAdapterstripeController creates a new instance of AdapterstripeController
func NewAdapterstripeController(
	paymentMethodService pkginterfaces.IPaymentMethodService,
	paymentService pkginterfaces.IPaymentService,
	refundService pkginterfaces.IRefundService,
) *AdapterstripeController {
	return &AdapterstripeController{
		paymentMethodService: paymentMethodService,
		paymentService:       paymentService,
		refundService:        refundService,
	}
}

// Payment-related methods

// CreatePayment implements the CreatePayment RPC method
func (c *AdapterstripeController) CreatePayment(ctx context.Context, req *adapterstripe.CreatePaymentRequest) (*adapterstripe.CreatePaymentResponse, error) {
	payment, err := mappers.MapProtoToDomainPayment(req.Payment)
	if err != nil {
		return nil, err
	}

	createdPayment, err := c.paymentService.CreatePayment(ctx, payment)
	if err != nil {
		return nil, err
	}

	protoPayment := mappers.MapDomainToProtoPayment(createdPayment)
	return &adapterstripe.CreatePaymentResponse{
		Payment: protoPayment,
	}, nil
}

// UpdatePayment implements the UpdatePayment RPC method
func (c *AdapterstripeController) UpdatePayment(ctx context.Context, req *adapterstripe.UpdatePaymentRequest) (*adapterstripe.UpdatePaymentResponse, error) {
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

// Refund-related methods

// CreateRefund implements the CreateRefund RPC method
func (c *AdapterstripeController) CreateRefund(ctx context.Context, req *adapterstripe.CreateRefundRequest) (*adapterstripe.CreateRefundResponse, error) {
	payment, err := mappers.MapProtoToDomainPayment(req.Payment)
	if err != nil {
		return nil, err
	}

	refund, err := mappers.MapProtoToDomainRefund(req.Refund)
	if err != nil {
		return nil, err
	}

	createdRefund, err := c.refundService.CreateRefund(ctx, payment, refund)
	if err != nil {
		return nil, err
	}

	protoRefund := mappers.MapDomainToProtoRefund(createdRefund)
	return &adapterstripe.CreateRefundResponse{
		Refund: protoRefund,
	}, nil
}

// UpdateRefund implements the UpdateRefund RPC method
func (c *AdapterstripeController) UpdateRefund(ctx context.Context, req *adapterstripe.UpdateRefundRequest) (*adapterstripe.UpdateRefundResponse, error) {
	refundID, err := uuid.Parse(req.RefundId)
	if err != nil {
		return nil, err
	}

	updatedRefund, err := c.refundService.UpdateRefund(
		ctx,
		refundID,
		req.PartnerRefundId,
		models.RefundStatus(req.Status),
		req.EventType,
		req.EventId,
		req.Metadata,
	)
	if err != nil {
		return nil, err
	}

	protoRefund := mappers.MapDomainToProtoRefund(updatedRefund)
	return &adapterstripe.UpdateRefundResponse{
		Refund: protoRefund,
	}, nil
}

// PaymentMethod-related methods

// GetPaymentMethod implements the GetPaymentMethod RPC method
func (c *AdapterstripeController) GetPaymentMethod(ctx context.Context, req *adapterstripe.GetPaymentMethodRequest) (*adapterstripe.GetPaymentMethodResponse, error) {
	paymentMethod, err := c.paymentMethodService.GetByCode(ctx, req.PaymentMethodCode)
	if err != nil {
		return nil, err
	}

	protoPaymentMethod := mappers.MapDomainToProtoPaymentMethod(paymentMethod)
	return &adapterstripe.GetPaymentMethodResponse{
		PaymentMethod: protoPaymentMethod,
	}, nil
}

// ListPaymentMethods implements the ListPaymentMethods RPC method
func (c *AdapterstripeController) ListPaymentMethods(ctx context.Context, req *adapterstripe.ListPaymentMethodsRequest) (*adapterstripe.ListPaymentMethodsResponse, error) {
	paymentMethods, err := c.paymentMethodService.List(ctx)
	if err != nil {
		return nil, err
	}

	protoPaymentMethods := make([]*adapterstripe.PaymentMethod, len(paymentMethods))
	for i, pm := range paymentMethods {
		protoPaymentMethods[i] = mappers.MapDomainToProtoPaymentMethod(pm)
	}

	return &adapterstripe.ListPaymentMethodsResponse{
		PaymentMethods: protoPaymentMethods,
	}, nil
}

// EnablePaymentMethod implements the EnablePaymentMethod RPC method
func (c *AdapterstripeController) EnablePaymentMethod(ctx context.Context, req *adapterstripe.EnablePaymentMethodRequest) (*adapterstripe.EnablePaymentMethodResponse, error) {
	// TODO: Implement this method
	return nil, nil
}

// DisablePaymentMethod implements the DisablePaymentMethod RPC method
func (c *AdapterstripeController) DisablePaymentMethod(ctx context.Context, req *adapterstripe.DisablePaymentMethodRequest) (*adapterstripe.DisablePaymentMethodResponse, error) {
	// TODO: Implement this method
	return nil, nil
}

// DeletePaymentMethod implements the DeletePaymentMethod RPC method
func (c *AdapterstripeController) DeletePaymentMethod(ctx context.Context, req *adapterstripe.DeletePaymentMethodRequest) (*adapterstripe.DeletePaymentMethodResponse, error) {
	// TODO: Implement this method
	return nil, nil
}
