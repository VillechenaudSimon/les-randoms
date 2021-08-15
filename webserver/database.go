package webserver

import (
	"les-randoms/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleDatabaseRoute(c *gin.Context) {
	data := &databaseData{}
	data.LayoutData.SubnavData.Title = "Database"

	data.LayoutData.SubnavData.SubnavItems = append(data.LayoutData.SubnavData.SubnavItems, subnavItem{Name: "Users"})
	data.LayoutData.SubnavData.SubnavItems = append(data.LayoutData.SubnavData.SubnavItems, subnavItem{Name: "Lists"})
	data.LayoutData.SubnavData.SubnavItems = append(data.LayoutData.SubnavData.SubnavItems, subnavItem{Name: "ListItems"})

	data.LayoutData.SubnavData.SelectedSubnavItemIndex = 0
	if c.Request.Method == "POST" {
		selectedItemName := c.PostForm("subnavSelectedItem")
		for i := 0; i < len(data.LayoutData.SubnavData.SubnavItems); i++ {
			if selectedItemName == data.LayoutData.SubnavData.SubnavItems[i].Name {
				data.LayoutData.SubnavData.SelectedSubnavItemIndex = i
				break
			}
		}
	}

	switch data.LayoutData.SubnavData.SelectedSubnavItemIndex {
	case 0:
		data.EntityTableData = newCustomTableDataFromDBStruct(database.User_GetType())
	case 1:
		data.EntityTableData = newCustomTableDataFromDBStruct(database.ListItem_GetType())
	case 2:
		data.EntityTableData = newCustomTableDataFromDBStruct(database.List_GetType())
	}

	c.HTML(http.StatusOK, "database.tmpl.html", data)
}
