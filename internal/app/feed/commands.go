package feed

type AddFeedCommand struct {
	URL    string
	UserID int
}

type DeleteFeedCommand struct {
	FeedID int
}

type MoveFeedCommand struct {
	FeedID   int
	UserID   int
	FolderID *int
}

type ReorderFeedCommand struct {
	FeedID   int
	UserID   int
	Position int
}

type BulkMoveFeedsCommand struct {
	FeedIDs  []int
	UserID   int
	FolderID *int
}
