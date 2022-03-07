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
	ProfileParameters struct {
		SummonerName string
	}
	LastGameParameters struct {
		SummonerName string
	}
	LolGameReviewData        lolGameReviewData
	LolProfileData           lolProfileData
	LadderChampPoolTableData customTableData
	LadderTableData          customTableData
}

type databaseData struct {
	LayoutData        layoutData
	ContentHeaderData contentHeaderData
	EntityTableData   customTableData
	SelectParameters  struct {
		SelectQueryBody string
	}
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
	lolNavItems = append(lolNavItems, navItem{IsGroup: false, Href: "/lol/players", ImgSrc: "/static/images/MPengu.png"})
	if getAccessStatus(s, "/lol/aram") > database.RightTypes.Hidden {
		lolNavItems = append(lolNavItems, navItem{IsGroup: false, Href: "/lol/aram", ImgSrc: "/static/images/HowlingAbyssIcon.png"})
	}
	data.NavItems = append(data.NavItems, navItem{IsGroup: true, ImgSrc: "/static/images/lol.ico", NavGroupItems: lolNavItems})

	if getAccessStatus(s, "/discord-bot") > database.RightTypes.Hidden {
		discordBotNavItems := make([]navItem, 0)
		discordBotNavItems = append(discordBotNavItems, navItem{IsGroup: false, Href: "/discord-bot/music", ImgSrc: "/static/images/music_note.png"})
		data.NavItems = append(data.NavItems, navItem{IsGroup: true, ImgSrc: "/static/images/discord.ico", NavGroupItems: discordBotNavItems})
	}

	if getAccessStatus(s, "/database") > database.RightTypes.Hidden {
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
	HeaderList      []string
	ColumnTypes     []customTableColumnType
	ItemList        []tableItemData
	SortColumnIndex int
	SortOrder       int // 0 Means descending order and 1 Means ascending order
}

type customTableColumnType int

const (
	customTableColumnTypeText   customTableColumnType = iota // 0
	customTableColumnTypeNumber                              // 1
	customTableColumnTypeImage                               // 2
	customTableColumnTypeDate                                // 3
)

func newCustomTableDataFromDBStruct(structType reflect.Type, dbStructs []database.DBStruct) customTableData {
	data := customTableData{}
	for i := 0; i < structType.NumField(); i++ {
		data.HeaderList = append(data.HeaderList, structType.Field(i).Name)
	}

	for _, dbStruct := range dbStructs {
		data.ItemList = append(data.ItemList, tableItemData{FieldList: dbStruct.ToStringSlice()})
	}
	data.SortColumnIndex = -1

	data.ColumnTypes = make([]customTableColumnType, len(data.HeaderList))

	for i := 0; i < len(data.ColumnTypes); i++ {
		data.ColumnTypes[i] = customTableColumnTypeText
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

type lolProfileData struct {
	SummonerName string
}

type discordBotMusicData struct {
	LayoutData              layoutData
	ContentHeaderData       contentHeaderData
	DiscordBotMusicPlayData discordBotMusicPlayData
}

type discordBotMusicPlayData struct {
	CurrentPlayStatus bool
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
