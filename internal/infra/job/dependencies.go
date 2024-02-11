package job

import (
	"github.com/cubny/lite-reader/internal/app/feed"
	"github.com/cubny/lite-reader/internal/app/item"
)

type FeedService interface {
	ListFeeds() ([]*feed.Feed, error)
	FetchItems(feedId int) ([]*item.Item, error)
}
type ItemService interface {
	UpsertItems(command *item.UpsertItemsCommand) error
}
