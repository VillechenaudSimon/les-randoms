package logic

import (
	"errors"
	"fmt"
	"les-randoms/utils"
	"os"
	"strings"
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
	if strings.Contains(input, "spotify") {
		// The program consider that if "spotify" in the given input, the user wants to use this source
		return bot.appendSpotifyPlaylistToQueue(gId, input)
	} else {
		// Otherwise the default youtube source is used
		client := youtube.Client{}
		id := ParseYoutubeId(input)
		playlist, err := client.GetPlaylist(id)
		if err != nil {
			time.Sleep(time.Millisecond * 50)
			video, err := client.GetVideo(id)
			if err != nil {
				return err
			}
			return bot.appendYoutubeVideoToQueue(gId, video)
		}
		return bot.appendYoutubePlaylistToQueue(gId, playlist, true)
	}
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

func buildMusicPath(i *MusicInfos) string {
	if i.Source == MusicInfosSources.Youtube {
		return buildYoutubeMusicPath(i.Id)
	} else if i.Source == MusicInfosSources.Spotify {
		return buildSpotifyMusicPath(i.Id)
	}
	utils.LogError("unknown music info source of id : " + fmt.Sprint(i.Source))
	return ""
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
		file = nil
		err = errors.New("unknown music source")
		if i.Source == MusicInfosSources.Youtube {
			file, err = bot.DownloadMusicFromYoutube(client, i)
		} else if i.Source == MusicInfosSources.Spotify {
			file, err = bot.DownloadMusicFromSpotify()
		}
		if err != nil {
			return err
		}
	} else if err == nil {
		bot.Log("Music found in cache (" + fmt.Sprint(i.Source) + ": " + i.Id + ")")
	} else {
		return err
	}
	defer file.Close()
	return nil
}
