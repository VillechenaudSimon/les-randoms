package webserver

import (
	"les-randoms/webserver/admin"
	discordbot "les-randoms/webserver/discord-bot"
	"les-randoms/webserver/logic"
	"les-randoms/webserver/lol"
)

func StartWebServer() {
	logic.CreateServer()

	lol.SetupRoutes()
	discordbot.SetupRoutes()
	admin.SetupRoutes()

	logic.RunServer()
}

func StopServer() {
	logic.StopServer()
}
