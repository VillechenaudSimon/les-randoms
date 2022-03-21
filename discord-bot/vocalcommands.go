package discordbot

import (
	"les-randoms/discord-bot/logic"

	"github.com/bwmarrin/discordgo"
)

func CommandTestPlay(bot *logic.DiscordBot, m *discordgo.MessageCreate) error {
	_, err := bot.DiscordSession.ChannelMessageSend(m.ChannelID, m.Author.Username+" asked me to play some music in order to test smth..")
	if err != nil {
		return err
	}

	vc, err := bot.JoinMessageChannel(m, false, true)
	if err != nil {
		return err
	}

	err = bot.TestPlayMusic(vc)
	if err != nil {
		return err
	}

	return bot.Disconnect(m.GuildID)
}

func CommandPlay(bot *logic.DiscordBot, m *discordgo.MessageCreate) error {
	if bot.DiscordSession.VoiceConnections[m.GuildID] == nil { // If bot is not currently in a voice channel
		_, err := bot.DiscordSession.ChannelMessageSend(m.ChannelID, m.Author.Username+" asked me to play some music..")
		if err != nil {
			return err
		}

		vc, err := bot.JoinMessageChannel(m, false, true)
		if err != nil {
			return err
		}

		return bot.PlayQueue(vc)
	} else {
		//return bot.AppendQueue(m.GuildID, &logic.MusicInfos{Title: "TestQueue", Url: "playing.mp3"})
		i, err := bot.GetYoutubeVideoFromId("iJqAeMHnmY0")
		if err != nil {
			return nil
		}
		return bot.AppendQueue(m.GuildID, i)
	}
}

func CommandPause(bot *logic.DiscordBot, m *discordgo.MessageCreate) error {
	_, err := bot.DiscordSession.ChannelMessageSend(m.ChannelID, "Pausing music..")
	if err != nil {
		return err
	}

	return bot.PauseMusic(m.GuildID)
}

func CommandResume(bot *logic.DiscordBot, m *discordgo.MessageCreate) error {
	_, err := bot.DiscordSession.ChannelMessageSend(m.ChannelID, "Resuming music..")
	if err != nil {
		return err
	}

	return bot.ResumeMusic(m.GuildID)
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

	return bot.Disconnect(m.GuildID)
}
