package config

import (
	"errors"
)

type PrometheusProxy struct {
	Url                   string `mapstructure:"url"`
	QueryRangeStepSeconds int    `mapstructure:"query_range_step_seconds"`
}

func (cfg *PrometheusProxy) Validate() error {
	if cfg.Url == "" {
		return errors.New("prometheus_proxy.url is required")
	}
	if cfg.QueryRangeStepSeconds == 0 {
		return errors.New("prometheus_proxy.query_range_step_seconds is required")
	}

	return nil
}
