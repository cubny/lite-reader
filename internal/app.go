package internal

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/cubny/lite-reader/internal/app/feed"
	"github.com/cubny/lite-reader/internal/app/item"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"github.com/cubny/lite-reader/internal/config"
	"github.com/cubny/lite-reader/internal/infra/http/api"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type App struct {
	ctx context.Context
	cfg *config.Config

	feedService api.FeedService
	itemService api.ItemService
	apiServer   *http.Server

	err error
}

//go:embed infra/sqlite/migrations/*.sql
var embedMigrations embed.FS

func Init(ctx context.Context) (*App, error) {
	a := &App{ctx: ctx}
	a.initConfig()
	a.migrate()
	a.initServices()
	a.initAPIServer()

	return a, a.err
}

func (a *App) ifNoError(fn func() *App) *App {
	if a.err != nil {
		return a
	}
	return fn()
}
func (a *App) initConfig() *App {
	return a.ifNoError(
		func() *App {
			cfg, err := config.New(a.ctx)
			if err != nil {
				log.Fatalf("failed to initiate config: %v", err)
			}

			a.cfg = cfg
			return a
		},
	)
}

func (a *App) migrate() *App {
	return a.ifNoError(func() *App {
		var db *sql.DB
		if db, a.err = sql.Open("sqlite3", a.cfg.DB.Path); a.err != nil {
			a.err = fmt.Errorf("failed to open db: %v", a.err)
			return a
		}

		goose.SetBaseFS(embedMigrations)

		if err := goose.SetDialect("sqlite3"); err != nil {
			a.err = fmt.Errorf("failed to set dialect: %v", err)
			return a
		}

		if err := goose.Up(db, "infra/sqlite/migrations"); err != nil {
			a.err = fmt.Errorf("failed to migrate: %v", err)
		}

		return a
	})
}

func (a *App) initServices() *App {
	return a.ifNoError(func() *App {
		feedService, err := feed.NewService()
		if err != nil {
			a.err = err
			return a
		}

		a.feedService = feedService

		itemService, err := item.NewService()
		if err != nil {
			a.err = err
			return a
		}
		a.itemService = itemService
		return a
	})
}

func (a *App) stopAPIServer() *App {
	log.Info("shutting down HTTP component")
	tctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := a.apiServer.Shutdown(tctx); err != nil {
		a.err = fmt.Errorf("failed to shut down api server, %v", err)
		return a
	}
	log.Infof("api server shut down successfully")
	return a
}

func (a *App) initAPIServer() *App {
	return a.ifNoError(func() *App {
		handler, err := api.New(a.itemService, a.feedService)
		if err != nil {
			a.err = fmt.Errorf("cannot create handler, %v", err)
			return a
		}
		a.apiServer = &http.Server{
			Addr:    fmt.Sprintf(":%d", a.cfg.HTTP.Port),
			Handler: handler,
		}

		go func() {
			log.Infof("starting API server %d", a.cfg.HTTP.Port)
			if err = a.apiServer.ListenAndServe(); err != nil {
				a.err = err
			}
		}()

		return a
	})
}

func (a *App) Stop() error {
	a.err = a.ctx.Err()
	if a.apiServer != nil {
		a.stopAPIServer()
	}
	return nil
}

func WaitTermination() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGTERM, unix.SIGINT)
	<-sigs
}
