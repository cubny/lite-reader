package item

type GetFeedItemsCommand struct {
	FeedID int
}

type GetFolderItemsCommand struct {
	FolderID int
	UserID   int
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

type ReadFolderItemsCommand struct {
	FolderID int
	UserID   int
}

type UnreadFeedItemsCommand struct {
	FeedID int
}

type DeleteFeedItemsCommand struct {
	FeedID int
}
