package mappers

import "strconv"

// ParseAmount parses a string amount to float64.
// Returns 0 if parsing fails.
func ParseAmount(amount string) float64 {
	value, _ := strconv.ParseFloat(amount, 64)
	return value
}
