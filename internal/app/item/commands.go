package item

type GetFeedItemsCommand struct {
	FeedID int
}

type UpsertItemsCommand struct {
	FeedID int
	Items  []*Item
}

type UpdateItemCommand struct {
	ID      int
	Starred bool
	IsNew   bool
}

type FetchFeedNewItemsCommand struct {
	FeedID int
}

type ReadFeedItemsCommand struct {
	FeedID int
}

type UnreadFeedItemsCommand struct {
	FeedID int
}

type DeleteFeedItemsCommand struct {
	FeedID int
}
