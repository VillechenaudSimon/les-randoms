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
	LayoutData         layoutData
	ContentHeaderData  contentHeaderData
	LolGameReviewData  lolGameReviewData
	LastGameParameters struct {
		SummonerName string
	}
	LadderChampPoolTableData customTableData
	LadderTableData          customTableData
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
	NavItems []navItem
}

type navItem struct {
	IsGroup       bool
	Href          string
	ImgSrc        string
	NavGroupItems []navItem
}

func setupNavData(data *navData, s *sessions.Session) {
	data.NavItems = append(data.NavItems, navItem{IsGroup: false, Href: "/", ImgSrc: "/static/images/favicon.png"})

	lolNavItems := make([]navItem, 0)
	if getAccessStatus(s, "/aram") {
		lolNavItems = append(lolNavItems, navItem{IsGroup: false, Href: "/aram", ImgSrc: "/static/images/HowlingAbyssIcon.png"})
	}
	lolNavItems = append(lolNavItems, navItem{IsGroup: false, Href: "/players", ImgSrc: "/static/images/MPengu.png"})
	data.NavItems = append(data.NavItems, navItem{IsGroup: true, ImgSrc: "/static/images/lol.ico", NavGroupItems: lolNavItems})

	if getAccessStatus(s, "/database") {
		data.NavItems = append(data.NavItems, navItem{IsGroup: false, Href: "/database", ImgSrc: "/static/images/databaseConfig.png"})
	}
}

type subnavData struct {
	Title                   string
	SubnavItems             []subnavItem
	SelectedSubnavItemIndex int
}

// Returns the selectedItemName
// subnavItemsMap map the displayable name of the item from the 'not-displayable' name ("GoldenList" -> "Golden List")
// subnavItemsArray stands for the items order
// Can be optimized
func setupSubnavData(data *subnavData, c *gin.Context, title string, subnavItemsArray []string, subnavItemsDisplableNames map[string]string) string {
	data.Title = title

	for _, name := range subnavItemsArray {
		data.SubnavItems = append(data.SubnavItems, subnavItem{Name: subnavItemsDisplableNames[name]})
	}

	selectedItemName := subnavItemsDisplableNames[c.Param("subNavItem")]
	if selectedItemName == "" {
		selectedItemName = data.SubnavItems[0].Name
	} else {
		i := 0
		for _, name := range subnavItemsArray {
			if selectedItemName == subnavItemsDisplableNames[name] {
				break
			}
			i++
		}
		data.SelectedSubnavItemIndex = i
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

type lolGameReviewData struct {
	GameMode     string
	GameDuration string
	BlueTeam     lolTeamGameReviewData
	RedTeam      lolTeamGameReviewData
}

type lolTeamGameReviewData struct {
	Players []lolPlayerGameReviewData
}

type lolPlayerGameReviewData struct {
	Version         string
	SummonerName    string
	ChampionName    string
	SummonerSpell1  string
	SummonerSpell2  string
	CreepScore      int
	GoldEarned      string
	Kills           int
	Deaths          int
	Assists         int
	KDA             string
	WardsPlaced     int
	WardsKilled     int
	PinkWardsBought int
	VisionScore     int
	Trinket         string
	Items           []string
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
