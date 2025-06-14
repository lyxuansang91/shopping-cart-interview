package mappers

import (
	"github.com/cinchprotocol/cinch-api/packages/core/pkg/models"
	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
	"github.com/cinchprotocol/cinch-api/packages/proto/pkg/proto/assets/adapterstripe"
)

// MapProtoToDomainPaymentAttempt converts a proto PaymentAttempt to a models PaymentAttempt
func MapProtoToDomainPaymentAttempt(protoAttempt *adapterstripe.PaymentAttempt) (*models.PaymentAttempt, error) {
	if protoAttempt == nil {
		return nil, nil
	}

	paymentID, err := uuid.Parse(protoAttempt.PaymentId)
	if err != nil {
		return nil, err
	}

	paymentMethodID, err := uuid.Parse(protoAttempt.PaymentMethodId)
	if err != nil {
		return nil, err
	}

	partnerID, err := uuid.Parse(protoAttempt.PartnerId)
	if err != nil {
		return nil, err
	}

	createdAt, err := ParseTime(protoAttempt.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &models.PaymentAttempt{
		ID:               uuid.MustParse(protoAttempt.Id),
		PaymentID:        paymentID,
		PaymentMethodID:  paymentMethodID,
		PartnerID:        partnerID,
		PartnerCode:      protoAttempt.PartnerCode,
		Status:           models.PaymentAttemptStatus(protoAttempt.Status),
		PartnerPaymentID: protoAttempt.PartnerPaymentId,
		RedirectURL:      protoAttempt.RedirectUrl,
		ErrorMessage:     protoAttempt.ErrorMessage,
		CreatedAt:        createdAt,
	}, nil
}

// MapDomainToProtoPaymentAttempt converts a models PaymentAttempt to a proto PaymentAttempt
func MapDomainToProtoPaymentAttempt(attempt *models.PaymentAttempt) *adapterstripe.PaymentAttempt {
	if attempt == nil {
		return nil
	}

	return &adapterstripe.PaymentAttempt{
		Id:               string(attempt.ID),
		PaymentId:        string(attempt.PaymentID),
		PaymentMethodId:  string(attempt.PaymentMethodID),
		PartnerId:        string(attempt.PartnerID),
		PartnerCode:      attempt.PartnerCode,
		Status:           string(attempt.Status),
		PartnerPaymentId: attempt.PartnerPaymentID,
		RedirectUrl:      attempt.RedirectURL,
		ErrorMessage:     attempt.ErrorMessage,
		CreatedAt:        FormatTime(attempt.CreatedAt),
	}
}

// MapPaymentAttemptsToProto converts a slice of domain PaymentAttempts to proto PaymentAttempts
func MapPaymentAttemptsToProto(attempts []*models.PaymentAttempt) []*adapterstripe.PaymentAttempt {
	if attempts == nil {
		return nil
	}

	protoAttempts := make([]*adapterstripe.PaymentAttempt, 0, len(attempts))
	for _, attempt := range attempts {
		if attempt == nil {
			continue
		}
		protoAttempts = append(protoAttempts, MapDomainToProtoPaymentAttempt(attempt))
	}
	return protoAttempts
}
