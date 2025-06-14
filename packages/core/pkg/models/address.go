package models

import "errors"

type Address struct {
	Street      string `json:"street"`
	City        string `json:"city"`
	State       string `json:"state"`
	Zip         string `json:"zip"`
	CountryCode string `json:"country_code"`
}

// Validate performs validation on the address
func (a *Address) Validate() error {
	if a.Street == "" {
		return errors.New("street is required")
	}
	if a.City == "" {
		return errors.New("city is required")
	}
	if a.State == "" {
		return errors.New("state is required")
	}
	if a.Zip == "" {
		return errors.New("zip code is required")
	}
	if a.CountryCode == "" {
		return errors.New("country code is required")
	}
	return nil
}
