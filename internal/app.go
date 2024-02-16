package internal

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mmcdole/gofeed"
	"github.com/nikhil1raghav/feedfinder"
	"github.com/pressly/goose/v3"
	log "github.com/sirupsen/logrus"

	"github.com/cubny/lite-reader/internal/app/feed"
	"github.com/cubny/lite-reader/internal/app/item"
	"github.com/cubny/lite-reader/internal/config"
	"github.com/cubny/lite-reader/internal/infra/http/api"
	"github.com/cubny/lite-reader/internal/infra/job"
	feedRepo "github.com/cubny/lite-reader/internal/infra/sqlite/feed"
	itemRepo "github.com/cubny/lite-reader/internal/infra/sqlite/item"
)

type App struct {
	ctx context.Context
	cfg *config.Config

	feedService    api.FeedService
	jobFeedService job.FeedService
	itemService    api.ItemService
	jobItemService job.ItemService
	apiServer      *http.Server
	sqlClient      *sql.DB
	feedRepository feed.Repository
	itemRepository item.Repository
	scheduler      *job.Scheduler

	err error
}

const (
	gracePeriod       = 10 * time.Second
	readHeaderTimeout = 5 * time.Second
)

//go:embed infra/sqlite/migrations/*.sql
var embedMigrations embed.FS

func Init(ctx context.Context, runMigration bool) (*App, error) {
	a := &App{ctx: ctx}

	a.initConfig()
	a.initSQLClient()
	if runMigration {
		a.migrate()
	}
	a.initRepo()
	a.initServices()
	a.initScheduler()
	a.initAPIServer()

	return a, a.err
}

func (a *App) ifNoError(fn func() *App) *App {
	if a.err != nil {
		return a
	}
	return fn()
}
func (a *App) initSQLClient() *App {
	return a.ifNoError(func() *App {
		var sqlClient *sql.DB
		if sqlClient, a.err = sql.Open("sqlite3", a.cfg.DB.Path); a.err != nil {
			a.err = fmt.Errorf("failed to open db: %w", a.err)
			return a
		}
		a.sqlClient = sqlClient

		return a
	})
}

func (a *App) initRepo() *App {
	return a.ifNoError(func() *App {
		a.feedRepository = feedRepo.NewDB(a.sqlClient)
		a.itemRepository = itemRepo.NewDB(a.sqlClient)

		return a
	})
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
		goose.SetBaseFS(embedMigrations)

		if err := goose.SetDialect("sqlite3"); err != nil {
			a.err = fmt.Errorf("failed to set dialect: %w", err)
			return a
		}

		if err := goose.Up(a.sqlClient, "infra/sqlite/migrations"); err != nil {
			a.err = fmt.Errorf("failed to migrate: %w", err)
		}

		return a
	})
}

func (a *App) initServices() *App {
	return a.ifNoError(func() *App {
		finder := feedfinder.NewFeedFinder()
		parser := gofeed.NewParser()
		feedService := feed.NewService(a.feedRepository, parser, finder)
		a.feedService = feedService
		a.jobFeedService = feedService

		itemService := item.NewService(a.itemRepository)
		a.itemService = itemService
		a.jobItemService = itemService
		return a
	})
}

func (a *App) initScheduler() *App {
	return a.ifNoError(func() *App {
		a.scheduler = job.NewScheduler(1 * time.Hour)
		a.scheduler.Start()

		j := job.NewItemsJob(a.jobFeedService, a.jobItemService)
		a.scheduler.Queue <- j
		return a
	})
}

func (a *App) stopAPIServer() *App {
	log.Info("shutting down HTTP component")
	tctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
	defer cancel()
	if err := a.apiServer.Shutdown(tctx); err != nil {
		a.err = fmt.Errorf("failed to shut down api server, %w", err)
		return a
	}
	log.Infof("api server shut down successfully")
	return a
}

func (a *App) initAPIServer() *App {
	return a.ifNoError(func() *App {
		handler, err := api.New(a.itemService, a.feedService)
		if err != nil {
			a.err = fmt.Errorf("cannot create handler, %w", err)
			return a
		}
		a.apiServer = &http.Server{
			Addr:              fmt.Sprintf(":%d", a.cfg.HTTP.Port),
			Handler:           handler,
			ReadHeaderTimeout: readHeaderTimeout,
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
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
