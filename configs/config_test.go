package configs

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetConfig_Defaults(t *testing.T) {
	// Clear env and viper
	os.Clearenv()
	viper.Reset()

	config := GetConfig()
	assert.Equal(t, "8080", config.Port)
	assert.Equal(t, "http://localhost:8080", config.BaseURL)
	assert.Equal(t, "info", config.LogLevel)
	assert.True(t, config.EnableCORS)
}

func TestGetConfig_EnvOverride(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("BASE_URL", "https://example.com")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("ENABLE_CORS", "false")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("BASE_URL")
	defer os.Unsetenv("LOG_LEVEL")
	defer os.Unsetenv("ENABLE_CORS")
	viper.Reset()

	config := GetConfig()
	assert.Equal(t, "9090", config.Port)
	assert.Equal(t, "https://example.com", config.BaseURL)
	assert.Equal(t, "debug", config.LogLevel)
	assert.False(t, config.EnableCORS)
}

func TestGetConfig_DotEnvFile(t *testing.T) {
	// Write a temporary .env file
	dotenv := []byte("PORT=1234\nBASE_URL=https://dotenv.com\nLOG_LEVEL=warn\nENABLE_CORS=false\n")
	os.WriteFile(".env", dotenv, 0644)
	defer os.Remove(".env")
	os.Clearenv()
	viper.Reset()

	config := GetConfig()
	assert.Equal(t, "1234", config.Port)
	assert.Equal(t, "https://dotenv.com", config.BaseURL)
	assert.Equal(t, "warn", config.LogLevel)
	assert.False(t, config.EnableCORS)
} 