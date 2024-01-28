package feed

import (
	"errors"
	"fmt"
	"github.com/cubny/lite-reader/internal/app/item"
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
	case errors.Is(err, gofeed.ErrFeedTypeNotDetected):
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

func (s *ServiceImpl) ListFeeds(command *ListFeedCommand) ([]*Feed, error) {
	return s.repository.ListFeeds()
}

func (s *ServiceImpl) FetchItems(feedId int) ([]*item.Item, error) {
	feed, err := s.repository.GetFeed(feedId)
	if err != nil {
		return nil, fmt.Errorf("cannot get feed: %w", err)
	}

	fp := gofeed.NewParser()
	parsedFeed, err := fp.ParseURL(feed.URL)
	if err != nil {
		return nil, fmt.Errorf("cannot parse feed: %w", err)
	}

	items := make([]*item.Item, 0)
	for _, t := range parsedFeed.Items {
		timestamp := t.PublishedParsed
		items = append(items, &item.Item{
			Title:     t.Title,
			Desc:      t.Content,
			Link:      t.Link,
			Timestamp: *timestamp,
			Dir:       t.Description,
			IsNew:     true,
			Starred:   false,
		})
	}

	return items, nil
}
