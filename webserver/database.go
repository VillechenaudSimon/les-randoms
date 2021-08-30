package webserver

import (
	"les-randoms/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleDatabaseRoute(c *gin.Context) {
	session := getSession(c)

	if isNotAuthentified(session) {
		redirectToAuth(c)
	}

	if isNotAdmin(session) {
		redirectToIndex(c)
	}

	data := &databaseData{}

	data.LayoutData.NavData = newNavData(session)

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
		users, err := database.User_SelectAll()
		if err == nil {
			data.EntityTableData = newCustomTableDataFromDBStruct(database.User_GetType(), database.Users_ToDBStructSlice(users))
		}
	case 1:
		lists, err := database.List_SelectAll()
		if err == nil {
			data.EntityTableData = newCustomTableDataFromDBStruct(database.List_GetType(), database.Lists_ToDBStructSlice(lists))
		}
	case 2:
		listItems, err := database.ListItem_SelectAll()
		if err == nil {
			data.EntityTableData = newCustomTableDataFromDBStruct(database.ListItem_GetType(), database.ListItems_ToDBStructSlice(listItems))
		}
	}

	c.HTML(http.StatusOK, "database.tmpl.html", data)
}
