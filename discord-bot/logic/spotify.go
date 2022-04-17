package logic

import (
	"errors"
	"os"
)

// TODO
func (bot *DiscordBot) DownloadMusicFromSpotify() (*os.File, error) {
	return nil, errors.New("not Implemented Yet")
}

func buildSpotifyMusicPath(sId string) string {
	return musicCacheFolderPath + musicCacheSpotifySubfolder + sId + ".m4a"
}

// TODO
func (bot *DiscordBot) appendSpotifyPlaylistToQueue(gId string, input string) error {
	return errors.New("not Implemented Yet")
}
