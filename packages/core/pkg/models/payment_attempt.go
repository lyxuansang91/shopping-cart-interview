package models

import (
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
)

type PaymentAttemptStatus string

const (
	PaymentAttemptStatusPending           PaymentAttemptStatus = "pending"
	PaymentAttemptStatusPaid              PaymentAttemptStatus = "paid"
	PaymentAttemptStatusFailed            PaymentAttemptStatus = "failed"
	PaymentAttemptStatusPartiallyRefunded PaymentAttemptStatus = "partially_refunded"
	PaymentAttemptStatusRefunded          PaymentAttemptStatus = "refunded"
)

// IsValidStatus checks if the given status is valid
func (s PaymentAttemptStatus) IsValid() bool {
	return s == PaymentAttemptStatusPending || s == PaymentAttemptStatusPaid || s == PaymentAttemptStatusFailed || s == PaymentAttemptStatusPartiallyRefunded || s == PaymentAttemptStatusRefunded
}

type PaymentAttempt struct {
	ID               uuid.ID
	PaymentID        uuid.ID
	PaymentMethodID  uuid.ID
	PartnerID        uuid.ID
	PartnerCode      string
	Metadata         map[string]interface{}
	Status           PaymentAttemptStatus
	PartnerPaymentID string
	RedirectURL      string
	ErrorMessage     string
	CreatedAt        *time.Time
}
