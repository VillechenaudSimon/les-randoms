package discordbot

import (
	"les-randoms/discord-bot/logic"

	"github.com/bwmarrin/discordgo"
)

func CommandUnknown(bot *logic.DiscordBot, m *discordgo.MessageCreate) {
	_, _ = bot.DiscordSession.ChannelMessageSend(m.ChannelID, "You have entered an unknown command. You must be really smart <:YEP:930600118787383336>")
}
