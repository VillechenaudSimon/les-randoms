package logic

import (
	"io"
	"les-randoms/utils"
	"les-randoms/ytinterface"
	"os"
	"strings"

	"github.com/kkdai/youtube/v2"
)

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
	os.Remove("music.m4a")
	file, err := os.Create("music.m4a")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	utils.LogClassic("Downloading video of id : " + vidId)
	_, err = io.Copy(file, stream)
	if err != nil {
		return nil, err
	}

	return video, bot.AppendQueue(gId, &MusicInfos{Title: video.Title, Url: "music.m4a"})
}

func ParseYoutubeId(input string) string {
	if strings.Contains(input, "youtube.com/watch?v=") {
		input = input[strings.Index(input, "youtube.com/watch?v=")+20:]
		if strings.Contains(input, "&") {
			input = input[:strings.Index(input, "&")]
		}
	}
	return input
}
