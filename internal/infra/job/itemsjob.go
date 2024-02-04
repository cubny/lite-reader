package job

import (
	"github.com/cubny/lite-reader/internal/app/feed"
	"github.com/cubny/lite-reader/internal/app/item"
	log "github.com/sirupsen/logrus"
)

type FeedService interface {
	ListFeeds() ([]*feed.Feed, error)
	FetchItems(feedId int) ([]*item.Item, error)
}
type ItemService interface {
	UpsertItems(command *item.UpsertItemsCommand) error
}

type ItemsJob struct {
	feedService FeedService
	itemService ItemService
}

func NewItemsJob(feedService FeedService, itemService ItemService) *ItemsJob {
	return &ItemsJob{feedService: feedService, itemService: itemService}
}

func (j *ItemsJob) Execute() {
	feeds, err := j.feedService.ListFeeds()
	if err != nil {
		return
	}
	log.Infof("Found %d feeds", len(feeds))
	for _, f := range feeds {
		items, err := j.feedService.FetchItems(f.Id)
		if err != nil {
			log.Errorf("Failed to fetch items for feed %d: %v", f.Id, err)
			continue
		}
		log.Infof("Fetched %d items for feed %d", len(items), f.Id)
		upsertItemsCommand := &item.UpsertItemsCommand{FeedId: f.Id, Items: items}
		if err := j.itemService.UpsertItems(upsertItemsCommand); err != nil {
			continue
		}
	}
}
