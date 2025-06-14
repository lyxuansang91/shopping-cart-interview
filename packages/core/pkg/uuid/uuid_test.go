package uuid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("generates valid UUID", func(t *testing.T) {
		id, err := New()
		assert.NoError(t, err)
		assert.NotEmpty(t, id)
		assert.Equal(t, 36, len(id)) // UUID string length is 36 characters
		assert.True(t, IsValid(string(id)))
	})
}

func TestMustNew(t *testing.T) {
	t.Run("generates valid UUID without error", func(t *testing.T) {
		id := MustNew()
		assert.NotEmpty(t, id)
		assert.Equal(t, 36, len(id)) // UUID string length is 36 characters
		assert.True(t, IsValid(string(id)))
	})
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid UUID",
			input:   "123e4567-e89b-12d3-a456-426614174000",
			wantErr: false,
		},
		{
			name:    "invalid UUID",
			input:   "not-a-uuid",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)
				assert.Equal(t, tt.input, string(got))
			}
		})
	}
}

func TestMustParse(t *testing.T) {
	t.Run("valid UUID", func(t *testing.T) {
		validUUID := "123e4567-e89b-12d3-a456-426614174000"
		id := MustParse(validUUID)
		assert.NotEmpty(t, id)
		assert.Equal(t, validUUID, string(id))
	})

	t.Run("panics on invalid UUID", func(t *testing.T) {
		assert.Panics(t, func() {
			MustParse("not-a-uuid")
		})
	})
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid UUID",
			input: "123e4567-e89b-12d3-a456-426614174000",
			want:  true,
		},
		{
			name:  "invalid UUID",
			input: "not-a-uuid",
			want:  false,
		},
		{
			name:  "empty string",
			input: "",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValid(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
