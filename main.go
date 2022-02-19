package main

import (
	"les-randoms/backgroundworker"
	"les-randoms/database"
	"les-randoms/discord-bot/bot"
	radbwrapper "les-randoms/radb-wrapper"
	"les-randoms/utils"
	"les-randoms/webserver"
)

var AppEnd chan bool

func main() {
	AppEnd := make(chan bool, 1)

	go bot.Start(AppEnd)

	database.OpenDatabase()
	database.VerifyDatabase()

	go backgroundworker.Start()
	radbwrapper.SetupJobs()

	go webserver.StartWebServer()

	<-AppEnd
	utils.LogSuccess("App successfully closed.")
}
