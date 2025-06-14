package repositories

import (
	"context"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/models"
	sqlcmodels "github.com/cinchprotocol/cinch-api/services/cart/internal/pkg/db"
	"github.com/cinchprotocol/cinch-api/services/cart/internal/pkg/mappers"
)

// PaymentMethodRepository defines the interface for payment method operations
type PaymentMethodRepository interface {
	GetByCode(ctx context.Context, code string) (*models.PaymentMethod, error)
	List(ctx context.Context) ([]*models.PaymentMethod, error)
	Enable(ctx context.Context, code string) error
	Disable(ctx context.Context, code string) error
	Delete(ctx context.Context, code string) error
}

// paymentMethodRepository implements PaymentMethodRepository
type paymentMethodRepository struct {
	queries *sqlcmodels.Queries
}

// NewPaymentMethodRepository creates a new payment method repository
func NewPaymentMethodRepository(queries *sqlcmodels.Queries) PaymentMethodRepository {
	return &paymentMethodRepository{
		queries: queries,
	}
}

// GetByCode retrieves a payment method by its code
func (r *paymentMethodRepository) GetByCode(ctx context.Context, code string) (*models.PaymentMethod, error) {
	pm, err := r.queries.GetPaymentMethodByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return mappers.ToDomainPaymentMethod(&pm), nil
}

// List retrieves all payment methods
func (r *paymentMethodRepository) List(ctx context.Context) ([]*models.PaymentMethod, error) {
	pms, err := r.queries.ListPaymentMethods(ctx)
	if err != nil {
		return nil, err
	}
	return mappers.ToDomainPaymentMethods(pms), nil
}

// Enable enables a payment method by its code
func (r *paymentMethodRepository) Enable(ctx context.Context, code string) error {
	return r.queries.EnablePaymentMethodByCode(ctx, code)
}

// Disable disables a payment method by its code
func (r *paymentMethodRepository) Disable(ctx context.Context, code string) error {
	return r.queries.DisablePaymentMethodByCode(ctx, code)
}

// Delete deletes a payment method by its code
func (r *paymentMethodRepository) Delete(ctx context.Context, code string) error {
	return r.queries.DeletePaymentMethodByCode(ctx, code)
}
