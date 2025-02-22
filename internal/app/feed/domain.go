package feed

import "time"

type Feed struct {
	ID          int
	Title       string
	Description string
	Link        string
	URL         string
	Lang        string
	UpdatedAt   time.Time
	UnreadCount int
	UserID      int
}
