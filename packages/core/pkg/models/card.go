package models

import (
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
)

// CardBrand represents the brand/type of a payment card
type CardBrand string

const (
	CardBrandVisa            CardBrand = "visa"
	CardBrandMastercard      CardBrand = "mastercard"
	CardBrandAmericanExpress CardBrand = "amex"
	CardBrandDiscover        CardBrand = "discover"
	CardBrandDinersClub      CardBrand = "diners"
	CardBrandJCB             CardBrand = "jcb"
	CardBrandUnionPay        CardBrand = "unionpay"
)

// CardStatus represents the status of a card
type CardStatus string

const (
	CardStatusPending  CardStatus = "pending"
	CardStatusActive   CardStatus = "active"
	CardStatusInactive CardStatus = "inactive"
)

// CreateCardRequest represents the request to create a new card
type CreateCardRequest struct {
	OrganisationID uuid.ID
}

// CreateCardResponse represents the response from creating a new card
type CreateCardResponse struct {
	Card *Card
}

// Card represents a payment card in the system
type Card struct {
	ID              uuid.ID
	OrganisationID  uuid.ID
	VGSVaultToken   string // Tokenized card holder name
	CardNumberToken string // Tokenized card number
	CardCVCToken    string // Tokenized CVC
	CardExpToken    string // Tokenized expiration date
	LastFour        string // Last 4 digits of card number
	Brand           CardBrand
	ExpiryMonth     int
	ExpiryYear      int
	BindCardURL     string
	Status          CardStatus
	IsDefault       bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewCard creates a new card instance
func NewCard(
	organisationID uuid.ID,
	bindCardURL string,
) *Card {
	now := time.Now()
	return &Card{
		ID:              uuid.MustNew(),
		OrganisationID:  organisationID,
		VGSVaultToken:   "",
		CardNumberToken: "",
		CardCVCToken:    "",
		CardExpToken:    "",
		LastFour:        "",
		Brand:           CardBrandVisa, // Default to Visa
		ExpiryMonth:     0,
		ExpiryYear:      0,
		BindCardURL:     bindCardURL,
		Status:          CardStatusPending,
		IsDefault:       false,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// SetDefault marks the card as the default payment method
func (c *Card) SetDefault() {
	c.IsDefault = true
	c.UpdatedAt = time.Now()
}

// UnsetDefault marks the card as not being the default payment method
func (c *Card) UnsetDefault() {
	c.IsDefault = false
	c.UpdatedAt = time.Now()
}

// IsExpired checks if the card has expired
func (c *Card) IsExpired() bool {
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	// If the expiry year is less than current year, card is expired
	if c.ExpiryYear < currentYear {
		return true
	}

	// If the expiry year is current year and expiry month is less than current month, card is expired
	if c.ExpiryYear == currentYear && c.ExpiryMonth < currentMonth {
		return true
	}

	return false
}
