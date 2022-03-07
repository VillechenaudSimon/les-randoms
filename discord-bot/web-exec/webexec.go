package webexec

import (
	"les-randoms/discord-bot/logic"
	"les-randoms/utils"
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
	return bot.PauseMusic(bot.DiscordSession.VoiceConnections[mainGuildId])
}

func ExecuteMusicResume() error {
	return bot.ResumeMusic(bot.DiscordSession.VoiceConnections[mainGuildId])
}
