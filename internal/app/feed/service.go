package feed

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	"time"
)

type ServiceImpl struct {
}

func (s ServiceImpl) AddFeed(command *AddFeedCommand) (*Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(command.URL)
	if err != nil {
		return nil, fmt.Errorf("cannot parse feed: %w", err)
	}
	return &Feed{
		Id:          uuid.NewString(),
		Title:       feed.Title,
		Description: feed.Description,
		Link:        feed.Link,
		URL:         feed.FeedLink,
		Updated:     feed.Updated,
		Lang:        feed.Language,
		UpdatedAt:   time.Now(),
	}, nil
}

func NewService() (*ServiceImpl, error) {
	return &ServiceImpl{}, nil
}
