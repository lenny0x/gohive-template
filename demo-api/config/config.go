package config

import (
	"github.com/gohive/core/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
}

var Cfg *Config

func Load(path string) error {
	cfg, err := config.LoadConfig[Config](path)
	if err != nil {
		return err
	}
	Cfg = cfg
	return nil
}
