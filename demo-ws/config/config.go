package config

import (
	"github.com/gohive/core/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	WebSocket         WebSocketConfig `mapstructure:"websocket"`
}

type WebSocketConfig struct {
	ReadBufferSize  int      `mapstructure:"read_buffer_size"`
	WriteBufferSize int      `mapstructure:"write_buffer_size"`
	AllowedOrigins  []string `mapstructure:"allowed_origins"`
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
