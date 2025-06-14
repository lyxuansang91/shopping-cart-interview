package models

import (
	"encoding/json"
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
)

// Organisation represents a business entity in the system
type Organisation struct {
	ID              uuid.ID
	Config          json.RawMessage
	Status          OrganisationStatus
	KYCContact      ContactInfo `json:"kyc_contact"`
	ShippingContact ContactInfo `json:"shipping_contact"`
	BillingContact  ContactInfo `json:"billing_contact"`
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
	DeletedAt       *time.Time
}

// NewOrganisation creates a new organisation with the given configuration
func NewOrganisation(config json.RawMessage) (*Organisation, error) {
	id, err := uuid.New()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Organisation{
		ID:        id,
		Config:    config,
		Status:    OrganisationStatusActive,
		CreatedAt: &now,
		UpdatedAt: &now,
	}, nil
}

// IsActive returns true if the organisation is active
func (o *Organisation) IsActive() bool {
	return o.Status == OrganisationStatusActive
}

// IsDeleted returns true if the organisation has been soft deleted
func (o *Organisation) IsDeleted() bool {
	return o.DeletedAt != nil
}

// Enable sets the organisation status to active
func (o *Organisation) Enable() {
	now := time.Now()
	o.Status = OrganisationStatusActive
	o.UpdatedAt = &now
}

// Disable sets the organisation status to inactive
func (o *Organisation) Disable() {
	now := time.Now()
	o.Status = OrganisationStatusInactive
	o.UpdatedAt = &now
}

// Delete performs a soft delete of the organisation
func (o *Organisation) Delete() {
	now := time.Now()
	o.DeletedAt = &now
	o.UpdatedAt = &now
}

// UpdateConfig updates the organisation's configuration
func (o *Organisation) UpdateConfig(config json.RawMessage) {
	now := time.Now()
	o.Config = config
	o.UpdatedAt = &now
}

// UpdateStatus updates the organisation's status
func (o *Organisation) UpdateStatus(status OrganisationStatus) {
	now := time.Now()
	o.Status = status
	o.UpdatedAt = &now
}

// UpdateKYCContact updates the KYC contact information and records the update time
// Note: Changes to KYC information must be communicated to finance for internal updates
func (o *Organisation) UpdateKYCContact(contact ContactInfo) error {
	if err := contact.Validate(); err != nil {
		return err
	}
	now := time.Now()
	contact.UpdatedAt = &now
	o.KYCContact = contact
	o.UpdatedAt = &now
	return nil
}

// UpdateShippingContact updates the shipping contact information
func (o *Organisation) UpdateShippingContact(contact ContactInfo) error {
	if err := contact.Validate(); err != nil {
		return err
	}
	now := time.Now()
	contact.UpdatedAt = &now
	o.ShippingContact = contact
	o.UpdatedAt = &now
	return nil
}

// UpdateBillingContact updates the billing contact information
func (o *Organisation) UpdateBillingContact(contact ContactInfo) error {
	if err := contact.Validate(); err != nil {
		return err
	}
	now := time.Now()
	contact.UpdatedAt = &now
	o.BillingContact = contact
	o.UpdatedAt = &now
	return nil
}

// HasKYCContactChanged returns true if the KYC contact information has been updated
func (o *Organisation) HasKYCContactChanged() bool {
	return o.KYCContact.HasChanged()
}
