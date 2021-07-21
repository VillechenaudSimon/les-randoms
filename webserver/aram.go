package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleAramRoute(c *gin.Context) {
	c.HTML(http.StatusOK, "aram.tmpl.html", nil)
}
