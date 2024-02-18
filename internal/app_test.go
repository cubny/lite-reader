package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithContext(t *testing.T) {
	app, err := Init(context.Background(), false)
	if err != nil {
		t.Fatal(err)
	}
	require.NoError(t, err)
	assert.Equal(t, context.Background(), app.ctx)
	stopErr := app.Stop()
	require.NoError(t, stopErr)
}

func TestApp_StopOnError(t *testing.T) {
	app := &App{err: fmt.Errorf("bOOm"), ctx: context.Background()}
	testFn := func(fnToTest func() *App) func(t *testing.T) {
		return func(t *testing.T) {
			returned := fnToTest()
			assert.Equal(t, app, returned)
		}
	}

	t.Run("initConfig", testFn(app.initConfig))
	t.Run("initAPIServer", testFn(app.initAPIServer))
	t.Run("initDBFile", testFn(app.initDBFile))
	t.Run("initSQLClient", testFn(app.initSQLClient))
	t.Run("migrate", testFn(app.migrate))
	t.Run("initRepo", testFn(app.initRepo))
	t.Run("initServices", testFn(app.initServices))
	t.Run("initScheduler", testFn(app.initScheduler))
	stopErr := app.Stop()
	require.NoError(t, stopErr)
}
