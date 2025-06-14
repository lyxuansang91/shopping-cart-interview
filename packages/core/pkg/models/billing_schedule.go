package models

import (
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
)

// BillingScheduleStrategy represents the strategy used for billing
type BillingScheduleStrategy string

const (
	BillingScheduleStrategy12Months            BillingScheduleStrategy = "12_months"
	BillingScheduleStrategy12MonthsWithDeposit BillingScheduleStrategy = "12_months_with_deposit"
	BillingScheduleStrategyPaygo               BillingScheduleStrategy = "paygo"
)

type BillingSchedule struct {
	ID               uuid.ID
	Strategy         BillingScheduleStrategy
	IntervalUnit     BillingPeriodUnit
	IntervalValue    int
	MonthlyRate      float64
	ProrationEnabled bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}
