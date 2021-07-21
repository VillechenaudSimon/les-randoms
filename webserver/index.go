package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleIndexRoute(c *gin.Context) {
	data := &indexData{}
	data.LayoutData.SubnavData.Title = "Index"

	c.HTML(http.StatusOK, "index.tmpl.html", data)
}
