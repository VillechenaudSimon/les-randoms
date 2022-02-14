package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleIndexRoute(c *gin.Context) {
	session := getSession(c)

	data := indexData{}

	setupNavData(&data.LayoutData.NavData, session)

	data.LayoutData.SubnavData.Title = "Index"

	setupContentHeaderData(&data.ContentHeaderData, session)
	data.ContentHeaderData.Title = "Test"

	c.HTML(http.StatusFound, "index.tmpl.html", data)
}

func redirectToIndex(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}
