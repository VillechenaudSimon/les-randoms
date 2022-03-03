package discordbot

import (
	"les-randoms/database"
	"les-randoms/utils"
	"les-randoms/webserver"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var appEnd chan bool // If a value is sent to this, the complete app stop

type DiscordBot struct {
	DiscordSession  *discordgo.Session
	Token           string
	Id              string
	VoiceConnection *discordgo.VoiceConnection
	LogChannelId    string
	Prefix          string
}

var Bot *DiscordBot

func Start(applicationEnd chan bool) {
	appEnd = applicationEnd

	var err error
	Bot.DiscordSession, err = discordgo.New("Bot " + Bot.Token)
	if err != nil {
		utils.LogError(err.Error())
		return
	}

	u, err := Bot.DiscordSession.User("@me")
	if err != nil {
		utils.LogError(err.Error())
		return
	}

	Bot.Id = u.ID

	Bot.DiscordSession.AddHandler(messageHandler)

	err = Bot.DiscordSession.Open()
	if err != nil {
		utils.LogError(err.Error())
		return
	}

	utils.LogSuccess("Discord bot successfully started")
}

func Close() {
	utils.LogNotNilError(Bot.DiscordSession.Close())
	utils.LogSuccess("Discord bot session successfully closed")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Bot musn't reply to it's own messages , to confirm it we perform this check
	if m.Author.ID == Bot.Id {
		return
	}

	//We test if the right prefix is detected
	if strings.HasPrefix(m.Content, Bot.Prefix) {
		Bot.Log("Message Read : " + m.Content)
		m.Content = m.Content[len(Bot.Prefix):]
	} else {
		return
	}

	detectAndCallCommand(s, m)
}

func detectAndCallCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch strings.Split(strings.ToUpper(m.Content), " ")[0] {
	case "SHUTDOWN":
		applicationShutdown(s, m)
	case "KANNA": //We ignore the MEE6 command that sends website url
		return
	case "PING":
		Bot.CommandPing(m)
	case "PLAY":
		Bot.CommandPlay(m)
	default:
		Bot.CommandUnknown(m)
	}
}

func applicationShutdown(s *discordgo.Session, m *discordgo.MessageCreate) {
	/*user, err := database.User_SelectFirst("WHERE discordid=" + utils.Esc(m.Author.ID))
	if err != nil {
		utilitycommands.CommandUnknown(s, m)
		utils.LogError(err.Error())
		return
	}*/
	if m.Author.ID != "178853941189148672" { // Discord id of Vemuni#4770
		Bot.CommandUnknown(m)
		return
	}
	_, _ = s.ChannelMessageSend(m.ChannelID, "Stopping webserver, database connection, discord bot and software..")
	webserver.StopServer()
	database.CloseDatabase()
	Close()
	appEnd <- true
}
