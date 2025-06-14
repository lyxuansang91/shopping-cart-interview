package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Run("successfully parses config with environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("ENV", "test")
		os.Setenv("HTTP_PORT", "8080")
		os.Setenv("GRPC_PORT", "9090")
		os.Setenv("DATABASE_URI", "mysql://user:pass@localhost/db")
		os.Setenv("REDIS_URI", "redis://localhost:6379")
		os.Setenv("LOG_JSON", "true")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("HTTP_PORT")
			os.Unsetenv("GRPC_PORT")
			os.Unsetenv("DATABASE_URI")
			os.Unsetenv("REDIS_URI")
			os.Unsetenv("LOG_JSON")
		}()

		var cfg BaseConfig
		err := Parse(&cfg)

		require.NoError(t, err)
		assert.Equal(t, "test", cfg.Env)
		assert.Equal(t, "8080", cfg.HttpPort)
		assert.Equal(t, "9090", cfg.GrpcPort)
		assert.Equal(t, "mysql://user:pass@localhost/db", cfg.DatabaseURI)
		assert.Equal(t, "redis://localhost:6379", cfg.RedisURI)
		assert.True(t, cfg.LogJSON)
	})

	t.Run("uses default values when environment variables not set", func(t *testing.T) {
		// Ensure required ENV is set but others are unset
		os.Setenv("ENV", "development")
		os.Unsetenv("HTTP_PORT")
		os.Unsetenv("GRPC_PORT")
		os.Unsetenv("DATABASE_URI")
		os.Unsetenv("REDIS_URI")
		os.Unsetenv("LOG_JSON")
		defer os.Unsetenv("ENV")

		var cfg BaseConfig
		err := Parse(&cfg)

		require.NoError(t, err)
		assert.Equal(t, "development", cfg.Env)
		assert.Equal(t, "", cfg.HttpPort)
		assert.Equal(t, "", cfg.GrpcPort)
		assert.Equal(t, "", cfg.DatabaseURI)
		assert.Equal(t, "", cfg.RedisURI)
		assert.False(t, cfg.LogJSON) // default is false
	})

	t.Run("returns error for invalid configuration", func(t *testing.T) {
		// Test with a custom struct that has a truly required field without default
		type InvalidConfig struct {
			BaseConfig
			RequiredField string `env:"REQUIRED_FIELD,required"`
		}

		// Ensure the required field is not set
		os.Unsetenv("REQUIRED_FIELD")
		os.Setenv("ENV", "test")
		defer func() {
			os.Unsetenv("REQUIRED_FIELD")
			os.Unsetenv("ENV")
		}()

		var cfg InvalidConfig
		err := Parse(&cfg)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to load config")
	})

	t.Run("works with custom config struct that embeds BaseConfig", func(t *testing.T) {
		type CustomConfig struct {
			BaseConfig
			CustomField string `env:"CUSTOM_FIELD" envDefault:"default_value"`
		}

		os.Setenv("ENV", "production")
		os.Setenv("CUSTOM_FIELD", "custom_value")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("CUSTOM_FIELD")
		}()

		var cfg CustomConfig
		err := Parse(&cfg)

		require.NoError(t, err)
		assert.Equal(t, "production", cfg.Env)
		assert.Equal(t, "custom_value", cfg.CustomField)
	})
}
