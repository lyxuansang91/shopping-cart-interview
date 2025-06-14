package formatting

import (
	"encoding/json"
	"fmt"
)

// Pretty prints a struct in a pretty format
func Pretty(v any) string {
	// if v is literally nil, or a nil pointer, we’ll catch that here:
	if v == nil {
		return "null"
	}
	// MarshalIndent into bytes
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		// on error we still return something sensible
		return fmt.Sprintf(`"error marshalling: %v"`, err)
	}
	return string(b)
}
