package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handlePlayersRoute(c *gin.Context) {
	c.HTML(http.StatusOK, "players.tmpl.html", nil)
}
