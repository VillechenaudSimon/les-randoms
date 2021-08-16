package database

import (
	"database/sql"
	"fmt"
	"reflect"
)

type List struct {
	Id          int
	Name        string
	ColumnCount int
}

func (list List) ToStringSlice() []string {
	return []string{fmt.Sprint(list.Id), list.Name, fmt.Sprint(list.ColumnCount)}
}

func Lists_ToDBStructSlice(lists []List) []DBStruct {
	var r []DBStruct
	for _, list := range lists {
		r = append(r, list)
	}
	return r
}

func List_GetType() reflect.Type {
	return reflect.Indirect(reflect.ValueOf(&List{})).Type()
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

func List_CreateNew(name string, columnCount int) (List, sql.Result, error) {
	result, err := InsertDatabase("List VALUES(" + name + ", " + fmt.Sprint(columnCount) + ")")
	if err != nil {
		return List{}, result, err
	}
	newId, err := result.LastInsertId()
	if err != nil {
		return List{}, result, err
	}
	return List{Id: int(newId), Name: name, ColumnCount: columnCount}, result, err
}
