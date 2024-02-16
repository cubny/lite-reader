package feed

type AddFeedCommand struct {
	URL string
}

type DeleteFeedCommand struct {
	FeedID int
}
