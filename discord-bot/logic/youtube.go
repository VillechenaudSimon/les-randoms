package logic

import (
	"os"
)

var youtubeToken string

func init() {
	youtubeToken = os.Getenv("YOUTUBE_API_KEY")
}

func getYoutubeVideo(url string) {
	//client := &http.Client{
	//	Transport: &transport.APIKey{Key: youtubeToken},
	//}

	//service, err := youtube.New(client)
	//if err != nil {
	//log.Fatalf("Error creating new YouTube client: %v", err)
	//	return
	//}
}
