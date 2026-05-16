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
	Path string `env:"DB_PATH"`
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
	if cfg.DB.Path == "" {
		defaultPath, err := defaultDBPath()
		if err != nil {
			return nil, err
		}
		cfg.DB.Path = defaultPath
	}
	return cfg, nil
}

// defaultDBPath returns the per-OS user-config location for the SQLite db.
// Used only when DB_PATH is unset, so a downloaded binary works regardless of
// where the user runs it from. When DB_PATH is explicitly set we honor it
// as-is (cwd-relative or absolute) to preserve standard CLI semantics.
func defaultDBPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("locate user config dir: %w", err)
	}
	return filepath.Join(configDir, "lite-reader", "agg.db"), nil
}
