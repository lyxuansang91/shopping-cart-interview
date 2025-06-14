package models

import (
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
)

// InvoiceStatus represents the current status of an invoice
type InvoiceStatus string

const (
	InvoiceStatusDraft   InvoiceStatus = "draft"
	InvoiceStatusSent    InvoiceStatus = "sent"
	InvoiceStatusPaid    InvoiceStatus = "paid"
	InvoiceStatusOverdue InvoiceStatus = "overdue"
)

type Invoice struct {
	ID                  uuid.ID
	SubscriptionID      uuid.ID
	BillingPeriodStart  time.Time
	BillingPeriodEnd    time.Time
	IssuedDate          time.Time
	DueDate             *time.Time
	TotalAmount         float64
	TemporalExecutionID string
	Status              InvoiceStatus
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
	Relations           InvoiceRelations
}
