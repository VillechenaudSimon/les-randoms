package logic

import (
	"errors"
	"io"
	"les-randoms/utils"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

func (bot *DiscordBot) JoinChannel(guildID string, channelID string, mute bool, deaf bool) (*discordgo.VoiceConnection, error) {
	var err error
	vc, err := bot.DiscordSession.ChannelVoiceJoin(guildID, channelID, mute, deaf)
	if err != nil {
		return &discordgo.VoiceConnection{}, err
	}
	return vc, nil
}

func (bot *DiscordBot) JoinUserInChannel(guildId string, userId string, mute bool, deaf bool) (*discordgo.VoiceConnection, error) {
	voiceState, err := bot.DiscordSession.State.VoiceState(guildId, userId)
	if err != nil {
		return &discordgo.VoiceConnection{}, err
	}

	return bot.JoinChannel(voiceState.GuildID, voiceState.ChannelID, mute, deaf)
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
	if bot.queueAppender[gId] != nil {
		delete(bot.queueAppender, gId)
	}
	if bot.queuePlayer[gId] != nil {
		delete(bot.queuePlayer, gId)
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

func (bot *DiscordBot) GetCurrentTitle(gId string) string {
	if bot.DiscordSession.VoiceConnections[gId] == nil {
		return "Not Connected"
	}
	is := bot.musicQueues[gId]
	if is == nil || len(is) <= 0 {
		return "Not Playing Anything"
	}
	return is[0].Title
}

func (bot *DiscordBot) GetMusicQueue(gId string) []*MusicInfos {
	return bot.musicQueues[gId]
}

func (bot *DiscordBot) DCA(vc *discordgo.VoiceConnection, i *MusicInfos) error {
	gId := vc.GuildID
	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 96
	opts.Application = "lowdelay"
	opts.Volume = 32

	var err error

	bot.encodeSessions[gId], err = dca.EncodeFile(i.Url, opts)
	if err != nil {
		return errors.New(" Failed creating an encoding session: " + err.Error())
	}
	//bot.musicQueues[gId] = append(bot.musicQueues[gId], i)
	//v.encoder = encodeSession
	done := make(chan error)
	bot.streamingSessions[gId] = dca.NewStream(bot.encodeSessions[gId], vc, done)
	//v.stream = stream
	for err := range done {
		// Clean up incase something happened and ffmpeg is still running
		bot.encodeSessions[gId].Cleanup()

		delete(bot.streamingSessions, gId)
		delete(bot.encodeSessions, gId)
		if err != nil && err != io.EOF {
			return errors.New("An error occured " + err.Error())
		} else {
			return nil
		}
	}
	return errors.New("unreachable code")
}

var ErrorsPlay = struct {
	UnknownError           error
	UserNotInVoiceChannel  error
	TimeoutWaitingForVoice error
	QueueAppend            error
}{
	UnknownError:           errors.New("unknown error"),
	UserNotInVoiceChannel:  errors.New("user is not in a voice channel and ask for music"),
	TimeoutWaitingForVoice: errors.New("timeout waiting for voice -> probably missing rights to join channel"),
	QueueAppend:            errors.New("the append to queue function gone wrong"),
}

func (bot *DiscordBot) PlayOrder(gId string, uId string, songInput string) error {
	botVoiceState, err := bot.DiscordSession.State.VoiceState(gId, bot.Id)
	if err != nil {
		utils.LogError(err.Error())
		return ErrorsPlay.UnknownError
	}
	userVoiceState, err := bot.DiscordSession.State.VoiceState(gId, uId)
	if err != nil {
		utils.LogError(err.Error())
		return ErrorsPlay.UnknownError
	}
	if botVoiceState.ChannelID != userVoiceState.ChannelID { // If bot is not currently in the same voice channel as the user
		_, err := bot.JoinUserInChannel(gId, uId, false, true)
		if err != nil {
			if errors.Is(err, discordgo.ErrStateNotFound) {
				return ErrorsPlay.UserNotInVoiceChannel
			} else if err.Error() == "timeout waiting for voice" {
				bot.DiscordSession.VoiceConnections[gId] = nil // If the connection is wrong we delete it
				return ErrorsPlay.TimeoutWaitingForVoice
			} else {
				utils.LogError(err.Error())
				return ErrorsPlay.UnknownError
			}
		}

	}

	if bot.GetMusicQueue(gId) == nil { // If the queue is not playing, start it
		bot.PlayQueue(bot.DiscordSession.VoiceConnections[gId])
	}

	err = bot.AppendQueueFromInput(gId, songInput)
	if err != nil {
		utils.LogError(err.Error())
		return ErrorsPlay.QueueAppend
	}
	return nil
}
