package feed

type AddFeedCommand struct {
	URL    string
	UserID int64
}

type DeleteFeedCommand struct {
	FeedID int
}
