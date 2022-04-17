package logic

import (
	"context"
	"errors"
	"fmt"
	"les-randoms/utils"
	"math/rand"
	"os"
	"strings"

	"github.com/kkdai/youtube/v2"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var spotifyCtx context.Context
var spotifyConfig *clientcredentials.Config
var spotifyToken *oauth2.Token

func init() {
	spotifyCtx = context.Background()
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}

	var err error
	spotifyToken, err = config.Token(spotifyCtx)
	if err != nil {
		utils.LogError("couldn't get spotify token: " + err.Error())
	}
}

// Yeah a youtube client for this
// Because spotify tracks are downloaded through a youtube search
func (bot *DiscordBot) DownloadMusicFromSpotify(client *youtube.Client, i *MusicInfos) (*os.File, error) {
	file, err := os.Create(buildMusicPath(i))
	if err != nil {
		return nil, err
	}
	as := ""
	for _, a := range i.Artists {
		as = as + a + " "
	}
	bot.Log("Searching youtube for : " + as + i.Title + " (Duration: " + fmt.Sprint(int(i.Duration.Seconds())) + ")")
	yId, err := GetYoutubeId(as+i.Title, int(i.Duration.Seconds()))
	if err != nil {
		return nil, err
	}
	return bot.logicYoutubeDownload(file, client, yId)
}

func buildSpotifyMusicPath(sId string) string {
	return musicCacheFolderPath + musicCacheSpotifySubfolder + sId + ".m4a"
}

func (bot *DiscordBot) appendSpotifyPlaylistToQueue(gId string, input string, shuffle bool) error {
	client := spotify.New(spotifyauth.New().Client(spotifyCtx, spotifyToken))
	id, err := GetSpotifyPlaylistId(input)
	if err != nil {
		return err
	}
	tracks, err := client.GetPlaylistTracks(spotifyCtx, spotify.ID(id))
	if err != nil {
		return err
	}
	count := 0
	s := make([]*MusicInfos, 0)
	for page := 1; ; page++ {
		for _, t := range tracks.Tracks {
			count++
			artists := make([]string, 0)
			for _, a := range t.Track.Artists {
				artists = append(artists, a.Name)
			}
			s = append(s, NewMusicInfos(string(t.Track.ID), t.Track.Name, artists, buildSpotifyMusicPath(string(t.Track.ID)), t.Track.TimeDuration(), MusicInfosSources.Spotify))
		}
		err = client.NextPage(spotifyCtx, tracks)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			return err
		}
	}
	if shuffle {
		rand.Shuffle(len(s), func(i, j int) {
			s[i], s[j] = s[j], s[i]
		})
	}
	return bot.AppendEltsQueue(gId, s)
}

func GetSpotifyPlaylistId(s string) (string, error) {
	if strings.Contains(s, "spotify.com/playlist/") {
		id := s[strings.Index(s, "playlist/")+9:]
		if strings.Contains(id, "?") {
			id = id[:strings.Index(id, "?")]
		}
		return id, nil
	} else {
		return "", errors.New("can't parse given spotify url to find playlist id")
	}
}
