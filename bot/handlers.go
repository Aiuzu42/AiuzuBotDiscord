package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func CommandsHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "AiuzuBot" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Hello")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
