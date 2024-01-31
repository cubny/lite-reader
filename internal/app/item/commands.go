package item

type GetUnreadItemsCommand struct {
}

type GetStarredItemsCommand struct {
}

type GetFeedItemsCommand struct {
	FeedId int
}

type UpsertItemsCommand struct {
	FeedId int
	Items  []*Item
}

type UpdateItemCommand struct {
	Id      int
	Starred bool
	IsNew   bool
}

type FetchFeedNewItemsCommand struct {
	FeedId int
}

type ReadFeedItemsCommand struct {
	FeedId int
}

type UnreadFeedItemsCommand struct {
	FeedId int
}
