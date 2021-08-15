package database

import "time"

type ListItem struct {
	Id      int
	ListId  int
	OwnerId int
	Value   string
	Date    time.Time
}
