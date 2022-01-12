package utilitycommands

import "github.com/bwmarrin/discordgo"

func CommandPing(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, _ = s.ChannelMessageSend(m.ChannelID, "Pong !")
}
