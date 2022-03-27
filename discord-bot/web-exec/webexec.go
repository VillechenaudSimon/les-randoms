package webexec

import (
	"les-randoms/discord-bot/logic"
	"time"
)

var mainGuildId string
var bot *logic.DiscordBot

func init() {
	mainGuildId = "414802719673417729"
}

func Setup(b *logic.DiscordBot) {
	bot = b
}

func ExecuteMusicPlay(dsUserId string, input string) error {
	vs, err := bot.DiscordSession.State.VoiceState(mainGuildId, dsUserId)
	if err != nil {
		return err
	}

	if bot.DiscordSession.VoiceConnections[vs.GuildID] == nil { // If bot is not currently in a voice channel
		vc, err := bot.JoinChannel(vs.GuildID, vs.ChannelID, false, true)
		if err != nil {
			return err
		}

		bot.PlayQueue(vc)
	}

	_, err = bot.DownloadAndAppendQueue(vs.GuildID, logic.ParseYoutubeId(input))
	return err
}

func ExecuteMusicPause() error {
	return bot.PauseMusic(mainGuildId)
}

func ExecuteMusicResume() error {
	return bot.ResumeMusic(mainGuildId)
}

// True for paused, false for playing
func GetPlayStatus() bool {
	return bot.GetPlayStatus(mainGuildId)
}

func GetCurrentTime() time.Duration {
	return bot.GetCurrentTime(mainGuildId)
}

func GetCurrentTitle() string {
	return bot.GetCurrentTitle(mainGuildId)
}

func GetMusicQueue() []*logic.MusicInfos {
	return bot.GetMusicQueue(mainGuildId)
}
