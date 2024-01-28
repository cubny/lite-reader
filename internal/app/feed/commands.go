package feed

type AddFeedCommand struct {
	URL string
}

type GetFeedCommand struct {
	ID string
}

type ListFeedCommand struct {
}
