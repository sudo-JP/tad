package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleMessage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Message.Author.ID == s.State.User.ID {
		return
	}
	content := strings.TrimPrefix(i.Message.Content, "")
	//args := 

}

func handlePrefix()
