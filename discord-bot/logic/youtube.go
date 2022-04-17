package logic

import (
	"io"
	"les-randoms/ytinterface"
	"math/rand"
	"os"
	"strings"

	"github.com/kkdai/youtube/v2"
)

/* not used anymore
func (bot *DiscordBot) DownloadAndAppendQueue(gId string, vidId string) (*youtube.Video, error) {
	client := youtube.Client{}
	video, err := client.GetVideo(vidId)
	if err != nil {
		return nil, err
	}
	format, err := ytinterface.GetBestAudioOnlyFormat(video.Formats)
	if err != nil {
		return nil, err
	}
	// Download as file is mandatory since stream of more than 2m40s are ended without error thrown (probably because of youtube limitations)
	stream, _, err := client.GetStream(video, format)
	if err != nil {
		return nil, err
	}
	err = os.Mkdir(musicCacheFolderPath, os.ModeAppend)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return nil, err
	}
	file, err := os.Open(musicCacheFolderPath + vidId + ".m4a")
	if errors.Is(err, os.ErrNotExist) {
		file, err = os.Create(musicCacheFolderPath + vidId + ".m4a")
		if err != nil {
			return nil, err
		}
		utils.LogClassic("Downloading video of id : " + vidId)
		_, err = io.Copy(file, stream)
		if err != nil {
			return nil, err
		}
	} else if err == nil {
		utils.LogClassic("Found in cache video of id : " + vidId)
	} else {
		return nil, err
	}
	defer file.Close()

	return video, bot.AppendEltQueue(gId, &MusicInfos{Title: video.Title, Url: musicCacheFolderPath + vidId + ".m4a"})
}
*/

func buildYoutubeMusicPath(yId string) string {
	return musicCacheFolderPath + musicCacheYoutubeSubfolder + yId + ".m4a"
}

func ParseYoutubeId(input string) string {
	if strings.Contains(input, "youtube.com/watch?v=") { // Handle music.youtube.com links too
		input = input[strings.Index(input, "=")+1:]
		if strings.Contains(input, "&") {
			input = input[:strings.Index(input, "&")]
		}
	} else if strings.Contains(input, "youtube.com/playlist?list=") { // Handle music.youtube.com links too
		input = input[strings.Index(input, "=")+1:]
	} else if strings.Contains(input, "youtu.be/") {
		input = input[strings.Index(input, ".be/")+4:]
	}
	return input
}

func (bot *DiscordBot) DownloadMusicFromYoutube(client *youtube.Client, i *MusicInfos) (*os.File, error) {
	file, err := os.Create(buildMusicPath(i))
	if err != nil {
		return nil, err
	}
	// Download as file is mandatory since stream of more than 2m40s are ended without error thrown (probably because of youtube limitations)
	return bot.logicYoutubeDownload(file, client, i.Id)
}

func (bot *DiscordBot) logicYoutubeDownload(file *os.File, client *youtube.Client, yId string) (*os.File, error) {
	video, err := client.GetVideo(yId)
	if err != nil {
		return nil, err
	}
	format, err := ytinterface.GetBestAudioOnlyFormat(video.Formats)
	if err != nil {
		return nil, err
	}
	stream, _, err := client.GetStream(video, format)
	if err != nil {
		return nil, err
	}
	bot.Log("Music download start (" + yId + ")")
	_, err = io.Copy(file, stream)
	bot.Log("Music download end (" + yId + ")")
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (bot *DiscordBot) appendYoutubeVideoToQueue(gId string, v *youtube.Video) error {
	return bot.AppendEltQueue(gId, NewMusicInfos(v.ID, v.Title, []string{v.Author}, buildYoutubeMusicPath(v.ID), v.Duration, MusicInfosSources.Youtube))
}

func (bot *DiscordBot) appendYoutubePlaylistToQueue(gId string, p *youtube.Playlist, shuffle bool) error {
	s := make([]*MusicInfos, 0)
	for _, e := range p.Videos {
		s = append(s, NewMusicInfos(e.ID, e.Title, []string{e.Author}, buildYoutubeMusicPath(e.ID), e.Duration, MusicInfosSources.Youtube))
	}
	if shuffle {
		rand.Shuffle(len(s), func(i, j int) {
			s[i], s[j] = s[j], s[i]
		})
	}
	return bot.AppendEltsQueue(gId, s)
}

/* Wont implement for now
func GetYoutubeIdByVerbalQuery(query string) (string, error) {
	client := youtube.Client{}
	return "", nil
}
*/
