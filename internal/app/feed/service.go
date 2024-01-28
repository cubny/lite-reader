package feed

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

type ServiceImpl struct {
	repository Repository
}

func NewService(repository Repository) (*ServiceImpl, error) {
	return &ServiceImpl{repository: repository}, nil
}

func (s ServiceImpl) AddFeed(command *AddFeedCommand) (*Feed, error) {
	fp := gofeed.NewParser()
	parsedFeed, err := fp.ParseURL(command.URL)
	if err != nil {
		return nil, fmt.Errorf("cannot parse feed: %w", err)
	}

	feed := &Feed{
		Title:       parsedFeed.Title,
		Description: parsedFeed.Description,
		Link:        parsedFeed.Link,
		URL:         parsedFeed.FeedLink,
		Lang:        parsedFeed.Language,
		UpdatedAt:   time.Now(),
	}

	id, err := s.repository.AddFeed(feed)
	if err != nil {
		return nil, fmt.Errorf("cannot add feed: %w", err)
	}

	feed.Id = id

	return feed, nil
}
