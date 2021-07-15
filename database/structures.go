package database

import "time"

type ListItem struct {
	Id     int
	ListId int
	Name   string
	Date   time.Time
}
