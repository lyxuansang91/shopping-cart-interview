package services

import (
	"context"
	"errors"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/models"
	"github.com/cinchprotocol/cinch-api/services/cart/internal/pkg/repositories"
	"github.com/cinchprotocol/cinch-api/services/cart/internal/pkg/services/interfaces"
)

var (
	ErrPaymentMethodNotFound = errors.New("payment method not found")
	ErrInvalidStatus         = errors.New("invalid payment method status")
)

// PaymentMethodService handles business logic for payment methods
type PaymentMethodService struct {
	repo repositories.PaymentMethodRepository
}

// NewPaymentMethodService creates a new payment method service
func NewPaymentMethodService(repo repositories.PaymentMethodRepository) interfaces.IPaymentMethodService {
	return &PaymentMethodService{
		repo: repo,
	}
}

// GetByCode retrieves a payment method by its code
func (s *PaymentMethodService) GetByCode(ctx context.Context, code string) (*models.PaymentMethod, error) {
	pm, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if pm == nil {
		return nil, ErrPaymentMethodNotFound
	}
	return pm, nil
}

// List retrieves all payment methods
func (s *PaymentMethodService) List(ctx context.Context) ([]*models.PaymentMethod, error) {
	return s.repo.List(ctx)
}

// Enable enables a payment method by its code
func (s *PaymentMethodService) Enable(ctx context.Context, code string) error {
	// Verify the payment method exists
	pm, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return err
	}
	if pm == nil {
		return ErrPaymentMethodNotFound
	}

	return s.repo.Enable(ctx, code)
}

// Disable disables a payment method by its code
func (s *PaymentMethodService) Disable(ctx context.Context, code string) error {
	// Verify the payment method exists
	pm, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return err
	}
	if pm == nil {
		return ErrPaymentMethodNotFound
	}

	return s.repo.Disable(ctx, code)
}

// Delete deletes a payment method by its code
func (s *PaymentMethodService) Delete(ctx context.Context, code string) error {
	// Verify the payment method exists
	pm, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return err
	}
	if pm == nil {
		return ErrPaymentMethodNotFound
	}

	return s.repo.Delete(ctx, code)
}
