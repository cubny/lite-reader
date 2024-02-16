package item

type Repository interface {
	UpsertItems(feedID int, items []*Item) error
	GetUnreadItems() ([]*Item, error)
	GetStarredItems() ([]*Item, error)
	GetFeedItems(feedID int) ([]*Item, error)
	UpdateItem(id int, starred bool, isNew bool) error
	ReadFeedItems(feedID int) error
	UnreadFeedItems(feedID int) error
	GetStarredItemsCount() (int, error)
	GetUnreadItemsCount() (int, error)
	DeleteFeedItems(feedID int) error
}
