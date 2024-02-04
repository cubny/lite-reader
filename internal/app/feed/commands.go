package feed

type AddFeedCommand struct {
	URL string
}

type DeleteFeedCommand struct {
	FeedId int
}
