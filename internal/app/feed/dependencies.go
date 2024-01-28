package feed

type Repository interface {
	AddFeed(feed *Feed) (int, error)
	GetFeed(id int) (*Feed, error)
	ListFeeds() ([]*Feed, error)
}
