package discordbot

import (
	"fmt"
	"io"
	"les-randoms/discord-bot/logic"
	"les-randoms/utils"
	"les-randoms/ytinterface"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/kkdai/youtube/v2"
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
		args, err := parseArgs(m.Content)
		if err != nil {
			return err
		}
		for k, v := range args.Params {
			utils.LogDebug(fmt.Sprint(k) + " -> " + v)
		}
		for k, v := range args.Options {
			utils.LogDebug(k + " -> " + v)
		}

		if len(args.Params) <= 0 {
			_, err := bot.DiscordSession.ChannelMessageSend(m.ChannelID, "Wrong usage.")
			return err
		}

		client := youtube.Client{}
		video, err := client.GetVideo(args.Params[0])
		if err != nil {
			return err
		}
		format, err := ytinterface.GetBestAudioOnlyFormat(video.Formats)
		if err != nil {
			return err
		}
		// Download as file is mandatory since stream of more than 2m40s are ended without error thrown (probably because of youtube limitations)
		stream, _, err := client.GetStream(video, format)
		if err != nil {
			return err
		}
		os.Remove("music.m4a")
		file, err := os.Create("music.m4a")
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, stream)
		if err != nil {
			return err
		}

		return bot.AppendQueue(m.GuildID, &logic.MusicInfos{Title: video.Title, Url: "music.m4a"})
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
