package discordbot

import (
	"github.com/bwmarrin/discordgo"
)

func (bot *DiscordBot) CommandPlay(m *discordgo.MessageCreate) {
	_, err := bot.DiscordSession.ChannelMessageSend(m.ChannelID, m.Author.Username+" asked me to join the vocal channel")
	if err != nil {
		bot.CommandError(err.Error(), m)
		return
	}

	voiceState, err := bot.DiscordSession.State.VoiceState(m.GuildID, m.Author.ID)
	if err != nil {
		bot.CommandError(err.Error(), m)
		return
	}

	err = bot.JoinChannel(m.GuildID, voiceState.ChannelID, false, true)
	if err != nil {
		bot.CommandError(err.Error(), m)
		return
	}

	err = bot.PlayMusic()
	if err != nil {
		bot.CommandError(err.Error(), m)
	}

	err = bot.VoiceConnection.Disconnect()
	if err != nil {
		bot.CommandError(err.Error(), m)
		return
	}
}
