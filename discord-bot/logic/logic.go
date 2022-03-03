package logic

import (
	"errors"
	"io"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

func (bot *DiscordBot) JoinChannel(guildID string, channelID string, mute bool, deaf bool) error {
	var err error
	bot.VoiceConnection, err = bot.DiscordSession.ChannelVoiceJoin(guildID, channelID, mute, deaf)
	if err != nil {
		bot.CommandError(err.Error(), nil)
		return err
	}
	return nil
}

func (bot *DiscordBot) PlayMusic() error {
	time.Sleep(250 * time.Millisecond)

	err := bot.VoiceConnection.Speaking(true)
	if err != nil {
		return err
	}

	time.Sleep(250 * time.Millisecond)

	err = DCA(bot.VoiceConnection, "playing.mp3")
	if err != nil {
		return err
	}

	time.Sleep(250 * time.Millisecond)

	err = bot.VoiceConnection.Speaking(false)
	if err != nil {
		return err
	}

	time.Sleep(250 * time.Millisecond)
	return nil
}

func PauseMusic() {

}

func /*(v *VoiceInstance)*/ DCA(vc *discordgo.VoiceConnection, url string) error {
	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 96
	opts.Application = "lowdelay"

	encodeSession, err := dca.EncodeFile(url, opts)
	if err != nil {
		return errors.New(" Failed creating an encoding session: " + err.Error())
	}
	//v.encoder = encodeSession
	done := make(chan error)
	/*stream := */ dca.NewStream(encodeSession, vc, done)

	//v.stream = stream
	for {
		select {
		case err := <-done:
			if err != nil && err != io.EOF {
				return errors.New("An error occured " + err.Error())
			}
			// Clean up incase something happened and ffmpeg is still running
			encodeSession.Cleanup()
			return nil
		}
	}
}
