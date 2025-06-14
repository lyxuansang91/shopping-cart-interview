package models

import (
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
)

type BillingEventLog struct {
	ID             uuid.ID
	EventType      string
	SubscriptionID *uuid.ID
	OrderID        *uuid.ID
	InvoiceID      *uuid.ID
	Payload        map[string]interface{}
	OccurredAt     time.Time
	CreatedAt      time.Time
}
