package logic

import (
	"errors"
	"io"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

func (bot *DiscordBot) JoinMessageChannel(m *discordgo.MessageCreate, mute bool, deaf bool) (*discordgo.VoiceConnection, error) {
	voiceState, err := bot.DiscordSession.State.VoiceState(m.GuildID, m.Author.ID)
	if err != nil {
		return &discordgo.VoiceConnection{}, err
	}

	vc, err := bot.JoinChannel(m.GuildID, voiceState.ChannelID, false, true)
	if err != nil {
		return &discordgo.VoiceConnection{}, err
	}
	return vc, nil
}

func (bot *DiscordBot) JoinChannel(guildID string, channelID string, mute bool, deaf bool) (*discordgo.VoiceConnection, error) {
	var err error
	vc, err := bot.DiscordSession.ChannelVoiceJoin(guildID, channelID, mute, deaf)
	if err != nil {
		return &discordgo.VoiceConnection{}, err
	}
	return vc, nil
}

func (bot *DiscordBot) PlayMusic(vc *discordgo.VoiceConnection) error {
	time.Sleep(250 * time.Millisecond)

	err := vc.Speaking(true)
	if err != nil {
		return err
	}

	time.Sleep(250 * time.Millisecond)

	err = bot.DCA(vc, "playing.mp3")
	if err != nil {
		return err
	}

	time.Sleep(250 * time.Millisecond)

	err = vc.Speaking(false)
	if err != nil {
		return err
	}

	time.Sleep(250 * time.Millisecond)
	return nil
}

func (bot *DiscordBot) PauseMusic(gId string) error {
	s := bot.streamingSessions[gId]
	if s == nil {
		return errors.New("bot is not playing music in this guild")
	}
	s.SetPaused(true)
	return nil
}

func (bot *DiscordBot) ResumeMusic(gId string) error {
	s := bot.streamingSessions[gId]
	if s == nil {
		return errors.New("bot is not playing music in this guild")
	}
	s.SetPaused(false)
	return nil
}

func (bot *DiscordBot) Disconnect(gId string) error {
	var err error
	if bot.DiscordSession.VoiceConnections[gId] != nil {
		err = bot.DiscordSession.VoiceConnections[gId].Disconnect()
		delete(bot.DiscordSession.VoiceConnections, gId)
	}
	if bot.encodeSessions[gId] != nil {
		bot.encodeSessions[gId].Cleanup()
		delete(bot.encodeSessions, gId)
	}
	if bot.streamingSessions[gId] != nil {
		delete(bot.streamingSessions, gId)
	}
	return err
}

// True for paused, false for playing
func (bot *DiscordBot) GetPlayStatus(gId string) bool {
	s := bot.streamingSessions[gId]
	if s == nil {
		return true
	}
	return s.Paused()
}

func (bot *DiscordBot) GetCurrentTime(gId string) time.Duration {
	s := bot.streamingSessions[gId]
	if s == nil {
		return -1
	}
	return s.PlaybackPosition()
}

func (bot *DiscordBot) DCA(vc *discordgo.VoiceConnection, url string) error {
	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 96
	opts.Application = "lowdelay"
	opts.Volume = 32

	var err error
	bot.encodeSessions[vc.GuildID], err = dca.EncodeFile(url, opts)
	if err != nil {
		return errors.New(" Failed creating an encoding session: " + err.Error())
	}
	//v.encoder = encodeSession
	done := make(chan error)
	bot.streamingSessions[vc.GuildID] = dca.NewStream(bot.encodeSessions[vc.GuildID], vc, done)
	//v.stream = stream
	for err := range done {
		// Clean up incase something happened and ffmpeg is still running
		bot.encodeSessions[vc.GuildID].Cleanup()
		if err != nil && err != io.EOF {
			return errors.New("An error occured " + err.Error())
		}
	}
	return nil
}
