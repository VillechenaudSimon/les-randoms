package logic

import (
	"errors"
	"les-randoms/utils"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kkdai/youtube/v2"
)

func (bot *DiscordBot) PlayQueue(vc *discordgo.VoiceConnection) {
	client := youtube.Client{}
	gId := vc.GuildID
	bot.musicQueues[gId] = make([]*MusicInfos, 0)
	bot.DiscordSession.VoiceConnections[gId] = vc
	bot.queueAppender[gId] = make(chan []*MusicInfos)
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
	go func() {
		for bot.DiscordSession.VoiceConnections[gId] != nil {
			//utils.LogDebug("Waiting: " + fmt.Sprint(bot.queueAppender[gId]))
			select {
			case s := <-bot.queueAppender[gId]:
				//utils.LogDebug("Append Signal received")
				go func() {
					bot.musicQueues[gId] = append(bot.musicQueues[gId], s...)
					if bot.streamingSessions[gId] == nil { // If not streaming a music, start to play the queue
						//utils.LogDebug("FromAppend: Sending play signal to " + fmt.Sprint(bot.queuePlayer[gId]))
						bot.queuePlayer[gId] <- bot.musicQueues[gId][0]
					}
				}()
			case i := <-bot.queuePlayer[gId]:
				//utils.LogDebug("Play Signal received")
				go func() {
					err := bot.downloadIfNecesary(&client, i)
					if err == nil {
						bot.DCA(bot.DiscordSession.VoiceConnections[gId], i)
					} else {
						bot.LogError(err.Error())
					}
					bot.musicQueues[gId] = bot.musicQueues[gId][1:]
					//utils.LogDebug("EndPlay: Check length : " + fmt.Sprint(len(bot.musicQueues[gId])))
					if len(bot.musicQueues[gId]) > 0 {
						//utils.LogDebug("FromPlay: Sending play signal to " + fmt.Sprint(bot.queuePlayer[gId]))
						bot.queuePlayer[gId] <- bot.musicQueues[gId][0]
					}
				}()
			}
		}
		utils.LogNotNilError(bot.Disconnect(gId))
		//utils.LogDebug("End of PlayQueue()")
	}()
}

func (bot *DiscordBot) AppendQueueFromInput(gId string, input string) error {
	client := youtube.Client{}
	id := ParseYoutubeId(input)
	playlist, err := client.GetPlaylist(id)
	if err != nil {
		time.Sleep(time.Millisecond * 50)
		video, err := client.GetVideo(id)
		if err != nil {
			return err
		}
		return bot.appendVideoToQueue(gId, video)
	}
	return bot.appendPlaylistToQueue(gId, playlist, true)
}

func (bot *DiscordBot) AppendEltQueue(gId string, i *MusicInfos) error {
	if bot.musicQueues[gId] == nil {
		return errors.New("no music queue detected in this guild")
	}
	//utils.LogDebug("Sending append signal to " + fmt.Sprint(bot.queueAppender[gId]))
	s := make([]*MusicInfos, 1)
	s[0] = i
	bot.queueAppender[gId] <- s
	return nil
}

func (bot *DiscordBot) AppendEltsQueue(gId string, s []*MusicInfos) error {
	if bot.musicQueues[gId] == nil {
		return errors.New("no music queue detected in this guild")
	}
	bot.queueAppender[gId] <- s
	return nil
}

func (bot *DiscordBot) appendVideoToQueue(gId string, v *youtube.Video) error {
	return bot.AppendEltQueue(gId, NewMusicInfos(v.ID, v.Title, buildVideoURL(v.ID)))
}

func (bot *DiscordBot) appendPlaylistToQueue(gId string, p *youtube.Playlist, shuffle bool) error {
	s := make([]*MusicInfos, 0)
	for _, e := range p.Videos {
		s = append(s, NewMusicInfos(e.ID, e.Title, buildVideoURL(e.ID)))
	}
	if shuffle {
		rand.Shuffle(len(s), func(i, j int) {
			s[i], s[j] = s[j], s[i]
		})
	}
	return bot.AppendEltsQueue(gId, s)
}

func buildVideoURL(yId string) string {
	return musicCacheFolderPath + yId + ".m4a"
}
