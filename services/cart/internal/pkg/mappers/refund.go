package mappers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/models"
	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
	"github.com/cinchprotocol/cinch-api/packages/proto/pkg/proto/assets/cart"
)

// MapProtoToDomainRefund maps a proto refund to a domain refund
func MapProtoToDomainRefund(refund *cart.Refund) (*models.Refund, error) {
	if refund == nil {
		return nil, nil
	}

	id, err := uuid.Parse(refund.Id)
	if err != nil {
		return nil, err
	}

	paymentID, err := uuid.Parse(refund.PaymentId)
	if err != nil {
		return nil, err
	}

	createdAt, err := time.Parse(time.RFC3339, refund.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := time.Parse(time.RFC3339, refund.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &models.Refund{
		ID:        id,
		PaymentID: paymentID,
		Amount:    fmt.Sprintf("%f", refund.Amount),
		Reason:    refund.Reason,
		Status:    models.RefundStatus(refund.Status),
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}, nil
}

// MapDomainToProtoRefund maps a domain refund to a proto refund
func MapDomainToProtoRefund(refund *models.Refund) *cart.Refund {
	if refund == nil {
		return nil
	}

	amount, _ := strconv.ParseFloat(refund.Amount, 64)

	return &cart.Refund{
		Id:        string(refund.ID),
		PaymentId: string(refund.PaymentID),
		Amount:    amount,
		Reason:    refund.Reason,
		Status:    string(refund.Status),
		CreatedAt: refund.CreatedAt.Format(time.RFC3339),
		UpdatedAt: refund.UpdatedAt.Format(time.RFC3339),
	}
}

// MapProtoToDomainRefundAttempt converts a proto RefundAttempt to a models RefundAttempt
func MapProtoToDomainRefundAttempt(protoAttempt *cart.RefundAttempt) (*models.RefundAttempt, error) {
	refundID, err := uuid.Parse(protoAttempt.RefundId)
	if err != nil {
		return nil, err
	}

	partnerID, err := uuid.Parse(protoAttempt.PartnerId)
	if err != nil {
		return nil, err
	}

	createdAt, err := time.Parse(time.RFC3339, protoAttempt.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := time.Parse(time.RFC3339, protoAttempt.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &models.RefundAttempt{
		ID:              uuid.MustParse(protoAttempt.Id),
		RefundID:        refundID,
		PartnerID:       partnerID,
		Status:          models.RefundStatus(protoAttempt.Status),
		PartnerRefundID: protoAttempt.PartnerRefundId,
		ErrorMessage:    protoAttempt.ErrorMessage,
		CreatedAt:       &createdAt,
		UpdatedAt:       &updatedAt,
	}, nil
}

// MapDomainToProtoRefundAttempt converts a models RefundAttempt to a proto RefundAttempt
func MapDomainToProtoRefundAttempt(attempt *models.RefundAttempt) *cart.RefundAttempt {
	return &cart.RefundAttempt{
		Id:              string(attempt.ID),
		RefundId:        string(attempt.RefundID),
		PartnerId:       string(attempt.PartnerID),
		Status:          string(attempt.Status),
		PartnerRefundId: attempt.PartnerRefundID,
		ErrorMessage:    attempt.ErrorMessage,
		CreatedAt:       attempt.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       attempt.UpdatedAt.Format(time.RFC3339),
	}
}
