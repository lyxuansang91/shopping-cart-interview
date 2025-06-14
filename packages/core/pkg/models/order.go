package models

import (
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
)

// OrderStatus represents the current status of an order
type OrderStatus string

const (
	OrderStatusActive    OrderStatus = "active"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusCompleted OrderStatus = "completed"
)

// PricingType represents the type of pricing for an order
type PricingType string

const (
	PricingTypeFixed   PricingType = "fixed"
	PricingTypeMonthly PricingType = "monthly"
)

type Order struct {
	ID                uuid.ID
	SubscriptionID    uuid.ID
	ItemID            uuid.ID
	StartDate         time.Time
	EndDate           *time.Time
	Status            OrderStatus
	PricingType       PricingType
	TotalCost         float64
	DepositPaid       float64
	BillingScheduleID *uuid.ID
	IsPaygo           bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
	Relations         OrderRelations
}
