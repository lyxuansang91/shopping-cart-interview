package models

import (
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
)

// PaymentMethodDetails contains all possible payment method specific details
type PaymentMethodDetails struct {
	General     *GeneralDetails
	Card        *CardDetails
	EWallet     *EWalletDetails
	PromptPay   *PromptPayDetails
	PayNow      *PayNowDetails
	DuitNow     *DuitNowDetails
	DirectDebit *DirectDebitDetails
	FPX         *FPXDetails
}

type GeneralDetails struct {
	OnSuccessURL string
	OnFailureURL string
}

// CardDetails represents credit/debit card payment details
type CardDetails struct {
	CardID uuid.ID
}

// GrabPayDetails represents GrabPay payment details
type EWalletDetails struct {
	EWalletID uuid.ID
}

// PromptPayDetails represents PromptPay payment details
type PromptPayDetails struct {
	PromptPayID   string    // PromptPay ID/number
	QRCodeURL     string    // URL to the QR code image
	QRCodeContent string    // Raw QR code content
	ExpiryDate    time.Time // When the QR code expires
}

// PayNowDetails represents PayNow payment details
type PayNowDetails struct {
	PayNowID      string    // PayNow ID/number
	QRCodeURL     string    // URL to the QR code image
	QRCodeContent string    // Raw QR code content
	ExpiryDate    time.Time // When the QR code expires
}

// DuitNowDetails represents DuitNow payment details
type DuitNowDetails struct {
	DuitNowID     string    // DuitNow ID/number
	QRCodeURL     string    // URL to the QR code image
	QRCodeContent string    // Raw QR code content
	ExpiryDate    time.Time // When the QR code expires
}

// DirectDebitDetails represents virtual account payment details
type DirectDebitDetails struct {
	AccountNumber    string    // Virtual account number
	AccountFirstName string    // Account holder name
	AccountLastName  string    // Account holder name
	ExpiryDate       time.Time // When the virtual account expires
}

// FPXDetails represents FPX (Financial Process Exchange) payment details
type FPXDetails struct {
	BankCode      string // Bank code
	BankName      string // Bank name
	AccountNumber string // Bank account number
	AccountName   string // Account holder name
}

// PaymentMethodOutputs contains all possible payment method specific output details
type PaymentMethodOutputs struct {
	Card           *CardOutputs
	GrabPay        *GrabPayOutputs
	PromptPay      *PromptPayOutputs
	PayNow         *PayNowOutputs
	DuitNow        *DuitNowOutputs
	VirtualAccount *VirtualAccountOutputs
	FPX            *FPXOutputs
	ShopBack       *ShopBackOutputs
}

// CardOutputs represents credit/debit card payment output details
type CardOutputs struct {
	PaymentID   string
	Status      string
	Amount      float64
	Currency    string
	LastFour    string
	Brand       CardBrand
	ExpiryMonth int
	ExpiryYear  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GrabPayOutputs represents GrabPay payment output details
type GrabPayOutputs struct {
	PaymentID string
	Status    string
	Amount    float64
	Currency  string
	AccountID string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PromptPayOutputs represents PromptPay payment output details
type PromptPayOutputs struct {
	PaymentID     string
	Status        string
	Amount        float64
	Currency      string
	PromptPayID   string
	QRCodeURL     string
	QRCodeContent string
	ExpiryDate    time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// PayNowOutputs represents PayNow payment output details
type PayNowOutputs struct {
	PaymentID     string
	Status        string
	Amount        float64
	Currency      string
	PayNowID      string
	QRCodeURL     string
	QRCodeContent string
	ExpiryDate    time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// DuitNowOutputs represents DuitNow payment output details
type DuitNowOutputs struct {
	PaymentID     string
	Status        string
	Amount        float64
	Currency      string
	DuitNowID     string
	QRCodeURL     string
	QRCodeContent string
	ExpiryDate    time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// VirtualAccountOutputs represents virtual account payment output details
type VirtualAccountOutputs struct {
	PaymentID     string
	Status        string
	Amount        float64
	Currency      string
	BankCode      string
	AccountNumber string
	AccountName   string
	ExpiryDate    time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// FPXOutputs represents FPX payment output details
type FPXOutputs struct {
	PaymentID     string
	Status        string
	Amount        float64
	Currency      string
	BankCode      string
	BankName      string
	AccountNumber string
	AccountName   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ShopBackOutputs represents ShopBack payment output details
type ShopBackOutputs struct {
	PaymentID string
	Status    string
	Amount    float64
	Currency  string
	AccountID string
	CreatedAt time.Time
	UpdatedAt time.Time
}
