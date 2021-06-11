package main

import (
	"log"
	"net/http"
	"os"
	"database/sql"
	
	_ "github.com/go-sql-driver/mysql"
	"les-randoms/utils"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

  utils.Foo()
  
  db, err := sql.Open("mysql",
  		"217240:tD5w4$dA6$MC@tcp(127.0.0.1:3306)/BlackListItem")
  if err != nil {
    log.Fatal(err)
  }
  db.Ping()
  defer db.Close()

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.Run(":" + port)
}
