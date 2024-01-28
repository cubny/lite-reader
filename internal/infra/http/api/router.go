// Package api lite-reader
//
// Documentation of the lite-reader service.
// It is a service to schedule webhooks.
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//	Host: lite-reader
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package api

import (
	"net/http"

	"github.com/cubny/lite-reader/internal/app/feed"
	"github.com/cubny/lite-reader/internal/app/item"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"

	"github.com/cubny/lite-reader/internal/infra/http/api/middleware"
)

type FeedService interface {
	AddFeed(*feed.AddFeedCommand) (*feed.Feed, error)
	ListFeeds(*feed.ListFeedsCommand) ([]*feed.Feed, error)
	//GetFeed(id string) (*Feed, error)
}

type ItemService interface {
	GetUnreadItems(*item.GetUnreadItemsCommand) ([]*item.Item, error)
	GetStarredItems(*item.GetStarredItemsCommand) ([]*item.Item, error)
	GetFeedItems(*item.GetFeedItemsCommand) ([]*item.Item, error)
}

// Router handles http requests
type Router struct {
	http.Handler
	feedService FeedService
	itemService ItemService
}

// New creates a new handler to handle http requests
func New(itemService ItemService, feedService FeedService) (*Router, error) {
	h := &Router{
		itemService: itemService,
		feedService: feedService,
	}
	router := httprouter.New()

	chain := middleware.NewChain(middleware.ContentTypeJSON)

	router.GET("/health", h.health)
	router.POST("/agg/add", chain.Wrap(h.addFeed))
	router.GET("/feeds", chain.Wrap(h.listFeeds))
	router.GET("/feed/:id", chain.Wrap(h.getFeed))
	router.GET("/items/unread", chain.Wrap(h.getUnreadItems))
	router.GET("/items/starred", chain.Wrap(h.getStarredItems))
	router.GET("/items/feed/:id", chain.Wrap(h.getFeedItems))

	// serve static files for GET /
	router.NotFound = http.FileServer(http.Dir("public"))

	h.Handler = router
	return h, nil
}

func (h *Router) health(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		log.Error("failed to compose body of the response")
	}
}
