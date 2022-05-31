package admin

import (
	"io"
	"les-randoms/database"
	"les-randoms/utils"
	webserver "les-randoms/webserver/logic"
	"net/http"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() {
	database := webserver.Router.Group("/database")
	database.GET("", handleDatabaseRoute)
	database.POST("", handleDatabaseRoute)
	database.GET("/:subNavItem", handleDatabaseRoute)
	database.POST("/:subNavItem", handleDatabaseRoute)
}

func handleDatabaseRoute(c *gin.Context) {
	session := webserver.GetSession(c)

	if c.Param("subNavItem") == "" {
		c.Redirect(http.StatusFound, "/database/Users")
	}

	if webserver.IsNotAuthentified(session) {
		webserver.RedirectToAuth(c)
		return
	}

	if webserver.GetAccessStatus(session, "/database") <= database.RightTypes.Forbidden {
		webserver.RedirectToIndex(c)
		return
	}

	if handleDBFileDownload(c) {
		return
	}

	handleDBFileUpload(c)

	data := databaseData{}

	webserver.SetupNavData(&data.LayoutData.NavData, session)

	selectedItemName := webserver.SetupSubnavData(&data.LayoutData.SubnavData, c, "Database", []string{"Users", "Lists", "ListItems", "AccessRights", "Summoners"}, map[string]string{"Users": "Users", "Lists": "Lists", "ListItems": "List Items", "AccessRights": "Access Rights", "Summoners": "Summoners"})

	webserver.SetupContentHeaderData(&data.ContentHeaderData, session)
	data.ContentHeaderData.Title = selectedItemName

	data.SelectParameters.SelectQueryBody = c.PostForm("db-select-query-body-textbar")
	err := setupDatabaseEntityTableData(&data)
	if err != nil {
		c.Redirect(http.StatusFound, "/database/Users")
	}

	c.HTML(http.StatusFound, "database.tmpl.html", data)
}

type databaseData struct {
	LayoutData        webserver.LayoutData
	ContentHeaderData webserver.ContentHeaderData
	EntityTableData   webserver.CustomTableData
	SelectParameters  struct {
		SelectQueryBody string
	}
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
	case 4:
		summoners, err := database.Summoner_SelectAll(data.SelectParameters.SelectQueryBody)
		if err == nil {
			data.EntityTableData = newCustomTableDataFromDBStruct(database.Summoner_GetType(), database.Summoners_ToDBStructSlice(summoners))
		}
	}
	return nil
}

func handleDBFileDownload(c *gin.Context) bool {
	if c.Request.URL.Path == "/database/sqlite-database.db" {
		utils.LogInfo("Requesting to download database file")
		c.File("sqlite-database.db")
		return true
	}
	return false
}

func handleDBFileUpload(c *gin.Context) {
	err := c.Request.ParseMultipartForm(1024 * 1024 * 256)
	if err != nil {
		//utils.LogError(err.Error())
		return
	}
	file, fileHeader, err := c.Request.FormFile("upload-db-file")
	if err != nil {
		//utils.LogError(err.Error())
		return
	}

	if fileHeader.Filename != "sqlite-database.db" {
		return
	}

	database.CloseDatabase()

	os.Remove("sqlite-database.db")
	dst, err := os.Create("sqlite-database.db")
	if err != nil {
		utils.HandlePanicError(err)
	}
	_, err = io.Copy(dst, file)
	if err != nil {
		utils.HandlePanicError(err)
	}

	database.OpenDatabase()
	database.VerifyDatabase()

	defer file.Close()
}

func newCustomTableDataFromDBStruct(structType reflect.Type, dbStructs []database.DBStruct) webserver.CustomTableData {
	data := webserver.CustomTableData{}
	for i := 0; i < structType.NumField(); i++ {
		data.HeaderList = append(data.HeaderList, structType.Field(i).Name)
	}

	for _, dbStruct := range dbStructs {
		data.ItemList = append(data.ItemList, webserver.TableItemData{FieldList: dbStruct.ToStringSlice()})
	}
	data.SortColumnIndex = -1

	data.ColumnTypes = make([]webserver.CustomTableColumnType, len(data.HeaderList))

	for i := 0; i < len(data.ColumnTypes); i++ {
		data.ColumnTypes[i] = webserver.CustomTableColumnTypeText
	}

	return data
}
