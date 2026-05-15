package item

import "time"

// Scrape status values persisted on item.scrape_status. The empty string
// means "never tried"; ScrapeStatusOK means the article was extracted and
// stored in full_content; ScrapeStatusError means the last attempt failed.
const (
	ScrapeStatusOK    = "ok"
	ScrapeStatusError = "error"
)

type Item struct {
	ID           int
	Title        string
	Desc         string
	Dir          string
	Link         string
	IsNew        bool
	Starred      bool
	Timestamp    time.Time
	FullContent  string
	ScrapedAt    *time.Time
	ScrapeStatus string
}
