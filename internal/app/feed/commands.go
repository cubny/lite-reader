package feed

type AddFeedCommand struct {
	URL    string
	UserID int
}

type DeleteFeedCommand struct {
	FeedID int
}
