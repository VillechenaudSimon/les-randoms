package logic

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (bot *DiscordBot) HandleMessage(m *discordgo.MessageCreate) error {
	if !bot.IsMessageHandlerActivated() {
		return errors.New("Message handler not activated on this bot yet")
	}
	//Bot musn't reply to it's own messages , to confirm it we perform this check
	if m.Author.ID == bot.Id {
		return nil
	}
	//We test if the right prefix is detected
	if strings.HasPrefix(m.Content, bot.Prefix) {
		bot.Log("Message Read : " + m.Content)
		m.Content = m.Content[len(bot.Prefix):]
	} else {
		return nil
	}
	return bot.ExecuteCommand(bot.commandParser(bot, m), m)
}

func (bot *DiscordBot) ActivateMessageHandler(commandParser func(bot *DiscordBot, m *discordgo.MessageCreate) string) {
	bot.commandParser = commandParser

	bot.DiscordSession.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		err := bot.HandleMessage(m)
		if err != nil {
			bot.CommandError(err.Error(), m)
		}
	})

	bot.isMessageHandlerActivated = true
}

func (bot *DiscordBot) AddCommand(name string, f func(bot *DiscordBot, m *discordgo.MessageCreate) error) error {
	if _, ok := bot.commands[name]; ok {
		return errors.New("Command already exists : " + name)
	}
	bot.commands[name] = f
	return nil
}

func (bot *DiscordBot) SetDefaultCommand(f func(bot *DiscordBot, m *discordgo.MessageCreate) error) {
	bot.defaultCommand = f
}

func (bot *DiscordBot) ExecuteCommand(name string, m *discordgo.MessageCreate) error {
	if f, ok := bot.commands[name]; ok {
		return f(bot, m)
	} else {
		return bot.defaultCommand(bot, m)
	}
}

func (bot *DiscordBot) ExecuteDefaultCommand(m *discordgo.MessageCreate) error {
	return bot.defaultCommand(bot, m)
}
