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

func (bot *DiscordBot) PlayQueue(vc *discordgo.VoiceConnection) error {
	gId := vc.GuildID
	bot.musicQueues[gId] = make([]*MusicInfos, 0)
	bot.DiscordSession.VoiceConnections[gId] = vc
	bot.queueAppender[gId] = make(chan *MusicInfos)
	bot.queuePlayer[gId] = make(chan *MusicInfos)
	/*for i := range bot.queueSignals[gId] {
		utils.LogDebug("Signal received")
		bot.musicQueues[gId] = append(bot.musicQueues[gId], i)
		if bot.streamingSessions[gId] == nil { // if not playing a music
			go func() {
				if len(bot.musicQueues[gId]) > 0 {
					bot.DCA(bot.DiscordSession.VoiceConnections[gId], bot.musicQueues[gId][0])
					bot.musicQueues[gId] = bot.musicQueues[gId][1:]
				}
			}()
		}
	}*/
	for bot.DiscordSession.VoiceConnections[gId] != nil {
		//utils.LogDebug("Waiting: " + fmt.Sprint(bot.queueAppender[gId]))
		select {
		case i := <-bot.queueAppender[gId]:
			//utils.LogDebug("Append Signal received")
			go func() {
				bot.musicQueues[gId] = append(bot.musicQueues[gId], i)
				if bot.streamingSessions[gId] == nil { // If not streaming a music, start to play the queue
					//utils.LogDebug("FromAppend: Sending play signal to " + fmt.Sprint(bot.queuePlayer[gId]))
					bot.queuePlayer[gId] <- bot.musicQueues[gId][0]
				}
			}()
		case i := <-bot.queuePlayer[gId]:
			//utils.LogDebug("Play Signal received")
			go func() {
				bot.DCA(bot.DiscordSession.VoiceConnections[gId], i)
				bot.musicQueues[gId] = bot.musicQueues[gId][1:]
				//utils.LogDebug("EndPlay: Check length : " + fmt.Sprint(len(bot.musicQueues[gId])))
				if len(bot.musicQueues[gId]) > 0 {
					//utils.LogDebug("FromPlay: Sending play signal to " + fmt.Sprint(bot.queuePlayer[gId]))
					bot.queuePlayer[gId] <- bot.musicQueues[gId][0]
				}
			}()

		}
	}
	//utils.LogDebug("End of PlayQueue()")
	return bot.Disconnect(gId)
}

func (bot *DiscordBot) TestPlayMusic(vc *discordgo.VoiceConnection) error {
	time.Sleep(250 * time.Millisecond)

	err := vc.Speaking(true)
	if err != nil {
		return err
	}

	time.Sleep(250 * time.Millisecond)

	bot.musicQueues[vc.GuildID] = append(bot.musicQueues[vc.GuildID], &MusicInfos{Title: "TeST", Url: "playing.mp3"})

	//err = bot.DCA(vc, "https://www.youtube.com/watch?v=hRGIrrjuLYA", &MusicInfos{Title: "TeST"})
	err = bot.DCA(vc, bot.musicQueues[vc.GuildID][0])
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
	if bot.queueAppender[gId] != nil {
		delete(bot.queueAppender, gId)
	}
	if bot.queuePlayer[gId] != nil {
		delete(bot.queuePlayer, gId)
	}
	return err
}

func (bot *DiscordBot) AppendQueue(gId string, i *MusicInfos) error {
	if bot.musicQueues[gId] == nil {
		return errors.New("no music queue detected in this guild")
	}
	//utils.LogDebug("Sending append signal to " + fmt.Sprint(bot.queueAppender[gId]))
	bot.queueAppender[gId] <- i
	return nil
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
	is := bot.musicQueues[gId]
	if is == nil {
		return ""
	}
	if len(is) <= 0 {
		return ""
	}
	return is[0].Title
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
