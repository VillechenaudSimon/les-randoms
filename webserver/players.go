package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handlePlayersRoute(c *gin.Context) {
	session := getSession(c)

	if isAuthentified(session) {
		RedirectToAuth(c)
	}

	data := &playersData{}
	data.LayoutData.SubnavData.Title = "Player Analyser"

	c.HTML(http.StatusOK, "players.tmpl.html", data)
}
