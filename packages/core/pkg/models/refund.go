package models

import (
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
)

// RefundStatus represents the status of a refund
type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "pending"
	RefundStatusCompleted RefundStatus = "completed"
	RefundStatusFailed    RefundStatus = "failed"
)

// Refund represents a refund in the system
type Refund struct {
	ID        uuid.ID
	PaymentID uuid.ID
	Amount    string
	Reason    string
	Status    RefundStatus
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Attempts  []*RefundAttempt
}

// RefundAttempt represents a refund attempt in the system
type RefundAttempt struct {
	ID              uuid.ID
	RefundID        uuid.ID
	PartnerID       uuid.ID
	Status          RefundStatus
	PartnerRefundID string
	ErrorMessage    string
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}
