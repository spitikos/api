package config

import "errors"

type Server struct {
	Port                  int `mapstructure:"port"`
	StreamIntervalSeconds int `mapstructure:"stream_interval_seconds"`
}

func (cfg *Server) Validate() error {
	if cfg.Port == 0 {
		return errors.New("server.port is required")
	}
	if cfg.StreamIntervalSeconds == 0 {
		return errors.New("prometheus_proxy.stream_interval_seconds is required")
	}

	return nil
}
