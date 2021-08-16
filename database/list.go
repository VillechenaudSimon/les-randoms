package database

import (
	"fmt"
	"reflect"
)

type List struct {
	Id          int
	Name        string
	ColumnCount int
}

func List_GetType() reflect.Type {
	return reflect.Indirect(reflect.ValueOf(&List{})).Type()
}

func List_ToStringSlice(list List) []string {
	return []string{fmt.Sprint(list.Id), list.Name, fmt.Sprint(list.ColumnCount)}
}

func List_SelectAll() ([]List, error) {
	rows, err := SelectDatabase("id, name, ColumnCount FROM List")
	if err != nil {
		return nil, err
	}
	lists := make([]List, 0)
	for rows.Next() {
		var id int
		var name string
		var columnCount int
		err = rows.Scan(&id, &name, &columnCount)
		if err != nil {
			return nil, err
		}
		lists = append(lists, List{Id: id, Name: name, ColumnCount: columnCount})
	}
	return lists, nil
}
