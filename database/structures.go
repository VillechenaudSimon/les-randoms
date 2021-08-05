package database

import "time"

type User struct {
	Id       int
	Name     string
	Password string
}

type List struct {
	Id   int
	Name string
}

type ListItem struct {
	Id      int
	ListId  int
	OwnerId int
	Value   string
	Date    time.Time
}
