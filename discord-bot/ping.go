package discordbot

import (
	"les-randoms/discord-bot/logic"

	"github.com/bwmarrin/discordgo"
)

func CommandPing(bot *logic.DiscordBot, m *discordgo.MessageCreate) {
	_, _ = bot.DiscordSession.ChannelMessageSend(m.ChannelID, "Pong !")
}
