package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleIndexRoute(c *gin.Context) {
	session := getSession(c)

	if isAuthentified(session) {
		RedirectToAuth(c)
	}

	data := &indexData{}
	data.LayoutData.SubnavData.Title = "Index"

	c.HTML(http.StatusOK, "index.tmpl.html", data)
}
