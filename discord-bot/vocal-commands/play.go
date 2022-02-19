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

	voiceConnection, err := s.ChannelVoiceJoin(m.GuildID, voiceState.ChannelID, false, false)
	if err != nil {
		botutils.BotCommandError(err.Error(), s, m)
		return
	}

	time.Sleep(1 * time.Second)

	//utils.LogDebug("voiceConnection.Ready:" + fmt.Sprint(voiceConnection.Ready))

	err = voiceConnection.Speaking(true)
	if err != nil {
		botutils.BotCommandError(err.Error(), s, m)
		return
	}

	//time.Sleep(100 * time.Millisecond)

	//DCA(voiceConnection, "playing.mp3")
	/*for _, buff := range buffer {
		voiceConnection.OpusSend <- buff
	}*/

	time.Sleep(1 * time.Second)
	time.Sleep(1 * time.Second)

	err = voiceConnection.Speaking(false)
	if err != nil {
		botutils.BotCommandError(err.Error(), s, m)
		return
	}

	time.Sleep(100 * time.Millisecond)

	err = voiceConnection.Disconnect()
	if err != nil {
		botutils.BotCommandError(err.Error(), s, m)
		return
	}
}
