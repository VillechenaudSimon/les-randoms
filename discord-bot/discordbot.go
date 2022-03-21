package discordbot

import (
	"les-randoms/discord-bot/logic"
	webexec "les-randoms/discord-bot/web-exec"
	"les-randoms/utils"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var appEnd *chan bool // If a value is sent to this, the complete app stop

var Bot *logic.DiscordBot

func Start(applicationEnd *chan bool) {
	appEnd = applicationEnd

	err := Bot.Start()
	if err != nil {
		utils.LogError(err.Error())
		return
	}

	Bot.SetDefaultCommand(CommandUnknown)
	Bot.AddCommand("SHUTDOWN", CommandShutdown)
	Bot.AddCommand("PING", CommandPing)
	Bot.AddCommand("JOIN", CommandJoin)
	Bot.AddCommand("STOP", CommandDisconnect)
	Bot.AddCommand("TESTP", CommandTestPlay)
	Bot.AddCommand("PLAY", CommandPlay)
	Bot.AddCommand("PAUSE", CommandPause)
	Bot.AddCommand("RESUME", CommandResume)
	Bot.AddCommand("KANNA", func(bot *logic.DiscordBot, m *discordgo.MessageCreate) error {
		return nil
	})
	Bot.ActivateMessageHandler(func(bot *logic.DiscordBot, m *discordgo.MessageCreate) string {
		return strings.Split(strings.ToUpper(m.Content), " ")[0]
	})

	webexec.Setup(Bot)

	utils.LogSuccess("Discord bot successfully started")
}

func Close() {
	utils.LogNotNilError(Bot.DiscordSession.Close())
	utils.LogSuccess("Discord bot session successfully closed")
}
