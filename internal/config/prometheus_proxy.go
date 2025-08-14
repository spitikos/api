package config

import (
	"errors"
)

type PrometheusProxyConfig struct {
	Url                   string `mapstructure:"url"`
	QueryRangeStepSeconds int    `mapstructure:"query_range_step_seconds"`
	StreamIntervalSeconds int    `mapstructure:"stream_interval_seconds"`
}

func (cfg *PrometheusProxyConfig) Validate() error {
	if cfg.Url == "" {
		return errors.New("prometheus_proxy.url is required")
	}
	if cfg.QueryRangeStepSeconds == 0 {
		return errors.New("prometheus_proxy.query_range_step_seconds is required")
	}
	if cfg.StreamIntervalSeconds == 0 {
		return errors.New("prometheus_proxy.stream_interval_seconds is required")
	}
	return nil
}
