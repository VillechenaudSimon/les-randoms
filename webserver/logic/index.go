package logic

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleIndexRoute(c *gin.Context) {
	session := GetSession(c)

	data := indexData{}

	SetupNavData(&data.LayoutData.NavData, session)

	data.LayoutData.SubnavData.Title = "Index"

	SetupContentHeaderData(&data.ContentHeaderData, session)
	data.ContentHeaderData.Title = "Test"

	c.HTML(http.StatusFound, "index.tmpl.html", data)
}

func RedirectToIndex(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}
