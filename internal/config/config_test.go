package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWithContext(t *testing.T) {
	got, err := New(context.Background())
	require.NoError(t, err)

	assert.Equal(t, got.HTTP.Port, 3000)
}
