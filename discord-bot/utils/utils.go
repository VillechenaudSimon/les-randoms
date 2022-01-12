package utils

import (
	"les-randoms/utils"

	"github.com/bwmarrin/discordgo"
)

func BotLog(msg string) {
	utils.LogClassic("[DISCORD-BOT] " + msg)
}

func BotCommandError(msg string, s *discordgo.Session, m *discordgo.MessageCreate) {
	_, _ = s.ChannelMessageSend(m.ChannelID, "An error occured while executing your command. (Sorry Master <:sardAYAYA:657703982839365637>)\nError details : "+msg)
	BotError(msg)
}

func BotError(msg string) {
	utils.LogError("[DISCORD-BOT] " + msg)
}
