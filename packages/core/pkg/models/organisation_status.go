package models

// OrganisationStatus represents the possible states of an organisation
type OrganisationStatus string

const (
	// OrganisationStatusActive indicates the organisation is active
	OrganisationStatusActive OrganisationStatus = "active"
	// OrganisationStatusInactive indicates the organisation is inactive
	OrganisationStatusInactive OrganisationStatus = "inactive"
)
