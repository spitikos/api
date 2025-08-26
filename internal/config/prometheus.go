package config

import (
	"errors"
)

type PrometheusCfg struct {
	Url                   string `mapstructure:"url"`
	QueryRangeStepSeconds int    `mapstructure:"queryRangeStepSeconds"`
}

func (cfg *PrometheusCfg) Validate() error {
	if cfg.Url == "" {
		return errors.New("prometheus.url is required")
	}
	if cfg.QueryRangeStepSeconds == 0 {
		return errors.New("prometheus.queryRangeStepSeconds is required")
	}

	return nil
}
