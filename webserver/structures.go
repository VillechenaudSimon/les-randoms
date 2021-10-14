package webserver

import (
	"les-randoms/database"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type indexData struct {
	LayoutData        layoutData
	ContentHeaderData contentHeaderData
}

type aramData struct {
	LayoutData        layoutData
	ContentHeaderData contentHeaderData
	ListTableData     customTableData
}

type playersData struct {
	LayoutData        layoutData
	ContentHeaderData contentHeaderData
}

type databaseData struct {
	LayoutData        layoutData
	ContentHeaderData contentHeaderData
	EntityTableData   customTableData
}

type loginData struct {
}

type layoutData struct {
	NavData    navData
	SubnavData subnavData
}

type navData struct {
	IsAdmin bool
}

func setupNavData(data *navData, s *sessions.Session) {
	data.IsAdmin = !isNotAdmin(s)
}

type subnavData struct {
	Title                   string
	SubnavItems             []subnavItem
	SelectedSubnavItemIndex int
}

// Returns the selectedItemName
func setupSubnavData(data *subnavData, c *gin.Context, title string, subnavItemsName []string) string {
	data.Title = title
	data.SelectedSubnavItemIndex = 0

	for _, name := range subnavItemsName {
		data.SubnavItems = append(data.SubnavItems, subnavItem{Name: name})
	}

	selectedItemName := data.SubnavItems[data.SelectedSubnavItemIndex].Name
	if c.Request.Method == "POST" {
		selectedItemName = c.PostForm("subnavSelectedItem")
	}

	for i := 0; i < len(data.SubnavItems); i++ {
		if selectedItemName == data.SubnavItems[i].Name {
			data.SelectedSubnavItemIndex = i
			break
		}
	}
	return selectedItemName
}

type subnavItem struct {
	Name string
}

type contentHeaderData struct {
	Title          string
	IsAuthentified bool
	DiscordId      string
	Username       string
	AvatarId       string
}

func setupContentHeaderData(data *contentHeaderData, s *sessions.Session) {
	data.IsAuthentified = !isNotAuthentified(s)
	data.DiscordId = getDiscordId(s)
	data.Username = getUsername(s)
	data.AvatarId = getAvatarId(s)
	data.Title = "Default"
}

type customTableData struct {
	HeaderList []string
	ItemList   []tableItemData
}

func newCustomTableDataFromDBStruct(structType reflect.Type, dbStructs []database.DBStruct) customTableData {
	data := customTableData{}
	for i := 0; i < structType.NumField(); i++ {
		data.HeaderList = append(data.HeaderList, structType.Field(i).Name)
	}

	for _, dbStruct := range dbStructs {
		data.ItemList = append(data.ItemList, tableItemData{FieldList: dbStruct.ToStringSlice()})
	}

	return data
}

type tableItemData struct {
	FieldList []string
}

/*
Unused in nav.tmpl.html :
{{ range $i, $navItem := .NavItems }}
{{ if $navItem.IsGroup }}
<div class="nav-group expanded">
  <div class="nav-group-header">
    <img src="{{ $navItem.ImgSrc }}">
  </div>
  <div class="nav-group-content">
    {{ range $j, $navGroupItem := $navItem.NavGroupItems }}
    <a class="nav-item" href="{{ $navItem.Href }}">
      <img src="{{ $navItem.ImgSrc }}" />
    </a>
    {{ end }}
  </div>
</div>
{{ else }}
<a class="nav-item" href="{{ $navItem.Href }}">
  <img src="{{ $navItem.ImgSrc }}" />
</a>
{{ end }}
{{ end }}

Unused in structures.go :
type navData struct {
	NavItems []navItem
}
type navItem struct {
	IsGroup       bool
	Href          string
	ImgSrc        string
	NavGroupItems []navItem
}

*/
