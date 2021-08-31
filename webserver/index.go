package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleIndexRoute(c *gin.Context) {
	session := getSession(c)

	if isNotAuthentified(session) {
		redirectToAuth(c)
	}

	data := indexData{}

	data.LayoutData.NavData = newNavData(session)

	data.LayoutData.SubnavData.Title = "Index"

	c.HTML(http.StatusOK, "index.tmpl.html", data)
}

func redirectToIndex(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}
