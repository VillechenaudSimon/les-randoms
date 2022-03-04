package discordbot

import (
	"les-randoms/discord-bot/logic"

	"github.com/bwmarrin/discordgo"
)

func CommandPlay(bot *logic.DiscordBot, m *discordgo.MessageCreate) error {
	_, err := bot.DiscordSession.ChannelMessageSend(m.ChannelID, m.Author.Username+" asked me to play some music")
	if err != nil {
		return err
	}

	vc, err := bot.JoinMessageChannel(m, false, true)
	if err != nil {
		return err
	}

	err = bot.PlayMusic(vc)
	if err != nil {
		return err
	}

	return bot.DiscordSession.VoiceConnections[m.GuildID].Disconnect()
}

func CommandJoin(bot *logic.DiscordBot, m *discordgo.MessageCreate) error {
	_, err := bot.DiscordSession.ChannelMessageSend(m.ChannelID, m.Author.Username+" asked me to join the vocal channel")
	if err != nil {
		return err
	}

	_, err = bot.JoinMessageChannel(m, false, true)
	return err
}

func CommandDisconnect(bot *logic.DiscordBot, m *discordgo.MessageCreate) error {
	_, err := bot.DiscordSession.ChannelMessageSend(m.ChannelID, "Disconnecting from vocal channel..")
	if err != nil {
		return err
	}

	return bot.DiscordSession.VoiceConnections[m.GuildID].Disconnect()
}
