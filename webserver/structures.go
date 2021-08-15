package webserver

import "reflect"

type indexData struct {
	LayoutData layoutData
}

type aramData struct {
	LayoutData layoutData
}

type playersData struct {
	LayoutData layoutData
}

type databaseData struct {
	LayoutData      layoutData
	EntityTableData customTableData
}

type layoutData struct {
	NavData    navData
	SubnavData subnavData
}

type navData struct {
}

type subnavData struct {
	Title                   string
	SubnavItems             []subnavItem
	SelectedSubnavItemIndex int
}

type subnavItem struct {
	Name string
}

type customTableData struct {
	HeaderList []string
	ItemList   []string
}

func newCustomTableDataFromDBStruct(structType reflect.Type) customTableData {
	data := customTableData{}
	for i := 0; i < structType.NumField(); i++ {
		data.HeaderList = append(data.HeaderList, structType.Field(i).Name)
	}
	return data
}

type tableItemData struct {
	FieldList []string
}
