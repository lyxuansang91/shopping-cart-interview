package models

import (
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
)

type PartnerStatus string

const (
	PartnerStatusActive   PartnerStatus = "active"
	PartnerStatusInactive PartnerStatus = "inactive"
)

type Partner struct {
	ID          uuid.ID
	PartnerCode string
	Name        string
	GrpcURL     string
	Status      PartnerStatus
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}
