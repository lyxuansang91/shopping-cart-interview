package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type BaseConfig struct {
	Env                string `env:"ENV,required" envDefault:"dev"`
	HttpPort           string `env:"HTTP_PORT"`
	GrpcPort           string `env:"GRPC_PORT"`
	DatabaseURI        string `env:"DATABASE_URI"`
	RedisURI           string `env:"REDIS_URI"`
	LogJSON            bool   `env:"LOG_JSON" envDefault:"false"`
	JaegerEndpoint     string `env:"JAEGER_ENDPOINT" envDefault:"http://jaeger:14268/api/traces"`
	CorsAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS" envDefault:"http://frontend:3001"`
	TemporalHost       string `env:"TEMPORAL_HOST" envDefault:"temporal:7233"`
}

// Parse loads config into any struct that embeds BaseConfig
func Parse[T any](cfg *T) error {
	if err := env.Parse(cfg); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	return nil
}
