package feed

type AddFeedCommand struct {
	URL string
}

type ListFeedsCommand struct {
}
type DeleteFeedCommand struct {
	FeedId int
}
