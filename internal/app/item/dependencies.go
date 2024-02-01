package item

type Repository interface {
	UpsertItems(feedId int, items []*Item) error
	GetUnreadItems() ([]*Item, error)
	GetStarredItems() ([]*Item, error)
	GetFeedItems(feedId int) ([]*Item, error)
	UpdateItem(id int, starred bool, isNew bool) error
	ReadFeedItems(feedId int) error
	UnreadFeedItems(feedId int) error
	GetStarredItemsCount() (int, error)
	GetUnreadItemsCount() (int, error)
	DeleteFeedItems(feedId int) error
}
