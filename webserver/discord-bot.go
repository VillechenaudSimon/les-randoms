package webserver

import (
	"les-randoms/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleDiscordBotRoute(c *gin.Context) {
	session := getSession(c)

	if c.Param("subNavItem") == "" {
		c.Redirect(http.StatusFound, "/discord-bot/music/play")
	}

	if isNotAuthentified(session) {
		redirectToAuth(c)
		return
	}

	if getAccessStatus(session, "/discord-bot") <= database.RightTypes.Forbidden {
		redirectToIndex(c)
		return
	}

	data := discordBotMusicData{}

	setupNavData(&data.LayoutData.NavData, session)

	selectedItemName := setupSubnavData(&data.LayoutData.SubnavData, c, "Music Bot", []string{"Play"}, map[string]string{"Play": "Play"})

	setupContentHeaderData(&data.ContentHeaderData, session)
	data.ContentHeaderData.Title = selectedItemName

	switch data.LayoutData.SubnavData.SelectedSubnavItemIndex {
	case 0:
		if setupMusicData(&data) != nil {
			c.Redirect(http.StatusFound, "/discord-bot/music/play")
		}
	}

	c.HTML(http.StatusFound, "discord-bot-music.tmpl.html", data)
}

func setupMusicData(data *discordBotMusicData) error {
	return nil
}
