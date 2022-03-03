package discordbot

import "github.com/bwmarrin/discordgo"

func (bot *DiscordBot) CommandPing(m *discordgo.MessageCreate) {
	_, _ = bot.DiscordSession.ChannelMessageSend(m.ChannelID, "Pong !")
}
