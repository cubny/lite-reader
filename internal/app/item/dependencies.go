package item

type Repository interface {
	UpsertItems(feedId int, items []*Item) error
	GetUnreadItems() ([]*Item, error)
	GetStarredItems() ([]*Item, error)
	GetFeedItems(feedId int) ([]*Item, error)
}
