package webserver

import (
	"les-randoms/database"
	"reflect"

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

func newNavData(s *sessions.Session) navData {
	data := navData{}
	data.IsAdmin = !isNotAdmin(s)
	return data
}

type subnavData struct {
	Title                   string
	SubnavItems             []subnavItem
	SelectedSubnavItemIndex int
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

func newContentHeaderData(s *sessions.Session) contentHeaderData {
	data := contentHeaderData{}
	data.IsAuthentified = !isNotAuthentified(s)
	data.DiscordId = getDiscordId(s)
	data.Username = getUsername(s)
	data.AvatarId = getAvatarId(s)
	return data
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
