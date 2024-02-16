package item

import "time"

type Item struct {
	ID        int
	Title     string
	Desc      string
	Dir       string
	Link      string
	IsNew     bool
	Starred   bool
	Timestamp time.Time
}
