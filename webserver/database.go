package webserver

import (
	"les-randoms/database"
	"les-randoms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleDatabaseRoute(c *gin.Context) {
	session := getSession(c)

	if c.Param("subNavItem") == "" {
		c.Redirect(http.StatusFound, "/database/Users")
	}

	if isNotAuthentified(session) {
		redirectToAuth(c)
		return
	}

	if getAccessStatus(session, "/database") <= database.RightTypes.Forbidden {
		redirectToIndex(c)
		return
	}

	if c.Request.URL.Path == "/database/sqlite-database.db" {
		utils.LogInfo("Requesting to download database file")
		c.File("sqlite-database.db")
		return
	}

	data := databaseData{}

	setupNavData(&data.LayoutData.NavData, session)

	selectedItemName := setupSubnavData(&data.LayoutData.SubnavData, c, "Database", []string{"Users", "Lists", "ListItems", "AccessRights"}, map[string]string{"Users": "Users", "Lists": "Lists", "ListItems": "List Items", "AccessRights": "Access Rights"})

	setupContentHeaderData(&data.ContentHeaderData, session)
	data.ContentHeaderData.Title = selectedItemName

	data.SelectParameters.SelectQueryBody = c.PostForm("db-select-query-body-textbar")
	err := setupDatabaseEntityTableData(&data)
	if err != nil {
		c.Redirect(http.StatusFound, "/database/Users")
	}

	c.HTML(http.StatusOK, "database.tmpl.html", data)
}

func setupDatabaseEntityTableData(data *databaseData) error {
	switch data.LayoutData.SubnavData.SelectedSubnavItemIndex {
	case 0:
		users, err := database.User_SelectAll(data.SelectParameters.SelectQueryBody)
		if err == nil {
			data.EntityTableData = newCustomTableDataFromDBStruct(database.User_GetType(), database.Users_ToDBStructSlice(users))
		}
	case 1:
		lists, err := database.List_SelectAll(data.SelectParameters.SelectQueryBody)
		if err == nil {
			data.EntityTableData = newCustomTableDataFromDBStruct(database.List_GetType(), database.Lists_ToDBStructSlice(lists))
		}
	case 2:
		listItems, err := database.ListItem_SelectAll(data.SelectParameters.SelectQueryBody)
		if err == nil {
			data.EntityTableData = newCustomTableDataFromDBStruct(database.ListItem_GetType(), database.ListItems_ToDBStructSlice(listItems))
		}
	case 3:
		accessRights, err := database.AccessRight_SelectAll(data.SelectParameters.SelectQueryBody)
		if err == nil {
			data.EntityTableData = newCustomTableDataFromDBStruct(database.AccessRight_GetType(), database.AccessRights_ToDBStructSlice(accessRights))
		}
	}
	return nil
}
