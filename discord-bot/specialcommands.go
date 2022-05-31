package discordbot

import (
	"les-randoms/database"
	"les-randoms/discord-bot/logic"
	"les-randoms/webserver"

	"github.com/bwmarrin/discordgo"
)

func CommandPing(bot *logic.DiscordBot, m *discordgo.MessageCreate) error {
	_, err := bot.DiscordSession.ChannelMessageSend(m.ChannelID, "Pong !")
	return err
}

func CommandShutdown(bot *logic.DiscordBot, m *discordgo.MessageCreate) error {
	/*user, err := database.User_SelectFirst("WHERE discordid=" + utils.Esc(m.Author.ID))
	if err != nil {
		utilitycommands.CommandUnknown(s, m)
		utils.LogError(err.Error())
		return
	}*/
	if m.Author.ID != "178853941189148672" { // Discord id of Vemuni#4770
		return bot.ExecuteDefaultCommand(m)
	}
	_, _ = bot.DiscordSession.ChannelMessageSend(m.ChannelID, "Stopping webserver, database connection, discord bot and software..")
	webserver.StopServer()
	database.CloseDatabase()
	Close()
	*appEnd <- true
	return nil
}

func CommandUnknown(bot *logic.DiscordBot, m *discordgo.MessageCreate) error {
	_, err := bot.DiscordSession.ChannelMessageSend(m.ChannelID, "You have entered an unknown command. You must be really smart <:YEP:930600118787383336>")
	return err
}
