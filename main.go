package main

import (
	"les-randoms/database"
	"les-randoms/discord-bot/bot"
	"les-randoms/webserver"
)

func main() {
	go bot.Start()

	database.OpenDatabase()
	defer database.CloseDatabase()

	webserver.StartWebServer()
}
