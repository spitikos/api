package config

import "errors"

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

func (s *ServerConfig) Validate() error {
	if s.Port == 0 {
		return errors.New("server.port is required")
	}
	return nil
}
