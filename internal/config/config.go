package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	PrometheusUrl         string
	Port                  int
	QueryRangeStepSeconds int
	StreamInvervalSeconds int
}

func New() (*Config, error) {
	v := viper.New()

	v.SetConfigName("prometheus_proxy")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if cfg.PrometheusUrl = v.GetString("prometheus_url"); cfg.PrometheusUrl == "" {
		return nil, fmt.Errorf("prometheus_url is required")
	}
	if cfg.Port = v.GetInt("port"); cfg.Port == 0 {
		return nil, fmt.Errorf("port is required")
	}
	if cfg.QueryRangeStepSeconds = v.GetInt("query_range_step_seconds"); cfg.QueryRangeStepSeconds == 0 {
		return nil, fmt.Errorf("query_range_step_seconds is required")
	}
	if cfg.StreamInvervalSeconds = v.GetInt("stream_interval_seconds"); cfg.StreamInvervalSeconds == 0 {
		return nil, fmt.Errorf("stream_interval_seconds is required")
	}

	return &cfg, nil
}
