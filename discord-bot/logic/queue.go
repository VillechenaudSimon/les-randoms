package logic

import (
	"errors"
	"fmt"
	"io"
	"les-randoms/utils"
	"les-randoms/ytinterface"
	"math/rand"
	"os"
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
						if bot.DiscordSession.VoiceConnections[gId] != nil {
							bot.DCA(bot.DiscordSession.VoiceConnections[gId], i)
						} else { // If there is no vocal connection the bot disconnects to prevent bugs later
							bot.Disconnect(gId)
						}
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
	return bot.AppendEltQueue(gId, NewMusicInfos(v.ID, v.Title, buildYoutubeMusicPath(v.ID), MusicInfosSources.Youtube))
}

func (bot *DiscordBot) appendPlaylistToQueue(gId string, p *youtube.Playlist, shuffle bool) error {
	s := make([]*MusicInfos, 0)
	for _, e := range p.Videos {
		s = append(s, NewMusicInfos(e.ID, e.Title, buildYoutubeMusicPath(e.ID), MusicInfosSources.Youtube))
	}
	if shuffle {
		rand.Shuffle(len(s), func(i, j int) {
			s[i], s[j] = s[j], s[i]
		})
	}
	return bot.AppendEltsQueue(gId, s)
}

func buildMusicPath(i *MusicInfos) string {
	if i.Source == MusicInfosSources.Youtube {
		return buildYoutubeMusicPath(i.Id)
	} else if i.Source == MusicInfosSources.Spotify {
		return buildSpotifyMusicPath(i.Id)
	}
	utils.LogError("unknown music info source of id : " + fmt.Sprint(i.Source))
	return ""
}

func buildYoutubeMusicPath(yId string) string {
	return musicCacheFolderPath + musicCacheYoutubeSubfolder + yId + ".m4a"
}

func buildSpotifyMusicPath(sId string) string {
	return musicCacheFolderPath + musicCacheSpotifySubfolder + sId + ".m4a"
}

func setupFolders() error {
	err := os.Mkdir(musicCacheFolderPath, os.ModeAppend)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	err = os.Mkdir(musicCacheFolderPath+musicCacheYoutubeSubfolder, os.ModeAppend)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	err = os.Mkdir(musicCacheFolderPath+musicCacheSpotifySubfolder, os.ModeAppend)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	return nil
}

// TODO : Refactor this function to separate the youtube part in youtube.go and start a spotify.go part
func (bot *DiscordBot) downloadIfNecesary(client *youtube.Client, i *MusicInfos) error {
	err := setupFolders()
	if err != nil {
		return err
	}
	file, err := os.Open(buildMusicPath(i))
	if errors.Is(err, os.ErrNotExist) {
		file, err = os.Create(buildMusicPath(i))
		if err != nil {
			return err
		}
		// Download as file is mandatory since stream of more than 2m40s are ended without error thrown (probably because of youtube limitations)
		video, err := client.GetVideo(i.Id)
		if err != nil {
			return err
		}
		format, err := ytinterface.GetBestAudioOnlyFormat(video.Formats)
		if err != nil {
			return err
		}
		stream, _, err := client.GetStream(video, format)
		if err != nil {
			return err
		}
		bot.Log("Music download start (" + i.Id + ")")
		_, err = io.Copy(file, stream)
		bot.Log("Music download end (" + i.Id + ")")
		if err != nil {
			return err
		}
	} else if err == nil {
		bot.Log("Found in cache video of id : " + i.Id)
	} else {
		return err
	}
	defer file.Close()
	return nil
}
