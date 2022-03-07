package webserver

import (
	"errors"
	"les-randoms/database"
	discordbot "les-randoms/discord-bot/web-exec"
	"les-randoms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleMusicRoute(c *gin.Context) {
	session := getSession(c)

	if c.Param("subNavItem") == "" {
		c.Redirect(http.StatusFound, "/discord-bot/music/playing")
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

	selectedItemName := setupSubnavData(&data.LayoutData.SubnavData, c, "Music Bot", []string{"Playing"}, map[string]string{"Playing": "Playing"})

	setupContentHeaderData(&data.ContentHeaderData, session)
	data.ContentHeaderData.Title = selectedItemName

	switch data.LayoutData.SubnavData.SelectedSubnavItemIndex {
	case 0:
		if c.Param("order") != "" {
			err := handlePlayingOrder(c)
			if err != nil {
				utils.LogError("Happened while executing discord bot playing order " + utils.Esc(c.Param("order")) + " : " + err.Error())
			} else {
				utils.LogClassic("Discord bot playing order " + utils.Esc(c.Param("order")) + " successfully executed")
				return
			}
		}
		if setupPlayingData(&data) != nil {
			c.Redirect(http.StatusFound, "/discord-bot/music/playing")
		}
	}

	c.HTML(http.StatusFound, "discord-bot-music.tmpl.html", data)
}

func handlePlayingOrder(c *gin.Context) error {
	order := c.Param("order")
	if order == "resume" {
		return discordbot.ExecuteMusicResume()
	} else if order == "pause" {
		return discordbot.ExecuteMusicPause()
	} else {
		return errors.New("Unknown order")
	}
}

func setupPlayingData(data *discordBotMusicData) error {
	data.DiscordBotMusicPlayData.CurrentPlayStatus = discordbot.GetPlayStatus()
	return nil
}
