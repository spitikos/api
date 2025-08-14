package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server          ServerConfig          `mapstructure:"server"`
	PrometheusProxy PrometheusProxyConfig `mapstructure:"prometheus_proxy"`
}

func New() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.Server.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate server config: %w", err)
	}

	if err := cfg.PrometheusProxy.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate prometheus proxy config: %w", err)
	}

	return &cfg, nil
}
