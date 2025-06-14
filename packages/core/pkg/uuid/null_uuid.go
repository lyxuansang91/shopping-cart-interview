package uuid

import (
	"database/sql/driver"
)

// NullUUID represents a nullable UUID
// Implements sql.Scanner and driver.Valuer
// Used for nullable BINARY(16) columns
type NullUUID struct {
	UUID  ID
	Valid bool
}

func (n *NullUUID) Scan(value interface{}) error {
	if value == nil {
		n.UUID, n.Valid = "", false
		return nil
	}
	var id ID
	if err := id.Scan(value); err != nil {
		return err
	}
	n.UUID, n.Valid = id, true
	return nil
}

func (n NullUUID) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.UUID.Value()
}

// MarshalBinary implements encoding.BinaryMarshaler interface
func (n NullUUID) MarshalBinary() ([]byte, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.UUID.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface
func (n *NullUUID) UnmarshalBinary(data []byte) error {
	if data == nil {
		n.UUID, n.Valid = "", false
		return nil
	}
	var id ID
	if err := id.UnmarshalBinary(data); err != nil {
		return err
	}
	n.UUID, n.Valid = id, true
	return nil
}
