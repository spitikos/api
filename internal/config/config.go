package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerCfg     `mapstructure:"server"`
	Prometheus PrometheusCfg `mapstructure:"prometheus"`
	Redis      RedisCfg      `mapstructure:"redis"`
	Github     GithubCfg     `mapstructure:"github"`
}

type Validatable interface {
	Validate() error
}

var envs = []string{
	"github.token",
}

func New() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")

	// bind env variables
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	for _, env := range envs {
		if err := v.BindEnv(env); err != nil {
			return nil, fmt.Errorf("failed to bind env variable %s: %w", env, err)
		}
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	toValidate := []Validatable{
		&cfg.Server,
		&cfg.Prometheus,
		&cfg.Redis,
		&cfg.Github,
	}

	for _, c := range toValidate {
		if err := c.Validate(); err != nil {
			name := reflect.TypeOf(c).Elem().Name()
			return nil, fmt.Errorf("failed to validate %s: %w", name, err)
		}
	}

	return &cfg, nil
}
