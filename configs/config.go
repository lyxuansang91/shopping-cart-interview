package configs

import (
	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	Port      string
	BaseURL   string
	LogLevel  string
	EnableCORS bool
}

// GetConfig loads configuration from .env file and environment variables
func GetConfig() *Config {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("BASE_URL", "http://localhost:8080")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("ENABLE_CORS", true)

	// Read .env file if present
	_ = viper.ReadInConfig()

	return &Config{
		Port:      viper.GetString("PORT"),
		BaseURL:   viper.GetString("BASE_URL"),
		LogLevel:  viper.GetString("LOG_LEVEL"),
		EnableCORS: viper.GetBool("ENABLE_CORS"),
	}
} 