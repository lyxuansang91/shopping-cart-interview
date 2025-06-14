package mappers

import (
	"fmt"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/models"
	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
	"github.com/cinchprotocol/cinch-api/packages/proto/pkg/proto/assets/cart"
)

// MapProtoToDomainPayment maps a proto payment to a domain payment
func MapProtoToDomainPayment(payment *cart.Payment) (*models.Payment, error) {
	if payment == nil {
		return nil, nil
	}

	id, err := uuid.Parse(payment.Id)
	if err != nil {
		return nil, err
	}

	paymentMethodID, err := uuid.Parse(payment.PaymentMethodId)
	if err != nil {
		return nil, err
	}

	invoiceID, err := uuid.Parse(payment.InvoiceId)
	if err != nil {
		return nil, err
	}

	createdAt, err := ParseTime(payment.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := ParseTime(payment.UpdatedAt)
	if err != nil {
		return nil, err
	}

	paidAt, err := ParseTime(payment.PaidAt)
	if err != nil {
		return nil, err
	}

	dueOn, err := ParseTime(payment.DueOn)
	if err != nil {
		return nil, err
	}

	return &models.Payment{
		ID:              id,
		PaymentMethodID: paymentMethodID,
		InvoiceID:       invoiceID,
		Amount:          fmt.Sprintf("%f", payment.Amount),
		Status:          models.PaymentStatus(payment.Status),
		DueOn:           dueOn,
		PaidAt:          paidAt,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}, nil
}

// MapDomainToProtoPayment maps a domain payment to a proto payment
func MapDomainToProtoPayment(payment *models.Payment) *cart.Payment {
	if payment == nil {
		return nil
	}

	protoPayment := &cart.Payment{
		Id:              string(payment.ID),
		PaymentMethodId: string(payment.PaymentMethodID),
		InvoiceId:       string(payment.InvoiceID),
		Amount:          ParseAmount(payment.Amount),
		Status:          string(payment.Status),
		DueOn:           FormatTime(payment.DueOn),
		PaidAt:          FormatTime(payment.PaidAt),
		CreatedAt:       FormatTime(payment.CreatedAt),
		UpdatedAt:       FormatTime(payment.UpdatedAt),
	}

	if payment.PaymentAttempts != nil {
		protoPayment.Attempts = MapPaymentAttemptsToProto(payment.PaymentAttempts)
	}

	return protoPayment
}
