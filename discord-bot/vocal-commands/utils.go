package vocalcommands

import (
	"io"
	"les-randoms/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

func /*(v *VoiceInstance)*/ DCA(vc *discordgo.VoiceConnection, url string) {
	utils.LogDebug("CALL DCA : " + url)
	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 96
	opts.Application = "lowdelay"

	encodeSession, err := dca.EncodeFile(url, opts)
	if err != nil {
		utils.LogError("FATA: Failed creating an encoding session: " + err.Error())
	}
	//v.encoder = encodeSession
	done := make(chan error)
	/*stream := */ dca.NewStream(encodeSession, vc, done)

	//v.stream = stream
	for {
		select {
		case err := <-done:
			if err != nil && err != io.EOF {
				utils.LogError("FATA: An error occured " + err.Error())
			}
			// Clean up incase something happened and ffmpeg is still running
			encodeSession.Cleanup()
			utils.LogDebug("END DCA")
			return
		}
	}
}
