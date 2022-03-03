package discordbot

import "github.com/bwmarrin/discordgo"

func (bot *DiscordBot) CommandUnknown(m *discordgo.MessageCreate) {
	_, _ = bot.DiscordSession.ChannelMessageSend(m.ChannelID, "You have entered an unknown command. You must be really smart <:YEP:930600118787383336>")
}
