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
