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
	"errors"
	"io/fs"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"

	"github.com/cubny/lite-reader/internal/app/auth"
	"github.com/cubny/lite-reader/internal/app/feed"
	"github.com/cubny/lite-reader/internal/app/folder"
	"github.com/cubny/lite-reader/internal/app/item"
	"github.com/cubny/lite-reader/internal/infra/http/api/middleware"
)

type FeedService interface {
	AddFeed(command *feed.AddFeedCommand) (*feed.Feed, error)
	ListFeeds(UserID int) ([]*feed.Feed, error)
	FetchItems(int) ([]*item.Item, error)
	DeleteFeed(command *feed.DeleteFeedCommand) error
	MoveFeed(command *feed.MoveFeedCommand) error
	ReorderFeed(command *feed.ReorderFeedCommand) error
	BulkMoveFeeds(command *feed.BulkMoveFeedsCommand) error
}

type FolderService interface {
	AddFolder(command *folder.AddFolderCommand) (*folder.Folder, error)
	ListFolders(userID int) ([]*folder.Folder, error)
	RenameFolder(command *folder.RenameFolderCommand) error
	ReorderFolder(command *folder.ReorderFolderCommand) error
	DeleteFolder(command *folder.DeleteFolderCommand) error
	GetFolder(command *folder.GetFolderCommand) (*folder.Folder, error)
}

type ItemService interface {
	GetUnreadItems() ([]*item.Item, error)
	GetStarredItems() ([]*item.Item, error)
	GetFeedItems(*item.GetFeedItemsCommand) ([]*item.Item, error)
	GetFolderItems(*item.GetFolderItemsCommand) ([]*item.Item, error)
	UpsertItems(command *item.UpsertItemsCommand) error
	UpdateItem(*item.UpdateItemCommand) error
	ReadFeedItems(*item.ReadFeedItemsCommand) error
	UnreadFeedItems(*item.UnreadFeedItemsCommand) error
	ReadFolderItems(*item.ReadFolderItemsCommand) error
	GetStarredItemsCount() (int, error)
	GetUnreadItemsCount() (int, error)
	DeleteFeedItems(*item.DeleteFeedItemsCommand) error
}

type AuthService interface {
	Login(command *auth.LoginCommand) (*auth.LoginResponse, error)
	Signup(command *auth.SignupCommand) error
	Setup(command *auth.SetupCommand) error
	NeedsSetup() (bool, error)
	IsSignupAllowed() bool
	SetAllowSignup(allow bool) error
	GetSession(token string) (*auth.Session, error)
	GetAllUsers() ([]*auth.User, error)
}

// Router handles http requests
type Router struct {
	http.Handler
	feedService   FeedService
	itemService   ItemService
	authService   AuthService
	folderService FolderService
}

// New creates a new handler to handle http requests. staticFS serves the
// frontend assets at "/" via the NotFound fallback.
func New(
	itemService ItemService,
	feedService FeedService,
	authService AuthService,
	folderService FolderService,
	staticFS fs.FS,
) (*Router, error) {
	if staticFS == nil {
		return nil, errors.New("staticFS must not be nil")
	}
	h := &Router{
		itemService:   itemService,
		feedService:   feedService,
		authService:   authService,
		folderService: folderService,
	}
	router := httprouter.New()

	chain := middleware.NewChain(middleware.ContentTypeJSON, middleware.AuthMiddleware(h.authService))

	router.GET("/health", h.health)

	router.GET("/feeds", chain.Wrap(h.listFeeds))
	router.POST("/feeds", chain.Wrap(h.addFeed))
	router.POST("/feeds-bulk-move", chain.Wrap(h.bulkMoveFeeds))

	router.DELETE("/feeds/:id", chain.Wrap(h.deleteFeed))
	router.PATCH("/feeds/:id", chain.Wrap(h.patchFeed))
	router.PUT("/feeds/:id/fetch", chain.Wrap(h.fetchFeedNewItems))
	router.POST("/feeds/:id/read", chain.Wrap(h.readFeedItems))
	router.POST("/feeds/:id/unread", chain.Wrap(h.unreadFeedItems))
	router.GET("/feeds/:id/items", chain.Wrap(h.getFeedItems))

	router.GET("/folders", chain.Wrap(h.listFolders))
	router.POST("/folders", chain.Wrap(h.addFolder))
	router.PATCH("/folders/:id", chain.Wrap(h.updateFolder))
	router.DELETE("/folders/:id", chain.Wrap(h.deleteFolder))
	router.GET("/folders/:id/items", chain.Wrap(h.getFolderItems))
	router.POST("/folders/:id/read", chain.Wrap(h.readFolderItems))

	router.PUT("/items/:id", chain.Wrap(h.updateItem))
	router.GET("/items/unread", chain.Wrap(h.getUnreadItems))
	router.GET("/items/starred", chain.Wrap(h.getStarredItems))
	router.GET("/items/unread/count", chain.Wrap(h.getUnreadItemsCount))
	router.GET("/items/starred/count", chain.Wrap(h.getStarredItemsCount))

	router.POST("/login", chain.Wrap(h.login))
	router.POST("/signup", chain.Wrap(h.signup))
	router.POST("/setup", chain.Wrap(h.setup))
	router.GET("/setup/status", chain.Wrap(h.needsSetup))

	router.GET("/settings", chain.Wrap(h.getSettings))
	router.PUT("/settings", chain.Wrap(h.updateSettings))
	// Serve static frontend assets via NotFound. Restrict to GET/HEAD so a
	// stray POST/PUT/DELETE to an unknown path returns 404 instead of the
	// 405 that http.FileServer would emit, which would mask missing API
	// endpoints from clients.
	fileServer := http.FileServer(http.FS(staticFS))
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.NotFound(w, r)
			return
		}
		fileServer.ServeHTTP(w, r)
	})

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
