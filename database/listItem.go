package database

import (
	"fmt"
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

func (listItem ListItem) ToStringSlice() []string {
	return []string{fmt.Sprint(listItem.Id), fmt.Sprint(listItem.ListId), fmt.Sprint(listItem.OwnerId), listItem.Value, listItem.Date.Local().Format("02/01/2006 15:04:05")}
}

func ListItems_ToDBStructSlice(listItems []ListItem) []DBStruct {
	var r []DBStruct
	for _, listItem := range listItems {
		r = append(r, listItem)
	}
	return r
}

func ListItem_GetType() reflect.Type {
	return reflect.Indirect(reflect.ValueOf(&ListItem{})).Type()
}

func ListItem_SelectAll(queryPart string) ([]ListItem, error) {
	rows, err := SelectDatabase("id, listId, ownerId, value, date FROM ListItem " + queryPart)
	if err != nil {
		return nil, err
	}
	listItems := make([]ListItem, 0)
	for rows.Next() {
		var id int
		var listId int
		var ownerId int
		var value string
		var date time.Time
		err = rows.Scan(&id, &listId, &ownerId, &value, &date)
		if err != nil {
			return nil, err
		}
		listItems = append(listItems, ListItem{Id: id, ListId: listId, OwnerId: ownerId, Value: value, Date: date})
	}
	return listItems, nil
}
