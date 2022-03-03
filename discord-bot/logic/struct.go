package logic

import "github.com/bwmarrin/discordgo"

type DiscordBot struct {
	DiscordSession  *discordgo.Session
	Token           string
	Id              string
	VoiceConnection *discordgo.VoiceConnection
	LogChannelId    string
	Prefix          string
}
