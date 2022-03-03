package logic

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

func (bot *DiscordBot) HandleMessage(m *discordgo.MessageCreate) error {
	if !bot.IsMessageHandlerActivated() {
		return errors.New("Message handler not activated on this bot yet")
	}
	return bot.ExecuteCommand(bot.commandParser(bot, m), m)
}

func (bot *DiscordBot) ActivateMessageHandler(commandParser func(bot *DiscordBot, m *discordgo.MessageCreate) string) {
	bot.commandParser = commandParser
	bot.isMessageHandlerActivated = true
}

func (bot *DiscordBot) AddCommand(name string, f func(bot *DiscordBot, m *discordgo.MessageCreate) error) error {
	if _, ok := bot.commands[name]; ok {
		return errors.New("Command already exists : " + name)
	}
	bot.commands[name] = f
	return nil
}

func (bot *DiscordBot) ExecuteCommand(name string, m *discordgo.MessageCreate) error {
	return bot.commands[name](bot, m)
}
