package database

import "reflect"

type List struct {
	Id   int
	Name string
}

func List_GetType() reflect.Type {
	return reflect.Indirect(reflect.ValueOf(&List{})).Type()
}
