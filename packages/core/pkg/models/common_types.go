package models

// BillingPeriodUnit represents the unit of time for billing periods
type BillingPeriodUnit string

const (
	BillingPeriodUnitDay   BillingPeriodUnit = "day"
	BillingPeriodUnitWeek  BillingPeriodUnit = "week"
	BillingPeriodUnitMonth BillingPeriodUnit = "month"
)
