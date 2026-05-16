package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	HTTP HTTP
	DB   DB
}

type DB struct {
	Path string `env:"DB_PATH, default=agg.db"`
}

type HTTP struct {
	Port int `env:"HTTP_PORT, default=3000"`
}

// New constructs the config.
// variables are populated using the envars and default values.
func New(ctx context.Context) (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process(ctx, cfg); err != nil {
		return nil, err
	}
	resolved, err := resolveDBPath(cfg.DB.Path)
	if err != nil {
		return nil, err
	}
	cfg.DB.Path = resolved
	return cfg, nil
}

// resolveDBPath returns raw unchanged when absolute; otherwise joins it under
// the OS user config dir (e.g. ~/.config/lite-reader on Linux,
// ~/Library/Application Support/lite-reader on macOS) so a downloaded binary
// works regardless of where the user runs it from.
func resolveDBPath(raw string) (string, error) {
	if filepath.IsAbs(raw) {
		return raw, nil
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("locate user config dir: %w", err)
	}
	return filepath.Join(configDir, "lite-reader", raw), nil
}
