package main

import (
	"les-randoms/backgroundworker"
	"les-randoms/database"
	"les-randoms/webserver"
)

func main() {
	database.OpenDatabase()
	defer database.CloseDatabase()
	database.VerifyDatabase()

	go backgroundworker.Start()

	webserver.StartWebServer()
}
