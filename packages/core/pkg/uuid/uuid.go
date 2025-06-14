package uuid

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

// ID represents a UUID
type ID string

const Nil = ID("")

// New generates a new UUID
func New() (ID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return ID(id.String()), nil
}

// MustNew generates a new UUID and panics if there's an error
func MustNew() ID {
	id := uuid.Must(uuid.NewV7())
	return ID(id.String())
}

// Parse parses a string into a UUID
func Parse(s string) (ID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return "", err
	}
	return ID(id.String()), nil
}

// MustParse parses a string into a UUID and panics if there's an error
func MustParse(s string) ID {
	id := uuid.Must(uuid.Parse(s))
	return ID(id.String())
}

// IsValid checks if a string is a valid UUID
func IsValid(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

// MarshalBinary encodes the UUID string into a 16-byte slice
func (id ID) MarshalBinary() ([]byte, error) {
	return uuid.MustParse(string(id)).MarshalBinary()
}

// UnmarshalBinary decodes a 16-byte UUID slice into an UUID
func (id *ID) UnmarshalBinary(data []byte) error {
	u, err := uuid.FromBytes(data)
	if err != nil {
		return err
	}
	*id = ID(u.String())
	return nil
}

func (id ID) Value() (driver.Value, error) {
	return id.MarshalBinary()
}

func (id *ID) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan type %T into ID", value)
	}
	return id.UnmarshalBinary(data)
}
