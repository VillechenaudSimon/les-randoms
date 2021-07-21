package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handlePlayersRoute(c *gin.Context) {
	data := &playersData{}
	data.LayoutData.SubnavData.Title = "Player Analyser"

	c.HTML(http.StatusOK, "players.tmpl.html", data)
}
