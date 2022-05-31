package lol

import (
	"fmt"
	"les-randoms/database"
	"les-randoms/utils"
	webserver "les-randoms/webserver/logic"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleAramRoute(c *gin.Context) {
	session := webserver.GetSession(c)

	if webserver.IsNotAuthentified(session) {
		webserver.RedirectToAuth(c)
		return
	}

	if webserver.GetAccessStatus(session, "/lol/aram") <= database.RightTypes.Forbidden {
		webserver.RedirectToIndex(c)
		return
	}

	data := aramData{}

	webserver.SetupNavData(&data.LayoutData.NavData, session)

	selectedItemName := webserver.SetupSubnavData(&data.LayoutData.SubnavData, c, "Aram Gaming", []string{"GoldenList", "BlackList", "BotList" /*, "TierList"*/}, map[string]string{"GoldenList": "Golden List", "BlackList": "Black List", "BotList": "Bot List" /*, "TierList": "Tier List"*/})

	webserver.SetupContentHeaderData(&data.ContentHeaderData, session)
	data.ContentHeaderData.Title = selectedItemName

	data.ListTableData = webserver.CustomTableData{}
	list, err := database.List_SelectFirst("WHERE name = " + utils.Esc(selectedItemName))
	if err != nil {
		webserver.RedirectToIndex(c)
		utils.LogError("Error while creating customTableData item with a DB List (Selected list : " + selectedItemName + ") : " + err.Error())
		return
	}
	data.ListTableData.HeaderList = list.Headers
	data.ListTableData.ColumnTypes = make([]webserver.CustomTableColumnType, len(data.ListTableData.HeaderList))
	for i := 0; i < len(data.ListTableData.ColumnTypes); i++ {
		data.ListTableData.ColumnTypes[i] = webserver.CustomTableColumnTypeText
	}
	data.ListTableData.ColumnTypes[0] = webserver.CustomTableColumnTypeDate

	listItems, err := database.ListItem_SelectAll("WHERE listId = " + fmt.Sprint(list.Id) + " ORDER BY date")
	if err == nil {
		data.ListTableData.ItemList = make([]webserver.TableItemData, 0)
		for _, listItem := range listItems {
			data.ListTableData.ItemList = append(data.ListTableData.ItemList, webserver.TableItemData{FieldList: append([]string{listItem.Date.Local().Format(utils.DateFormat)}, utils.ParseDatabaseStringList(listItem.Value)...)})
		}
	}
	data.ListTableData.SortColumnIndex = -1

	c.HTML(http.StatusFound, "aram.tmpl.html", data)
}

type aramData struct {
	LayoutData        webserver.LayoutData
	ContentHeaderData webserver.ContentHeaderData
	ListTableData     webserver.CustomTableData
}
