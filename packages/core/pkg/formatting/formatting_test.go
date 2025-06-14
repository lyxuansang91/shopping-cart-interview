package formatting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPretty(t *testing.T) {
	t.Run("handles nil value", func(t *testing.T) {
		result := Pretty(nil)
		assert.Equal(t, "null", result)
	})

	t.Run("formats simple struct", func(t *testing.T) {
		type TestStruct struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		input := TestStruct{Name: "John", Age: 30}
		result := Pretty(input)

		expected := `{
  "name": "John",
  "age": 30
}`
		assert.Equal(t, expected, result)
	})

	t.Run("formats nested struct", func(t *testing.T) {
		type Address struct {
			Street string `json:"street"`
			City   string `json:"city"`
		}

		type Person struct {
			Name    string  `json:"name"`
			Address Address `json:"address"`
		}

		input := Person{
			Name: "Jane",
			Address: Address{
				Street: "123 Main St",
				City:   "Springfield",
			},
		}

		result := Pretty(input)
		expected := `{
  "name": "Jane",
  "address": {
    "street": "123 Main St",
    "city": "Springfield"
  }
}`
		assert.Equal(t, expected, result)
	})

	t.Run("formats slice", func(t *testing.T) {
		input := []string{"apple", "banana", "cherry"}
		result := Pretty(input)

		expected := `[
  "apple",
  "banana",
  "cherry"
]`
		assert.Equal(t, expected, result)
	})

	t.Run("formats map", func(t *testing.T) {
		input := map[string]interface{}{
			"name":   "test",
			"count":  42,
			"active": true,
		}

		result := Pretty(input)

		// Since map iteration order is not guaranteed, we just check that it contains the expected fields
		assert.Contains(t, result, `"name": "test"`)
		assert.Contains(t, result, `"count": 42`)
		assert.Contains(t, result, `"active": true`)
	})

	t.Run("handles unmarshallable type gracefully", func(t *testing.T) {
		// channels cannot be marshalled to JSON
		ch := make(chan int)
		result := Pretty(ch)

		assert.Contains(t, result, "error marshalling:")
	})

	t.Run("formats pointer to struct", func(t *testing.T) {
		type TestStruct struct {
			Value string `json:"value"`
		}

		input := &TestStruct{Value: "pointer test"}
		result := Pretty(input)

		expected := `{
  "value": "pointer test"
}`
		assert.Equal(t, expected, result)
	})

	t.Run("handles nil pointer", func(t *testing.T) {
		var ptr *string
		result := Pretty(ptr)
		assert.Equal(t, "null", result)
	})

	t.Run("formats primitive types", func(t *testing.T) {
		tests := []struct {
			name     string
			input    interface{}
			expected string
		}{
			{"string", "hello", `"hello"`},
			{"int", 42, "42"},
			{"bool true", true, "true"},
			{"bool false", false, "false"},
			{"float", 3.14, "3.14"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := Pretty(tt.input)
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}
