package feed

import "github.com/mmcdole/gofeed"

type Repository interface {
	AddFeed(feed *Feed) (int, error)
	GetFeed(id int) (*Feed, error)
	ListFeeds(userID int) ([]*Feed, error)
	DeleteFeed(id int) error
	MoveFeed(feedID, userID int, folderID *int) error
	ReorderFeed(feedID, userID int, position int) error
	BulkMoveFeeds(feedIDs []int, userID int, folderID *int) error
}

type Parser interface {
	ParseURL(url string) (*gofeed.Feed, error)
}

type Finder interface {
	FindFeeds(url string) ([]string, error)
}
