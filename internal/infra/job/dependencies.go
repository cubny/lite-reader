package job

import (
	"github.com/cubny/lite-reader/internal/app/auth"
	"github.com/cubny/lite-reader/internal/app/feed"
	"github.com/cubny/lite-reader/internal/app/item"
)

type FeedService interface {
	ListFeeds(userID int64) ([]*feed.Feed, error)
	FetchItems(feedID int) ([]*item.Item, error)
}
type ItemService interface {
	UpsertItems(command *item.UpsertItemsCommand) error
}

type UserService interface {
	GetAllUsers() ([]*auth.User, error)
}
