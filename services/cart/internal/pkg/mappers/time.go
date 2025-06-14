package mappers

import "time"

// ParseTime parses a time string in RFC3339 format and returns a pointer to the time.
// Returns nil if the input string is empty.
func ParseTime(timeStr string) (*time.Time, error) {
	if timeStr == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// FormatTime formats a time pointer to RFC3339 format.
// Returns an empty string if the input is nil.
func FormatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
