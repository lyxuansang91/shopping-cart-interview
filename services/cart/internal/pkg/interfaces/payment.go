package interfaces

import "context"

type PaymentService interface {
	ProcessPayment(ctx context.Context, params PaymentParams) error
}

type PaymentParams struct {
	CardToken string
	Amount    int64
	Currency  string
	Reference string
}
