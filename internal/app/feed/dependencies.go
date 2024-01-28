package feed

type Repository interface {
	AddFeed(feed *Feed) (int, error)
	GetFeed(id string) (*Feed, error)
}
