package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleAramRoute(c *gin.Context) {
	data := &aramData{}
	data.LayoutData.SubnavData.Title = "Aram Gaming"

	c.HTML(http.StatusOK, "aram.tmpl.html", data)
}
