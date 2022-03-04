package logic

import (
	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	DiscordSession            *discordgo.Session
	Token                     string
	Id                        string
	LogChannelId              string
	Prefix                    string
	isMessageHandlerActivated bool
	commandParser             func(bot *DiscordBot, m *discordgo.MessageCreate) string
	commands                  map[string]func(bot *DiscordBot, m *discordgo.MessageCreate) error
	defaultCommand            func(bot *DiscordBot, m *discordgo.MessageCreate) error
}

/*
type Command struct {
	Name           string
	Execute        func(bot *DiscordBot, m *discordgo.MessageCreate) error
	PossibleErrors map[string]error
}
*/

func New(prefix string, token string, logChannelId string) *DiscordBot {
	return &DiscordBot{
		Prefix:                    prefix,
		Token:                     token,
		isMessageHandlerActivated: false,
		commands:                  make(map[string]func(bot *DiscordBot, m *discordgo.MessageCreate) error),
		defaultCommand:            func(bot *DiscordBot, m *discordgo.MessageCreate) error { return nil },
	}
}

func (bot *DiscordBot) Start() error {
	var err error
	bot.DiscordSession, err = discordgo.New("Bot " + bot.Token)
	if err != nil {
		return err
	}

	u, err := bot.DiscordSession.User("@me")
	if err != nil {
		return err
	}
	bot.Id = u.ID

	err = bot.DiscordSession.Open()
	if err != nil {
		return err
	}

	return nil
}

func (bot *DiscordBot) IsMessageHandlerActivated() bool {
	return bot.isMessageHandlerActivated
}
