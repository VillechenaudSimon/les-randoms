package logic

import (
	"les-randoms/utils"

	"github.com/bwmarrin/discordgo"
)

func (bot *DiscordBot) Log(msg string) {
	utils.LogClassic("[DISCORD-BOT] " + msg)
}

func (bot *DiscordBot) CommandError(msg string, m *discordgo.MessageCreate) {
	id := bot.LogChannelId
	if m == nil {
		id = m.ChannelID
	}
	_, _ = bot.DiscordSession.ChannelMessageSend(id, "An error occured while executing your command. (Sorry Master <:sardAYAYA:657703982839365637>)\nError details : "+msg)
	bot.LogError(msg)
}

func (bot *DiscordBot) LogError(msg string) {
	utils.LogError("[DISCORD-BOT] " + msg)
}
