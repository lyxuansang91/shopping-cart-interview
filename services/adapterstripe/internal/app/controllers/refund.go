package controllers

import (
	"context"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/models"
	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
	"github.com/cinchprotocol/cinch-api/packages/proto/pkg/proto/assets/adapterstripe"
	"github.com/cinchprotocol/cinch-api/services/adapterstripe/internal/pkg/mappers"
	"github.com/cinchprotocol/cinch-api/services/adapterstripe/internal/pkg/services/interfaces"
)

// RefundController implements the refund-related RPC methods
type RefundController struct {
	adapterstripe.UnimplementedAdapterstripeServiceServer
	refundService interfaces.IRefundService
}

// NewRefundController creates a new instance of RefundController
func NewRefundController(refundService interfaces.IRefundService) *RefundController {
	return &RefundController{
		refundService: refundService,
	}
}

// CreateRefund implements the CreateRefund RPC method
func (c *RefundController) CreateRefund(ctx context.Context, req *adapterstripe.CreateRefundRequest) (*adapterstripe.CreateRefundResponse, error) {
	// Convert proto Payment to models Payment
	payment, err := mappers.MapProtoToDomainPayment(req.Payment)
	if err != nil {
		return nil, err
	}

	// Convert proto Refund to models Refund
	refund, err := mappers.MapProtoToDomainRefund(req.Refund)
	if err != nil {
		return nil, err
	}

	// Create refund using service
	createdRefund, err := c.refundService.CreateRefund(ctx, payment, refund)
	if err != nil {
		return nil, err
	}

	// Convert models Refund back to proto Refund
	protoRefund := mappers.MapDomainToProtoRefund(createdRefund)

	return &adapterstripe.CreateRefundResponse{
		Refund: protoRefund,
	}, nil
}

// UpdateRefund implements the UpdateRefund RPC method
func (c *RefundController) UpdateRefund(ctx context.Context, req *adapterstripe.UpdateRefundRequest) (*adapterstripe.UpdateRefundResponse, error) {
	refundID, err := uuid.Parse(req.RefundId)
	if err != nil {
		return nil, err
	}

	// Update refund using service
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

	// Convert models Refund to proto Refund
	protoRefund := mappers.MapDomainToProtoRefund(updatedRefund)

	return &adapterstripe.UpdateRefundResponse{
		Refund: protoRefund,
	}, nil
}
