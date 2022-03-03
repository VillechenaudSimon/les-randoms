package discordbot

import (
	"les-randoms/discord-bot/logic"
	"les-randoms/utils"
	"os"
)

func init() {
	Bot = &logic.DiscordBot{
		Token:        os.Getenv("DISCORD_BOT_TOKEN"),
		LogChannelId: "784039117264388128",
	}

	if utils.DebugMode {
		Bot.Prefix = "k!"
	} else {
		Bot.Prefix = "!"
	}
}
