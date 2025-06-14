package models

import (
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
)

// EWalletBrand represents the brand/type of an e-wallet
type EWalletBrand string

const (
	EWalletBrandGrabPay  EWalletBrand = "grabpay"
	EWalletBrandShopBack EWalletBrand = "shopback"
	EWalletBrandShopee   EWalletBrand = "shopee"
)

// EWalletStatus represents the status of an e-wallet
type EWalletStatus string

const (
	EWalletStatusPending  EWalletStatus = "pending"
	EWalletStatusActive   EWalletStatus = "active"
	EWalletStatusInactive EWalletStatus = "inactive"
)

// EWalletOptions contains configuration options for an e-wallet
type EWalletOptions struct {
	SuccessURL string // Success URL for the e-wallet
	FailureURL string // Failure URL for the e-wallet
}

// EWallet represents an e-wallet in the system
type EWallet struct {
	ID               uuid.ID
	OrganisationID   uuid.ID
	Brand            EWalletBrand
	AccountID        string // E-wallet account identifier
	Email            string // Associated email address
	PhoneNumber      string // Associated phone number
	PartnerEWalletID string // Partner e-wallet account identifier
	PartnerToken     string // Partner e-wallet token
	RedirectURL      string // Redirect URL for the e-wallet
	Options          *EWalletOptions
	Status           EWalletStatus
	IsDefault        bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NewEWallet creates a new e-wallet instance
func NewEWallet(
	organisationID uuid.ID,
	brand EWalletBrand,
) *EWallet {
	now := time.Now()
	return &EWallet{
		ID:             uuid.MustNew(),
		OrganisationID: organisationID,
		Brand:          brand,
		AccountID:      "",
		Email:          "",
		PhoneNumber:    "",
		Status:         EWalletStatusPending,
		IsDefault:      false,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// SetDefault marks the e-wallet as the default payment method
func (e *EWallet) SetDefault() {
	e.IsDefault = true
	e.UpdatedAt = time.Now()
}

// UnsetDefault marks the e-wallet as not being the default payment method
func (e *EWallet) UnsetDefault() {
	e.IsDefault = false
	e.UpdatedAt = time.Now()
}

// Activate marks the e-wallet as active
func (e *EWallet) Activate() {
	e.Status = EWalletStatusActive
	e.UpdatedAt = time.Now()
}

// Deactivate marks the e-wallet as inactive
func (e *EWallet) Deactivate() {
	e.Status = EWalletStatusInactive
	e.UpdatedAt = time.Now()
}

// IsActive checks if the e-wallet is active
func (e *EWallet) IsActive() bool {
	return e.Status == EWalletStatusActive
}
