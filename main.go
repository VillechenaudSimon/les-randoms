package main

import (
	"les-randoms/database"
	"les-randoms/webserver"
)

func main() {
	database.OpenDatabase()
	defer database.Database.Close()

	webserver.StartWebServer()
}
