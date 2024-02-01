package feed

import "time"

type Feed struct {
	Id          int
	Title       string
	Description string
	Link        string
	URL         string
	Lang        string
	UpdatedAt   time.Time
	UnreadCount int
}
