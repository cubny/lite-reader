package job

import (
	log "github.com/sirupsen/logrus"

	"github.com/cubny/lite-reader/internal/app/item"
)

type ItemsJob struct {
	feedService FeedService
	itemService ItemService
	userService UserService
}

func NewItemsJob(feedService FeedService, itemService ItemService, userService UserService) *ItemsJob {
	return &ItemsJob{feedService: feedService, itemService: itemService, userService: userService}
}

func (j *ItemsJob) Execute() {
	users, err := j.userService.GetAllUsers()
	if err != nil {
		log.Errorf("Failed to get users: %v", err)
		return
	}

	for _, u := range users {
		log.Printf("Processing user %d", u.ID)
		feeds, err := j.feedService.ListFeeds(int64(u.ID))
		if err != nil {
			return
		}
		log.Infof("Found %d feeds", len(feeds))
		for _, f := range feeds {
			items, err := j.feedService.FetchItems(f.ID)
			if err != nil {
				log.Errorf("Failed to fetch items for feed %d: %v", f.ID, err)
				continue
			}
			upsertItemsCommand := &item.UpsertItemsCommand{FeedID: f.ID, Items: items}
			log.Infof("Upserting %d items for feed %d", len(items), f.ID)
			if err := j.itemService.UpsertItems(upsertItemsCommand); err != nil {
				log.Errorf("Failed to upsert items for feed %d: %v", f.ID, err)
				continue
			}
		}
	}
}
