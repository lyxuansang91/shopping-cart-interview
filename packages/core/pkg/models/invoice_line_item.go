package models

import (
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
)

type InvoiceLineItem struct {
	ID          uuid.ID
	InvoiceID   uuid.ID
	OrderID     uuid.ID
	Description string
	Amount      float64
	IsProrated  bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	Relations   InvoiceLineItemRelations
}
