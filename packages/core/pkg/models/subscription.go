package models

import (
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
)

// BillingInterval represents the interval at which a subscription is billed
type BillingInterval string

const (
	BillingIntervalMonthly   BillingInterval = "monthly"
	BillingIntervalQuarterly BillingInterval = "quarterly"
	BillingIntervalAnnually  BillingInterval = "annually"
)

// BillingStrategy represents how billing periods are calculated
type BillingStrategy string

const (
	// CalendarMonthStrategy bills from last day of previous month to last day of current month
	BillingStrategyCalendarMonth BillingStrategy = "calendar_month"
	// Rolling30DayStrategy bills in 30-day periods from the subscription start date
	BillingStrategyRolling30Day BillingStrategy = "rolling_30_day"
)

// SubscriptionStatus represents the current status of a subscription
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
	SubscriptionStatusPaused    SubscriptionStatus = "paused"
)

type Subscription struct {
	ID                 uuid.ID
	UserID             uuid.ID
	StartDate          time.Time
	EndDate            *time.Time
	BillingInterval    BillingInterval
	BillingStrategy    BillingStrategy
	BillingPeriodUnit  BillingPeriodUnit
	BillingPeriodValue int
	InvoiceDay         int
	FrontLoaded        bool
	DepositAmount      float64
	IsPaygo            bool
	Status             SubscriptionStatus
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time
	Relations          SubscriptionRelations
}
