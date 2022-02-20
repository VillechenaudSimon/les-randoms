package bot

import (
	"les-randoms/database"
	utilitycommands "les-randoms/discord-bot/utility-commands"
	botutils "les-randoms/discord-bot/utils"
	vocalcommands "les-randoms/discord-bot/vocal-commands"
	"les-randoms/utils"
	"les-randoms/webserver"
	"strings"

	"os"

	"github.com/bwmarrin/discordgo"
)

var goBot *discordgo.Session
var BotToken string
var BotId string
var appEnd chan bool // If a value is sent to this, the complete app stop

func init() {
	BotToken = os.Getenv("DISCORD_BOT_TOKEN")
}

func Start(applicationEnd chan bool) {
	appEnd = applicationEnd

	var err error
	goBot, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		utils.LogError(err.Error())
		return
	}

	u, err := goBot.User("@me")
	if err != nil {
		utils.LogError(err.Error())
		return
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	if err != nil {
		utils.LogError(err.Error())
		return
	}

	utils.LogSuccess("Discord bot successfully started")
}

func Close() {
	utils.LogNotNilError(goBot.Close())
	utils.LogSuccess("Discord bot session successfully closed")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Bot musn't reply to it's own messages , to confirm it we perform this check
	if m.Author.ID == BotId {
		return
	}

	//We test if the right prefix is detected
	if strings.HasPrefix(m.Content, Prefix) {
		botutils.BotLog("Message Read : " + m.Content)
		m.Content = m.Content[len(Prefix):]
	} else {
		return
	}

	detectAndCallCommand(s, m)

	s.Close()
}

func detectAndCallCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch strings.Split(strings.ToUpper(m.Content), " ")[0] {
	case "SHUTDOWN":
		applicationShutdown(s, m)
	case "KANNA": //We ignore the MEE6 command that sends website url
		return
	case "PING":
		utilitycommands.CommandPing(s, m)
	case "PLAY":
		vocalcommands.CommandPlay(s, m)
	default:
		utilitycommands.CommandUnknown(s, m)
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
		utilitycommands.CommandUnknown(s, m)
		return
	}
	_, _ = s.ChannelMessageSend(m.ChannelID, "Stopping webserver, database connection, discord bot and software..")
	webserver.StopServer()
	database.CloseDatabase()
	Close()
	appEnd <- true
}
