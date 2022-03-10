package webexec

import (
	"les-randoms/discord-bot/logic"
	"les-randoms/utils"
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

func ExecuteMusicPlay() error {
	return utils.NotImplementedYet
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
