package item

import "time"

type Item struct {
	Id        int
	Title     string
	Desc      string
	Dir       string
	Link      string
	IsNew     bool
	Starred   bool
	Timestamp time.Time
}
