package models

import (
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
)

// PaymentAttempt is defined in payment_attempt.go

type PaymentStatus string

const (
	PaymentStatusPending           PaymentStatus = "pending"
	PaymentStatusCompleted         PaymentStatus = "completed"
	PaymentStatusPaid              PaymentStatus = "paid"
	PaymentStatusFailed            PaymentStatus = "failed"
	PaymentStatusPartiallyRefunded PaymentStatus = "partially_refunded"
	PaymentStatusRefunded          PaymentStatus = "refunded"
)

// IsValidStatus checks if the given status is valid
func (s PaymentStatus) IsValid() bool {
	return s == PaymentStatusPending || s == PaymentStatusCompleted || s == PaymentStatusPaid || s == PaymentStatusFailed || s == PaymentStatusPartiallyRefunded || s == PaymentStatusRefunded
}

type Payment struct {
	ID                   uuid.ID
	PaymentMethodID      uuid.ID
	PaymentMethod        *PaymentMethod
	InvoiceID            uuid.ID
	Amount               string
	Status               PaymentStatus
	PaymentMethodDetails PaymentMethodDetails
	RedirectURL          string
	DueOn                *time.Time
	PaidAt               *time.Time
	CreatedAt            *time.Time
	UpdatedAt            *time.Time
	PaymentAttempts      []*PaymentAttempt
	Refunds              []*Refund
}
