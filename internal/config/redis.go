package config

import (
	"errors"
)

type RedisCfg struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
}

func (cfg *RedisCfg) Validate() error {
	if cfg.Address == "" {
		return errors.New("redis.address is required")
	}
	// Password can be empty if no password is set for Redis

	return nil
}
