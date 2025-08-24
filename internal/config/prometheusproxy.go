package config

import (
	"errors"
)

type PrometheusProxy struct {
	Url                   string `mapstructure:"url"`
	QueryRangeStepSeconds int    `mapstructure:"queryRangeStepSeconds"`
}

func (cfg *PrometheusProxy) Validate() error {
	if cfg.Url == "" {
		return errors.New("prometheusProxy.url is required")
	}
	if cfg.QueryRangeStepSeconds == 0 {
		return errors.New("prometheusProxy.queryRangeStepSeconds is required")
	}

	return nil
}
