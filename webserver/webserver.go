package webserver

import (
	"errors"
	"les-randoms/utils"
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
	router.GET("/", handleIndexRoute)
	router.GET("/aram", handleAramRoute)
	router.POST("/aram", handleAramRoute)
	router.GET("/players", handlePlayersRoute)
	router.GET("/database", handleDatabaseRoute)
	router.POST("/database", handleDatabaseRoute)
}
