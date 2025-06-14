package configs

import (
	"github.com/cinchprotocol/cinch-api/packages/core/pkg/config"
)

type Config struct {
	config.BaseConfig
	Temporal TemporalConfig `mapstructure:"temporal"`
}

type TemporalConfig struct {
	Host string `mapstructure:"host"`
}
