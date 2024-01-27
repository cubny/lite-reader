package item

import "time"

type Item struct {
	Id        string
	Title     string
	Desc      string
	Link      string
	IsNew     bool
	Starred   bool
	Timestamp time.Time
}
