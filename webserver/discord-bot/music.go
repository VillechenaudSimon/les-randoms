package discordbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"les-randoms/database"
	"les-randoms/discord-bot/logic"
	webexec "les-randoms/discord-bot/web-exec"
	"les-randoms/utils"
	webserver "les-randoms/webserver/logic"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func SetupRoutes() {
	discordbot := webserver.Router.Group("/discord-bot")

	music := discordbot.Group("/music")
	music.GET("", handleMusicRoute)
	music.POST("", handleMusicRoute)
	music.GET("/:subNavItem", handleMusicRoute)
	music.POST("/:subNavItem", handleMusicRoute)
	music.GET("/:subNavItem/:order", handleMusicRoute)
	music.POST("/:subNavItem/:order", handleMusicRoute)
	music.GET("/playing/ws", handlePlayingWs)
}

func handleMusicRoute(c *gin.Context) {
	session := webserver.GetSession(c)

	if c.Param("subNavItem") == "" {
		c.Redirect(http.StatusFound, "/discord-bot/music/playing")
	}

	if webserver.IsNotAuthentified(session) {
		webserver.RedirectToAuth(c)
		return
	}

	if webserver.GetAccessStatus(session, "/discord-bot") <= database.RightTypes.Forbidden {
		webserver.RedirectToIndex(c)
		return
	}

	data := discordBotMusicData{}

	webserver.SetupNavData(&data.LayoutData.NavData, session)

	selectedItemName := webserver.SetupSubnavData(&data.LayoutData.SubnavData, c, "Music Bot", []string{"Playing"}, map[string]string{"Playing": "Playing"})

	webserver.SetupContentHeaderData(&data.ContentHeaderData, session)
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

type discordBotMusicData struct {
	LayoutData              webserver.LayoutData
	ContentHeaderData       webserver.ContentHeaderData
	DiscordBotMusicPlayData discordBotMusicPlayData
}

type discordBotMusicPlayData struct {
	CurrentPlayStatus bool
	CurrentTime       string
	CurrentTitle      string
}

func handlePlayingOrder(c *gin.Context) error {
	order := c.Param("order")
	if order == "resume" {
		return webexec.ExecuteMusicResume()
	} else if order == "pause" {
		return webexec.ExecuteMusicPause()
	} else if order == "play" {
		return webexec.ExecuteMusicPlay(webserver.GetDiscordId(webserver.GetSession(c)), c.Request.PostFormValue("param1"))
	} else {
		return errors.New("unknown order")
	}
}

func setupPlayingData(data *discordBotMusicData) error {
	data.DiscordBotMusicPlayData.CurrentPlayStatus = webexec.GetPlayStatus()
	currentMusicDuration := webexec.GetCurrentTime()
	data.DiscordBotMusicPlayData.CurrentTime = fmt.Sprintf("%02d:%02d", int(currentMusicDuration.Minutes()), int(currentMusicDuration.Seconds())%60)
	data.DiscordBotMusicPlayData.CurrentTitle = webexec.GetCurrentTitle()
	return nil
}

func handlePlayingWs(c *gin.Context) {
	conn, err := webserver.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.LogError(err.Error())
		return
	}
	defer conn.Close()
	conn.SetCloseHandler(func(code int, text string) error {
		//utils.LogDebug("Connection closed : " + fmt.Sprint(code) + " : " + text)
		return nil
	})
	time.Sleep(time.Second)
	counter := 0
	for {
		//mt, msg, err := conn.ReadMessage()
		//if err != nil {
		//	utils.LogError(err.Error())
		//}
		//if mt != websocket.TextMessage {
		//	utils.LogError(err.Error())
		//}
		//var v message
		//json.Unmarshal(msg, &v)
		time.Sleep(time.Millisecond * 250)
		currentMusicDuration := webexec.GetCurrentTime()
		data := struct {
			DataType     int // 0 For just PlayStatus CurrentTime and CurrentTitle
			PlayStatus   bool
			CurrentTime  string
			CurrentTitle string
			Queue        []*logic.MusicInfos
		}{
			DataType:     0,
			PlayStatus:   webexec.GetPlayStatus(),
			CurrentTime:  fmt.Sprintf("%02d:%02d", int(currentMusicDuration.Minutes()), int(currentMusicDuration.Seconds())%60),
			CurrentTitle: webexec.GetCurrentTitle(),
		}
		if counter >= 20 { // Every 20 * 250ms = 5 seconds the whole queue is sent to fix desync problems
			data.DataType = 1 // 1 For 0 + Whole Queue informations
			data.Queue = webexec.GetMusicQueue()
			counter = 0
		}
		counter++
		jsonContent, err := json.Marshal(data)
		if err != nil {
			utils.LogError(err.Error())
			return
		}
		err = conn.WriteMessage(websocket.TextMessage, jsonContent)
		if err != nil {
			if !strings.Contains(err.Error(), "An established connection was aborted by the software in your host machine") {
				utils.LogError(err.Error())
			}
			return
		}
	}
}
