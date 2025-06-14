package controllers

import (
	"context"

	"github.com/cinchprotocol/cinch-api/packages/proto/pkg/proto/assets/adapterstripe"
	"github.com/cinchprotocol/cinch-api/services/adapterstripe/internal/pkg/mappers"
	"github.com/cinchprotocol/cinch-api/services/adapterstripe/internal/pkg/services/interfaces"
)

// PaymentMethodController handles gRPC requests for payment methods
type PaymentMethodController struct {
	adapterstripe.UnimplementedAdapterstripeServiceServer
	service interfaces.IPaymentMethodService
}

// NewPaymentMethodController creates a new payment method controller
func NewPaymentMethodController(service interfaces.IPaymentMethodService) *PaymentMethodController {
	return &PaymentMethodController{
		service: service,
	}
}

// GetPaymentMethod handles the GetPaymentMethod gRPC request
func (c *PaymentMethodController) GetPaymentMethod(ctx context.Context, req *adapterstripe.GetPaymentMethodRequest) (*adapterstripe.GetPaymentMethodResponse, error) {
	pm, err := c.service.GetByCode(ctx, req.PaymentMethodCode)
	if err != nil {
		return nil, err
	}

	return &adapterstripe.GetPaymentMethodResponse{
		PaymentMethod: mappers.MapDomainToProtoPaymentMethod(pm),
	}, nil
}

// ListPaymentMethods handles the ListPaymentMethods gRPC request
func (c *PaymentMethodController) ListPaymentMethods(ctx context.Context, req *adapterstripe.ListPaymentMethodsRequest) (*adapterstripe.ListPaymentMethodsResponse, error) {
	pms, err := c.service.List(ctx)
	if err != nil {
		return nil, err
	}

	protoPms := make([]*adapterstripe.PaymentMethod, len(pms))
	for i, pm := range pms {
		protoPms[i] = mappers.MapDomainToProtoPaymentMethod(pm)
	}

	return &adapterstripe.ListPaymentMethodsResponse{
		PaymentMethods: protoPms,
	}, nil
}

// EnablePaymentMethod handles the EnablePaymentMethod gRPC request
func (c *PaymentMethodController) EnablePaymentMethod(ctx context.Context, req *adapterstripe.EnablePaymentMethodRequest) (*adapterstripe.EnablePaymentMethodResponse, error) {
	err := c.service.Enable(ctx, req.PaymentMethodCode)
	if err != nil {
		return nil, err
	}

	return &adapterstripe.EnablePaymentMethodResponse{}, nil
}

// DisablePaymentMethod handles the DisablePaymentMethod gRPC request
func (c *PaymentMethodController) DisablePaymentMethod(ctx context.Context, req *adapterstripe.DisablePaymentMethodRequest) (*adapterstripe.DisablePaymentMethodResponse, error) {
	err := c.service.Disable(ctx, req.PaymentMethodCode)
	if err != nil {
		return nil, err
	}

	return &adapterstripe.DisablePaymentMethodResponse{}, nil
}

// DeletePaymentMethod handles the DeletePaymentMethod gRPC request
func (c *PaymentMethodController) DeletePaymentMethod(ctx context.Context, req *adapterstripe.DeletePaymentMethodRequest) (*adapterstripe.DeletePaymentMethodResponse, error) {
	err := c.service.Delete(ctx, req.PaymentMethodCode)
	if err != nil {
		return nil, err
	}

	return &adapterstripe.DeletePaymentMethodResponse{}, nil
}
