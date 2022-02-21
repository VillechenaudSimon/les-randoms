package vocalcommands

import (
	botutils "les-randoms/discord-bot/utils"
	"les-randoms/utils"
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

	voiceConnection, err := s.ChannelVoiceJoin(m.GuildID, voiceState.ChannelID, false, false)
	if err != nil {
		botutils.BotCommandError(err.Error(), s, m)
		return
	}

	time.Sleep(250 * time.Millisecond)

	err = voiceConnection.Speaking(true)
	if err != nil {
		botutils.BotCommandError(err.Error(), s, m)
		return
	}

	time.Sleep(250 * time.Millisecond)

	DCA(voiceConnection, utils.WEBSITE_URL+"/static/musics/playing.mp3")

	time.Sleep(250 * time.Millisecond)

	err = voiceConnection.Speaking(false)
	if err != nil {
		botutils.BotCommandError(err.Error(), s, m)
		return
	}

	time.Sleep(250 * time.Millisecond)

	err = voiceConnection.Disconnect()
	if err != nil {
		botutils.BotCommandError(err.Error(), s, m)
		return
	}
}
