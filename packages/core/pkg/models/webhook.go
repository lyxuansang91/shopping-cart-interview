package models

import (
	"encoding/json"
	"time"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
)

// Webhook represents a webhook event received from a partner
type Webhook struct {
	ID               uuid.ID         `json:"id"`
	Method           string          `json:"method"`
	URL              string          `json:"url"`
	Headers          json.RawMessage `json:"headers"`
	Payload          json.RawMessage `json:"payload"`
	PartnerID        uuid.ID         `json:"partner_id"`
	PartnerWebhookID string          `json:"partner_webhook_id"`
	PartnerEventType string          `json:"partner_event_type"`
	PartnerPaymentID string          `json:"partner_payment_id"`
	PartnerRefundID  string          `json:"partner_refund_id"`
	ReceivedAt       time.Time       `json:"received_at"`
	ProcessedAt      time.Time       `json:"processed_at"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}
