package discordbot

import (
	"errors"
	"les-randoms/discord-bot/logic"
	"les-randoms/utils"

	"github.com/bwmarrin/discordgo"
)

var ErrorsPlay = struct {
	InputParsing    error
	NoInputDetected error
}{
	InputParsing:    errors.New("input parsing gone wrong"),
	NoInputDetected: errors.New("no input detected (ie nothing was asked to play)"),
}

func CommandPlay(bot *logic.DiscordBot, m *discordgo.MessageCreate) error {
	/*
		if bot.DiscordSession.VoiceConnections[m.GuildID] == nil { // If bot is not currently in a voice channel
			_, err := bot.JoinUserInChannel(m.GuildID, m.Author.ID, false, true)
			if err != nil {
				if errors.Is(err, discordgo.ErrStateNotFound) {
					_, err = bot.DiscordSession.ChannelMessageSend(m.ChannelID, m.Author.Username+" is not in a voice channel..")
					if err != nil {
						return err
					}
					return logic.ErrorsPlay.UserNotInVoiceChannel
				} else if err.Error() == "timeout waiting for voice" {
					_, err = bot.DiscordSession.ChannelMessageSend(m.ChannelID, "I can't join your voice channel")
					if err != nil {
						return err
					}
					return logic.ErrorsPlay.TimeoutWaitingForVoice
				} else {
					return err
				}
			}

		}

		if bot.GetMusicQueue(m.GuildID) == nil { // If the queue is not playing, start it
			bot.PlayQueue(bot.DiscordSession.VoiceConnections[m.GuildID])
		}

		args, err := logic.ParseArgs(m.Content)
		if err != nil {
			return err
		}

		if len(args.Params) <= 0 {
			_, err := bot.DiscordSession.ChannelMessageSend(m.ChannelID, "Type "+bot.Prefix+"play [Music] in order to play the desired song.")
			return err
		}

		return bot.AppendQueueFromInput(m.GuildID, args.Params[0])*/

	args, err := logic.ParseArgs(m.Content)
	if err != nil {
		return ErrorsPlay.InputParsing
	}

	if len(args.Params) <= 0 {
		_, err = bot.DiscordSession.ChannelMessageSend(m.ChannelID, "Type "+bot.Prefix+"play [Music] in order to play the desired song.")
		utils.LogNotNilError(err)
		return ErrorsPlay.NoInputDetected
	}

	err = bot.PlayOrder(m.GuildID, m.Author.ID, m.Content)
	if err != nil {
		_, err = bot.DiscordSession.ChannelMessageSend(m.ChannelID, "Error while executing your order : "+err.Error())
		utils.LogNotNilError(err)
	}
	return err
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

	_, err = bot.JoinUserInChannel(m.GuildID, m.Author.ID, false, true)
	return err
}

func CommandDisconnect(bot *logic.DiscordBot, m *discordgo.MessageCreate) error {
	_, err := bot.DiscordSession.ChannelMessageSend(m.ChannelID, "Disconnecting from vocal channel..")
	if err != nil {
		return err
	}

	return bot.Disconnect(m.GuildID)
}
