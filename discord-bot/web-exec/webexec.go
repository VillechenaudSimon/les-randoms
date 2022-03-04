package webexec

import (
	"errors"
	"les-randoms/discord-bot/logic"
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
	return bot.PlayMusic(bot.DiscordSession.VoiceConnections[mainGuildId])
}

func ExecuteMusicPause() error {
	return errors.New("Not implemented yet")
}
