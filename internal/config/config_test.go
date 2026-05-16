package config

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWithContext(t *testing.T) {
	got, err := New(context.Background())
	require.NoError(t, err)

	assert.Equal(t, got.HTTP.Port, 3000)
}

func TestNew_DBPath_UnsetUsesUserConfigDir(t *testing.T) {
	t.Setenv("DB_PATH", "")
	// Redirect the per-OS user config dir to a tempdir so the test is hermetic.
	tempHome := t.TempDir()
	switch runtime.GOOS {
	case "windows":
		t.Setenv("AppData", tempHome)
	default:
		t.Setenv("HOME", tempHome)
		t.Setenv("XDG_CONFIG_HOME", filepath.Join(tempHome, ".config"))
	}

	cfg, err := New(context.Background())
	require.NoError(t, err)

	assert.True(t, filepath.IsAbs(cfg.DB.Path), "default path must be absolute, got %q", cfg.DB.Path)
	assert.Equal(t, "agg.db", filepath.Base(cfg.DB.Path))
	assert.Equal(t, "lite-reader", filepath.Base(filepath.Dir(cfg.DB.Path)))
}

func TestNew_DBPath_AbsolutePassthrough(t *testing.T) {
	abs := filepath.Join(t.TempDir(), "custom.db")
	t.Setenv("DB_PATH", abs)

	cfg, err := New(context.Background())
	require.NoError(t, err)
	assert.Equal(t, abs, cfg.DB.Path)
}

func TestNew_DBPath_RelativePassthrough(t *testing.T) {
	t.Setenv("DB_PATH", "data/test-agg.db")

	cfg, err := New(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "data/test-agg.db", cfg.DB.Path,
		"explicit relative DB_PATH must pass through unchanged so callers like make run-test-server keep working")
}
