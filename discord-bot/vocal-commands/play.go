package vocalcommands

import (
	botutils "les-randoms/discord-bot/utils"
	"time"

	"github.com/bwmarrin/discordgo"
)

func CommandPlay(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, m.Author.Username+" asked me to join the vocal channel")
	if err != nil {
		botutils.BotCommandError(err.Error(), s, m)
		return
	}

	voiceState, err := s.State.VoiceState(m.GuildID, m.Author.ID)
	if err != nil {
		botutils.BotCommandError(err.Error(), s, m)
		return
	}

	voiceConnection, err := s.ChannelVoiceJoin(m.GuildID, voiceState.ChannelID, false, true)
	if err != nil {
		botutils.BotCommandError(err.Error(), s, m)
		return
	}

	voiceConnection.Speaking(true)

	time.Sleep(3 * time.Second)

	voiceConnection.Speaking(false)

	err = voiceConnection.Disconnect()
	if err != nil {
		botutils.BotCommandError(err.Error(), s, m)
		return
	}
}
