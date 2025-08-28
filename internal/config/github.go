package config

import (
	"errors"
)

type GithubCfg struct {
	Owner string `mapstructure:"owner"`
	Token string `mapstructure:"token"`
}

func (cfg *GithubCfg) Validate() error {
	if cfg.Owner == "" {
		return errors.New("github.owner is required")
	}

	if cfg.Token == "" {
		return errors.New("GITHUB_TOKEN environment variable is required")
	}

	return nil
}
