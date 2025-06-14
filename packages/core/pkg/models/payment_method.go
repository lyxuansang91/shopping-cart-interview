package models

import (
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
)

type PaymentMethodType string

const (
	PaymentMethodTypeCard PaymentMethodType = "card"
)

// IsValidType checks if the given type is valid
func (t PaymentMethodType) IsValid() bool {
	return t == PaymentMethodTypeCard
}

type PaymentMethodStatus string

const (
	PaymentMethodStatusActive   PaymentMethodStatus = "active"
	PaymentMethodStatusInactive PaymentMethodStatus = "inactive"
)

// IsValidStatus checks if the given status is valid
func (s PaymentMethodStatus) IsValid() bool {
	return s == PaymentMethodStatusActive || s == PaymentMethodStatusInactive
}

type PaymentMethod struct {
	ID                uuid.ID
	PaymentMethodCode string
	Name              string
	Description       string
	PartnerID         uuid.ID
	PartnerCode       string
	PartnerMethodID   string
	Type              PaymentMethodType
	Status            PaymentMethodStatus
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
	DeletedAt         *time.Time
}
