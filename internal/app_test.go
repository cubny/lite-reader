package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithContext(t *testing.T) {
	runMigration := false
	app, err := Init(context.Background(), runMigration)
	require.NoError(t, err)
	assert.Equal(t, context.Background(), app.ctx)
}

func TestApp_StopOnError(t *testing.T) {
	app := &App{err: fmt.Errorf("bOOm")}
	testFn := func(fnToTest func() *App) func(t *testing.T) {
		return func(t *testing.T) {
			returned := fnToTest()
			assert.Equal(t, app, returned)
		}
	}

	t.Run("initConfig", testFn(app.initConfig))
	t.Run("initAPIServer", testFn(app.initAPIServer))
	t.Run("initSqlClient", testFn(app.initSqlClient))
	t.Run("migrate", testFn(app.migrate))
	t.Run("initRepo", testFn(app.initRepo))
	t.Run("initServices", testFn(app.initServices))
	t.Run("initScheduler", testFn(app.initScheduler))

}
