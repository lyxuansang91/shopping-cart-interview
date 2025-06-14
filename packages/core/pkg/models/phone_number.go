package models

import "errors"

type PhoneNumber struct {
	CountryCode string `json:"country_code"`
	Number      string `json:"number"`
}

// Validate performs validation on the phone number
func (p *PhoneNumber) Validate() error {
	if p.CountryCode == "" {
		return errors.New("country code is required")
	}
	if p.Number == "" {
		return errors.New("phone number is required")
	}
	return nil
}
