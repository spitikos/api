package config

import (
	"errors"
)

type Server struct {
	Port                  int `mapstructure:"port"`
	StreamIntervalSeconds int `mapstructure:"streamIntervalSeconds"`
}

func (cfg *Server) Validate() error {
	if cfg.Port == 0 {
		return errors.New("server.port is required")
	}
	if cfg.StreamIntervalSeconds == 0 {
		return errors.New("server.streamIntervalSeconds is required")
	}

	return nil
}
