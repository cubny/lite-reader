package item

import (
	"context"
	"errors"
	"time"
)

var ErrItemNotFound = errors.New("item not found")

type ServiceImpl struct {
	repository Repository
	scraper    Scraper
	now        func() time.Time
}

func NewService(repository Repository, scraper Scraper) *ServiceImpl {
	return &ServiceImpl{
		repository: repository,
		scraper:    scraper,
		now:        time.Now,
	}
}

func (s *ServiceImpl) GetUnreadItems() ([]*Item, error) {
	return s.repository.GetUnreadItems()
}

func (s *ServiceImpl) GetStarredItems() ([]*Item, error) {
	return s.repository.GetStarredItems()
}

func (s *ServiceImpl) GetFeedItems(command *GetFeedItemsCommand) ([]*Item, error) {
	return s.repository.GetFeedItems(command.FeedID)
}

func (s *ServiceImpl) GetFolderItems(command *GetFolderItemsCommand) ([]*Item, error) {
	return s.repository.GetFolderItems(command.FolderID, command.UserID)
}

func (s *ServiceImpl) UpsertItems(command *UpsertItemsCommand) error {
	return s.repository.UpsertItems(command.FeedID, command.Items)
}

func (s *ServiceImpl) UpdateItem(command *UpdateItemCommand) error {
	return s.repository.UpdateItem(command.ID, command.Starred, command.IsNew)
}

func (s *ServiceImpl) ReadFeedItems(command *ReadFeedItemsCommand) error {
	return s.repository.ReadFeedItems(command.FeedID)
}

func (s *ServiceImpl) UnreadFeedItems(command *UnreadFeedItemsCommand) error {
	return s.repository.UnreadFeedItems(command.FeedID)
}

func (s *ServiceImpl) ReadFolderItems(command *ReadFolderItemsCommand) error {
	return s.repository.ReadFolderItems(command.FolderID, command.UserID)
}

func (s *ServiceImpl) GetStarredItemsCount() (int, error) {
	return s.repository.GetStarredItemsCount()
}

func (s *ServiceImpl) GetUnreadItemsCount() (int, error) {
	return s.repository.GetUnreadItemsCount()
}

func (s *ServiceImpl) DeleteFeedItems(command *DeleteFeedItemsCommand) error {
	return s.repository.DeleteFeedItems(command.FeedID)
}

// ScrapeItem fetches the full article body for the given item, persists it,
// and returns the refreshed item. Errors from the scraper are recorded as
// scrape_status='error' on the row before being returned so the UI can
// distinguish "never tried" from "tried and failed".
func (s *ServiceImpl) ScrapeItem(ctx context.Context, command *ScrapeItemCommand) (*Item, error) {
	existing, err := s.repository.GetItemForUser(command.ID, command.UserID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrItemNotFound
	}

	now := s.now()
	content, scrapeErr := s.scraper.Scrape(ctx, existing.Link)
	status := "ok"
	if scrapeErr != nil {
		status = "error"
		content = ""
	}
	if err := s.repository.UpdateItemContent(existing.ID, content, status, now); err != nil {
		return nil, err
	}
	if scrapeErr != nil {
		return nil, scrapeErr
	}

	existing.FullContent = content
	existing.ScrapeStatus = status
	existing.ScrapedAt = &now
	return existing, nil
}
