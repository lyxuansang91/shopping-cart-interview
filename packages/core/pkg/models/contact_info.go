package models

import (
	"errors"
	"time"
)

// ContactInfo represents common contact information fields
// Note: Any changes to KYC information must be communicated to finance for internal updates
// as it may affect compliance and regulatory requirements.
type ContactInfo struct {
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Email     string      `json:"email"`
	Phone     PhoneNumber `json:"phone"`
	Address   Address     `json:"address"`
	UpdatedAt *time.Time  `json:"updated_at"`
}

// Validate performs validation on the contact information
func (c *ContactInfo) Validate() error {
	if c.FirstName == "" {
		return errors.New("first name is required")
	}
	if c.LastName == "" {
		return errors.New("last name is required")
	}
	if c.Email == "" {
		return errors.New("email is required")
	}
	if err := c.Phone.Validate(); err != nil {
		return err
	}
	if err := c.Address.Validate(); err != nil {
		return err
	}
	return nil
}

// HasChanged returns true if the contact information has been updated
func (c *ContactInfo) HasChanged() bool {
	return c.UpdatedAt != nil
}
