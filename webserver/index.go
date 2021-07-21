package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleIndexRoute(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl.html", nil)
}
