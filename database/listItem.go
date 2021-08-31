package database

import (
	"fmt"
	"les-randoms/utils"
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
	return []string{fmt.Sprint(listItem.Id), fmt.Sprint(listItem.ListId), fmt.Sprint(listItem.OwnerId), listItem.Value, listItem.Date.Format(utils.DBDateTimeFormat)}
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
		utils.LogError("Error while selecting on ListItem table : " + err.Error())
		return nil, err
	}
	listItems := make([]ListItem, 0)
	for rows.Next() {
		var id int
		var listId int
		var ownerId int
		var value string
		var date []uint8
		err = rows.Scan(&id, &listId, &ownerId, &value, &date)
		if err != nil {
			utils.LogError("Error while scanning a ListItem : " + err.Error())
			continue
		}
		parsedDate, err := time.Parse(utils.DBDateTimeFormat, string(date))
		if err != nil {
			utils.LogError("Error while parsing a listItem date : " + err.Error())
			continue
		}
		listItems = append(listItems, ListItem{Id: id, ListId: listId, OwnerId: ownerId, Value: value, Date: parsedDate})
	}
	return listItems, nil
}
