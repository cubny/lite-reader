// Package api lite-reader
//
// Documentation of the lite-reader service.
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

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"

	"github.com/cubny/lite-reader/internal/app/auth"
	"github.com/cubny/lite-reader/internal/app/feed"
	"github.com/cubny/lite-reader/internal/app/item"
	"github.com/cubny/lite-reader/internal/infra/http/api/middleware"
)

type FeedService interface {
	AddFeed(command *feed.AddFeedCommand) (*feed.Feed, error)
	ListFeeds() ([]*feed.Feed, error)
	FetchItems(int) ([]*item.Item, error)
	DeleteFeed(command *feed.DeleteFeedCommand) error
}

type ItemService interface {
	GetUnreadItems() ([]*item.Item, error)
	GetStarredItems() ([]*item.Item, error)
	GetFeedItems(*item.GetFeedItemsCommand) ([]*item.Item, error)
	UpsertItems(command *item.UpsertItemsCommand) error
	UpdateItem(*item.UpdateItemCommand) error
	ReadFeedItems(*item.ReadFeedItemsCommand) error
	UnreadFeedItems(*item.UnreadFeedItemsCommand) error
	GetStarredItemsCount() (int, error)
	GetUnreadItemsCount() (int, error)
	DeleteFeedItems(*item.DeleteFeedItemsCommand) error
}

type AuthService interface {
	Login(command *auth.LoginCommand) (string, error)
}

// Router handles http requests
type Router struct {
	http.Handler
	feedService FeedService
	itemService ItemService
	authService AuthService
}

// New creates a new handler to handle http requests
func New(itemService ItemService, feedService FeedService, authService AuthService) (*Router, error) {
	h := &Router{
		itemService: itemService,
		feedService: feedService,
		authService: authService,
	}
	router := httprouter.New()

	chain := middleware.NewChain(middleware.ContentTypeJSON)

	router.GET("/health", h.health)

	router.GET("/feeds", chain.Wrap(h.listFeeds))
	router.POST("/feeds", chain.Wrap(h.addFeed))

	router.DELETE("/feeds/:id", chain.Wrap(h.deleteFeed))
	router.PUT("/feeds/:id/fetch", chain.Wrap(h.fetchFeedNewItems))
	router.POST("/feeds/:id/read", chain.Wrap(h.readFeedItems))
	router.POST("/feeds/:id/unread", chain.Wrap(h.unreadFeedItems))
	router.GET("/feeds/:id/items", chain.Wrap(h.getFeedItems))

	router.PUT("/items/:id", chain.Wrap(h.updateItem))
	router.GET("/items/unread", chain.Wrap(h.getUnreadItems))
	router.GET("/items/starred", chain.Wrap(h.getStarredItems))
	router.GET("/items/unread/count", chain.Wrap(h.getUnreadItemsCount))
	router.GET("/items/starred/count", chain.Wrap(h.getStarredItemsCount))

	router.POST("/login", chain.Wrap(h.login))
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
