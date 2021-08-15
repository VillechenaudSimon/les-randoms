package database

import (
	"reflect"
	"time"
)

type ListItem struct {
	Id      int
	ListId  int
	OwnerId int
	Value   string
	Date    time.Time
}

func ListItem_GetType() reflect.Type {
	return reflect.Indirect(reflect.ValueOf(&ListItem{})).Type()
}
