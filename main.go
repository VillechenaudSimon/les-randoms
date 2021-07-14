package main

import (
	"log"
	"net/http"
	"os"

	"les-randoms/database"
	"les-randoms/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.Foo(2)

	database.Test()

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

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

	router.Run(":" + port)
}
