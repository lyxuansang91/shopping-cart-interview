package controllers

import (
	"context"

	"github.com/cinchprotocol/cinch-api/packages/proto/pkg/proto/assets/cart"
	"github.com/cinchprotocol/cinch-api/services/cart/internal/pkg/mappers"
	"github.com/cinchprotocol/cinch-api/services/cart/internal/pkg/services/interfaces"
)

// PaymentMethodController handles gRPC requests for payment methods
type PaymentMethodController struct {
	cart.UnimplementedCartServiceServer
	service interfaces.IPaymentMethodService
}

// NewPaymentMethodController creates a new payment method controller
func NewPaymentMethodController(service interfaces.IPaymentMethodService) *PaymentMethodController {
	return &PaymentMethodController{
		service: service,
	}
}

// GetPaymentMethod handles the GetPaymentMethod gRPC request
func (c *PaymentMethodController) GetPaymentMethod(ctx context.Context, req *cart.GetPaymentMethodRequest) (*cart.GetPaymentMethodResponse, error) {
	pm, err := c.service.GetByCode(ctx, req.PaymentMethodCode)
	if err != nil {
		return nil, err
	}

	return &cart.GetPaymentMethodResponse{
		PaymentMethod: mappers.MapDomainToProtoPaymentMethod(pm),
	}, nil
}

// ListPaymentMethods handles the ListPaymentMethods gRPC request
func (c *PaymentMethodController) ListPaymentMethods(ctx context.Context, req *cart.ListPaymentMethodsRequest) (*cart.ListPaymentMethodsResponse, error) {
	pms, err := c.service.List(ctx)
	if err != nil {
		return nil, err
	}

	protoPms := make([]*cart.PaymentMethod, len(pms))
	for i, pm := range pms {
		protoPms[i] = mappers.MapDomainToProtoPaymentMethod(pm)
	}

	return &cart.ListPaymentMethodsResponse{
		PaymentMethods: protoPms,
	}, nil
}

// EnablePaymentMethod handles the EnablePaymentMethod gRPC request
func (c *PaymentMethodController) EnablePaymentMethod(ctx context.Context, req *cart.EnablePaymentMethodRequest) (*cart.EnablePaymentMethodResponse, error) {
	err := c.service.Enable(ctx, req.PaymentMethodCode)
	if err != nil {
		return nil, err
	}

	return &cart.EnablePaymentMethodResponse{}, nil
}

// DisablePaymentMethod handles the DisablePaymentMethod gRPC request
func (c *PaymentMethodController) DisablePaymentMethod(ctx context.Context, req *cart.DisablePaymentMethodRequest) (*cart.DisablePaymentMethodResponse, error) {
	err := c.service.Disable(ctx, req.PaymentMethodCode)
	if err != nil {
		return nil, err
	}

	return &cart.DisablePaymentMethodResponse{}, nil
}

// DeletePaymentMethod handles the DeletePaymentMethod gRPC request
func (c *PaymentMethodController) DeletePaymentMethod(ctx context.Context, req *cart.DeletePaymentMethodRequest) (*cart.DeletePaymentMethodResponse, error) {
	err := c.service.Delete(ctx, req.PaymentMethodCode)
	if err != nil {
		return nil, err
	}

	return &cart.DeletePaymentMethodResponse{}, nil
}
