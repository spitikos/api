package config

import (
	"errors"
)

type PrometheusCfg struct {
	Address               string `mapstructure:"address"`
	QueryRangeStepSeconds int    `mapstructure:"queryRangeStepSeconds"`
}

func (cfg *PrometheusCfg) Validate() error {
	if cfg.Address == "" {
		return errors.New("prometheus.address is required")
	}
	if cfg.QueryRangeStepSeconds == 0 {
		return errors.New("prometheus.queryRangeStepSeconds is required")
	}

	return nil
}
