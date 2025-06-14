package models

import (
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
)

// ChangeEntityType represents the type of entity being tracked in change history
type ChangeEntityType string

const (
	ChangeEntityTypeSubscription    ChangeEntityType = "subscription"
	ChangeEntityTypeOrder           ChangeEntityType = "order"
	ChangeEntityTypeItem            ChangeEntityType = "item"
	ChangeEntityTypeInvoice         ChangeEntityType = "invoice"
	ChangeEntityTypeInvoiceLineItem ChangeEntityType = "invoice_line_item"
)

// ChangeType represents the type of change being tracked
type ChangeType string

const (
	ChangeTypeCreate ChangeType = "create"
	ChangeTypeUpdate ChangeType = "update"
	ChangeTypeDelete ChangeType = "delete"
)

type ChangeHistory struct {
	ID         uuid.ID
	EntityType ChangeEntityType
	EntityID   uuid.ID
	ChangeType ChangeType
	ChangedBy  *uuid.ID
	Snapshot   map[string]interface{}
	OccurredAt time.Time
}
