package bot

import (
	utilitycommands "les-randoms/discord-bot/utility-commands"
	botutils "les-randoms/discord-bot/utils"
	vocalcommands "les-randoms/discord-bot/vocal-commands"
	"les-randoms/utils"
	"strings"

	"os"

	"github.com/bwmarrin/discordgo"
)

var goBot *discordgo.Session
var BotToken string
var BotId string

func init() {
	BotToken = os.Getenv("DISCORD_BOT_TOKEN")
}

func Start() {
	goBot, err := discordgo.New("Bot " + BotToken)
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

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Bot musn't reply to it's own messages , to confirm it we perform this check
	if m.Author.ID == BotId {
		return
	}

	//We test if the right prefix is detected
	if strings.HasPrefix(m.Content, Prefix) {
		botutils.BotLog("Message Read : " + m.Content)
		m.Content = strings.TrimLeft(m.Content, Prefix)
	} else {
		return
	}

	detectAndCallCommand(s, m)
}

func detectAndCallCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch strings.Split(strings.ToUpper(m.Content), " ")[0] {
	case "PING":
		utilitycommands.CommandPing(s, m)
	case "PLAY":
		vocalcommands.CommandPlay(s, m)
	default:
		utilitycommands.CommandUnknown(s, m)
	}
}
