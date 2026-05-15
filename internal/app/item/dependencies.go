package item

import (
	"context"
	"time"
)

type Repository interface {
	UpsertItems(feedID int, items []*Item) error
	GetUnreadItems() ([]*Item, error)
	GetStarredItems() ([]*Item, error)
	GetFeedItems(feedID int) ([]*Item, error)
	GetFolderItems(folderID, userID int) ([]*Item, error)
	UpdateItem(id int, starred bool, isNew bool) error
	ReadFeedItems(feedID int) error
	UnreadFeedItems(feedID int) error
	ReadFolderItems(folderID, userID int) error
	GetStarredItemsCount() (int, error)
	GetUnreadItemsCount() (int, error)
	DeleteFeedItems(feedID int) error
	GetItemForUser(id, userID int) (*Item, error)
	UpdateItemContent(id int, fullContent, status string, scrapedAt time.Time) error
}

type Scraper interface {
	Scrape(ctx context.Context, url string) (string, error)
}
