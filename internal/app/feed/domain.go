package feed

import "time"

type Feed struct {
	Id          string
	Title       string
	Description string
	Link        string
	URL         string
	Updated     string
	Lang        string
	UpdatedAt   time.Time
}
