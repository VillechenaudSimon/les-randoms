package discordbot

import (
	"les-randoms/discord-bot/logic"

	"github.com/bwmarrin/discordgo"
)

func CommandPlay(bot *logic.DiscordBot, m *discordgo.MessageCreate) error {
	_, err := bot.DiscordSession.ChannelMessageSend(m.ChannelID, m.Author.Username+" asked me to join the vocal channel")
	if err != nil {
		return err
	}

	voiceState, err := bot.DiscordSession.State.VoiceState(m.GuildID, m.Author.ID)
	if err != nil {
		return err
	}

	err = bot.JoinChannel(m.GuildID, voiceState.ChannelID, false, true)
	if err != nil {
		return err
	}

	err = bot.PlayMusic()
	if err != nil {
		return err
	}

	err = bot.VoiceConnection.Disconnect()
	if err != nil {
		return err
	}

	return nil
}
