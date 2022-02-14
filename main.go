package main

import (
	"les-randoms/backgroundworker"
	"les-randoms/database"
	"les-randoms/discord-bot/bot"
	radbwrapper "les-randoms/radb-wrapper"
	"les-randoms/webserver"
)

func main() {
	go bot.Start()
	defer bot.Close()

	database.OpenDatabase()
	defer database.CloseDatabase()
	database.VerifyDatabase()

	go backgroundworker.Start()
	radbwrapper.SetupJobs()

	webserver.StartWebServer()
}
