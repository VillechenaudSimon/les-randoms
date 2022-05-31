package logic

import (
	"les-randoms/database"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type indexData struct {
	LayoutData        LayoutData
	ContentHeaderData ContentHeaderData
}

type LayoutData struct {
	NavData    NavData
	SubnavData SubnavData
}

type NavData struct {
	NavItems []NavItem
}

type NavItem struct {
	IsGroup       bool
	Href          string
	ImgSrc        string
	NavGroupItems []NavItem
}

func SetupNavData(data *NavData, s *sessions.Session) {
	data.NavItems = append(data.NavItems, NavItem{IsGroup: false, Href: "/", ImgSrc: "/static/images/favicon.png"})

	lolNavItems := make([]NavItem, 0)
	lolNavItems = append(lolNavItems, NavItem{IsGroup: false, Href: "/lol/players", ImgSrc: "/static/images/MPengu.png"})
	if GetAccessStatus(s, "/lol/aram") > database.RightTypes.Hidden {
		lolNavItems = append(lolNavItems, NavItem{IsGroup: false, Href: "/lol/aram", ImgSrc: "/static/images/HowlingAbyssIcon.png"})
	}
	data.NavItems = append(data.NavItems, NavItem{IsGroup: true, ImgSrc: "/static/images/lol.ico", NavGroupItems: lolNavItems})

	if GetAccessStatus(s, "/discord-bot") > database.RightTypes.Hidden {
		discordBotNavItems := make([]NavItem, 0)
		discordBotNavItems = append(discordBotNavItems, NavItem{IsGroup: false, Href: "/discord-bot/music", ImgSrc: "/static/images/music_note.png"})
		data.NavItems = append(data.NavItems, NavItem{IsGroup: true, ImgSrc: "/static/images/discord.ico", NavGroupItems: discordBotNavItems})
	}

	if GetAccessStatus(s, "/database") > database.RightTypes.Hidden {
		data.NavItems = append(data.NavItems, NavItem{IsGroup: false, Href: "/database", ImgSrc: "/static/images/databaseConfig.png"})
	}
}

type SubnavData struct {
	Title                   string
	SubnavItems             []SubnavItem
	SelectedSubnavItemIndex int
}

// Returns the selectedItemName
// subNavItemsMap map the displayable name of the item from the 'not-displayable' name ("GoldenList" -> "Golden List")
// subNavItemsArray stands for the items order
// Can be optimized
func SetupSubnavData(data *SubnavData, c *gin.Context, title string, subNavItemsArray []string, subNavItemsDisplableNames map[string]string) string {
	data.Title = title

	for _, name := range subNavItemsArray {
		data.SubnavItems = append(data.SubnavItems, SubnavItem{Name: subNavItemsDisplableNames[name]})
	}

	selectedItemName := subNavItemsDisplableNames[c.Param("subNavItem")]
	if selectedItemName == "" {
		selectedItemName = data.SubnavItems[0].Name
	} else {
		i := 0
		for _, name := range subNavItemsArray {
			if selectedItemName == subNavItemsDisplableNames[name] {
				break
			}
			i++
		}
		data.SelectedSubnavItemIndex = i
	}
	return selectedItemName
}

type SubnavItem struct {
	Name string
}

type ContentHeaderData struct {
	Title          string
	IsAuthentified bool
	DiscordId      string
	Username       string
	AvatarId       string
}

func SetupContentHeaderData(data *ContentHeaderData, s *sessions.Session) {
	data.IsAuthentified = !IsNotAuthentified(s)
	data.DiscordId = GetDiscordId(s)
	data.Username = GetUsername(s)
	data.AvatarId = GetAvatarId(s)
	data.Title = "Default"
}

type CustomTableData struct {
	HeaderList      []string
	ColumnTypes     []CustomTableColumnType
	ItemList        []TableItemData
	SortColumnIndex int
	SortOrder       int // 0 Means descending order and 1 Means ascending order
}

type CustomTableColumnType int

const (
	CustomTableColumnTypeText   CustomTableColumnType = iota // 0
	CustomTableColumnTypeNumber                              // 1
	CustomTableColumnTypeImage                               // 2
	CustomTableColumnTypeDate                                // 3
)

type TableItemData struct {
	FieldList []string
}

/*
Unused in nav.tmpl.html :
{{ range $i, $NavItem := .NavItems }}
{{ if $NavItem.IsGroup }}
<div class="nav-group expanded">
  <div class="nav-group-header">
    <img src="{{ $NavItem.ImgSrc }}">
  </div>
  <div class="nav-group-content">
    {{ range $j, $navGroupItem := $NavItem.NavGroupItems }}
    <a class="nav-item" href="{{ $NavItem.Href }}">
      <img src="{{ $NavItem.ImgSrc }}" />
    </a>
    {{ end }}
  </div>
</div>
{{ else }}
<a class="nav-item" href="{{ $NavItem.Href }}">
  <img src="{{ $NavItem.ImgSrc }}" />
</a>
{{ end }}
{{ end }}

Unused in structures.go :
type navData struct {
	NavItems []NavItem
}
type NavItem struct {
	IsGroup       bool
	Href          string
	ImgSrc        string
	NavGroupItems []NavItem
}

*/
