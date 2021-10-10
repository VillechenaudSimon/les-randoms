package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handlePlayersRoute(c *gin.Context) {
	session := getSession(c)

	data := playersData{}

	data.LayoutData.NavData = newNavData(session)

	data.LayoutData.SubnavData.Title = "Player Analyser"

	data.ContentHeaderData = newContentHeaderData(session)
	data.ContentHeaderData.Title = "WORK IN PROGRESS"

	c.HTML(http.StatusOK, "players.tmpl.html", data)
}
