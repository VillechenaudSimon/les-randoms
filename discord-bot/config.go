package discordbot

import (
	"les-randoms/discord-bot/logic"
	"les-randoms/utils"
	"os"
)

func init() {
	var prefix string
	if utils.DebugMode {
		prefix = "k!"
	} else {
		prefix = "!"
	}
	Bot = logic.New(prefix, os.Getenv("DISCORD_BOT_TOKEN"), "784039117264388128")
}
