package logic

import (
	"errors"
	"io"
	"les-randoms/utils"
	"les-randoms/ytinterface"
	"os"
	"strings"

	"github.com/kkdai/youtube/v2"
)

func (bot *DiscordBot) downloadIfNecesary(client *youtube.Client, i *MusicInfos) error {
	err := os.Mkdir(musicCacheFolderPath, os.ModeAppend)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	file, err := os.Open(buildVideoURL(i.Id))
	if errors.Is(err, os.ErrNotExist) {
		file, err = os.Create(buildVideoURL(i.Id))
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
		bot.Log("Downloading video of id : " + i.Id)
		stream, _, err := client.GetStream(video, format)
		if err != nil {
			return err
		}
		utils.LogClassic("Downloading video of id : " + i.Id)
		_, err = io.Copy(file, stream)
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

func ParseYoutubeId(input string) string {
	if strings.Contains(input, "youtube.com/watch?v=") {
		input = input[strings.Index(input, "=")+1:]
		if strings.Contains(input, "&") {
			input = input[:strings.Index(input, "&")]
		}
	} else if strings.Contains(input, "youtube.com/playlist?list=") {
		input = input[strings.Index(input, "=")+1:]
	}
	return input
}

/* Wont implement for now
func GetYoutubeIdByVerbalQuery(query string) (string, error) {
	client := youtube.Client{}
	return "", nil
}
*/
