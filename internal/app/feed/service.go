package feed

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/nikhil1raghav/feedfinder"
)

type ServiceImpl struct {
	repository Repository
}

func NewService(repository Repository) (*ServiceImpl, error) {
	return &ServiceImpl{repository: repository}, nil
}

func (s *ServiceImpl) AddFeed(command *AddFeedCommand) (*Feed, error) {
	fp := gofeed.NewParser()
	parsedFeed, err := fp.ParseURL(command.URL)
	switch {
	case err == gofeed.ErrFeedTypeNotDetected:
		f := feedfinder.NewFeedFinder()
		links, _ := f.FindFeeds(command.URL)
		for _, link := range links {
			parsedFeed, err = fp.ParseURL(link)
			if err == nil {
				break
			}
		}
	case err != nil:
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

func (s *ServiceImpl) ListFeeds(command *ListFeedsCommand) ([]*Feed, error) {
	return s.repository.ListFeeds()
}
