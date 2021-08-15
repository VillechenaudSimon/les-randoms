package webserver

import (
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

	c.HTML(http.StatusOK, "database.tmpl.html", data)
}
