package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handlePlayersRoute(c *gin.Context) {
	session := getSession(c)

	if isNotAuthentified(session) {
		redirectToAuth(c)
	}

	data := playersData{}

	data.LayoutData.NavData = newNavData(session)

	data.LayoutData.SubnavData.Title = "Player Analyser"

	c.HTML(http.StatusOK, "players.tmpl.html", data)
}
