package feed

import (
	"errors"
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/cubny/lite-reader/internal/app/item"
)

type ServiceImpl struct {
	repository Repository
	Parser     Parser
	finder     Finder
}

func NewService(repository Repository, parser Parser, finder Finder) *ServiceImpl {
	return &ServiceImpl{
		repository: repository,
		Parser:     parser,
		finder:     finder,
	}
}

func (s *ServiceImpl) AddFeed(command *AddFeedCommand) (*Feed, error) {
	parsedFeed, err := s.Parser.ParseURL(command.URL)
	switch {
	case errors.Is(err, gofeed.ErrFeedTypeNotDetected):
		links, findErr := s.finder.FindFeeds(command.URL)
		if findErr != nil {
			return nil, fmt.Errorf("cannot find feeds: %w", findErr)
		}
		for _, link := range links {
			parsedFeed, err = s.Parser.ParseURL(link)
			if err == nil {
				break
			}
		}
	case err != nil:
		return nil, fmt.Errorf("cannot parse feed: %w", err)
	}

	if parsedFeed == nil {
		return nil, fmt.Errorf("cannot parse feed: %w", err)
	}

	feed := &Feed{
		Title:       parsedFeed.Title,
		Description: parsedFeed.Description,
		Link:        parsedFeed.Link,
		URL:         parsedFeed.FeedLink,
		Lang:        parsedFeed.Language,
		UpdatedAt:   time.Now(),
		UnreadCount: len(parsedFeed.Items),
	}

	id, err := s.repository.AddFeed(feed)
	if err != nil {
		return nil, fmt.Errorf("cannot add feed: %w", err)
	}

	feed.ID = id

	return feed, nil
}

func (s *ServiceImpl) ListFeeds() ([]*Feed, error) {
	return s.repository.ListFeeds()
}

func (s *ServiceImpl) FetchItems(feedID int) ([]*item.Item, error) {
	feed, err := s.repository.GetFeed(feedID)
	if err != nil {
		return nil, fmt.Errorf("cannot get feed: %w", err)
	}

	parsedFeed, err := s.Parser.ParseURL(feed.URL)
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

func (s *ServiceImpl) DeleteFeed(command *DeleteFeedCommand) error {
	return s.repository.DeleteFeed(command.FeedID)
}
