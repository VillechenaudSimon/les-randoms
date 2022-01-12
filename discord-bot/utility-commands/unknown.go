package utilitycommands

import "github.com/bwmarrin/discordgo"

func CommandUnknown(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, _ = s.ChannelMessageSend(m.ChannelID, "You have entered an unknown command. You must be really smart <:YEP:930600118787383336>")
}
