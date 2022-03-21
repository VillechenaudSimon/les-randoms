package logic

import (
	"fmt"
	"io"
	"les-randoms/utils"
	"os"

	"github.com/kkdai/youtube/v2"
)

func (bot *DiscordBot) GetYoutubeVideoFromId(id string) (*MusicInfos, error) {
	i := &MusicInfos{}

	client := youtube.Client{}

	video, err := client.GetVideo(id)
	if err != nil {
		return nil, err
	}

	utils.LogDebug(fmt.Sprint(video.Formats.FindByItag(139)))
	utils.LogDebug(fmt.Sprint(video.Formats.FindByItag(140)))
	utils.LogDebug(fmt.Sprint(video.Formats.FindByItag(141)))
	stream, _, err := client.GetStream(video, video.Formats.FindByItag(140))
	if err != nil {
		return nil, err
	}

	file, err := os.Create("music.m4a")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return nil, err
	}

	i.Title = "DEFAULT"
	i.Url, err = client.GetStreamURL(video, video.Formats.FindByItag(140))
	if err != nil {
		return nil, err
	}
	return i, nil
	//client := &http.Client{
	//	Transport: &transport.APIKey{Key: youtubeToken},
	//}

	//service, err := youtube.New(client)
	//if err != nil {
	//log.Fatalf("Error creating new YouTube client: %v", err)
	//	return
	//}
}
