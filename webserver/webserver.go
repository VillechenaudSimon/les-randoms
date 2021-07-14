package webserver

import (
	"errors"
	"les-randoms/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func StartWebServer() {
	port := os.Getenv("PORT")
	if port == "" {
		utils.HandlePanicError(errors.New("$PORT must be set"))
	}

	router := gin.New()

	setupRouter(router)

	setupRoutes(router)

	router.Run(":" + port)
}

func setupRouter(router *gin.Engine) {
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")
}

func setupRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})
	router.GET("/lol", func(c *gin.Context) {
		c.HTML(http.StatusOK, "lol-index.tmpl.html", nil)
	})
	router.GET("/aram", func(c *gin.Context) {
		c.HTML(http.StatusOK, "aram.tmpl.html", nil)
	})
	router.GET("/players", func(c *gin.Context) {
		c.HTML(http.StatusOK, "players.tmpl.html", nil)
	})
}
