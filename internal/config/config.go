package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	HTTP HTTP
	DB   DB
}

type DB struct {
	Path string `env:"DB_PATH"`
}

type HTTP struct {
	Port int `env:"HTTP_PORT"`
}

// New constructs the config.
// variables are populated using the envars and default values.
func New(ctx context.Context) (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process(ctx, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
