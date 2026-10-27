package handlers

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sudo-JP/tad/internal/db"
)

func HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate, db *db.TadDatabase) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	const prefix = "$"

	if !strings.HasPrefix(m.Content, prefix) {
		return
	}
	content := strings.TrimPrefix(m.Content, prefix)
	args := strings.Fields(content)

	handleArgs(s, m, db, args)
}

func handleArgs(s *discordgo.Session, m *discordgo.MessageCreate, db *db.TadDatabase, args []string) {
	if len(args) == 0 {
		return
	}
	commands, err := db.GetCommands(context.Background())
	if err != nil {
		return
	}

	for _, command := range commands {
		if handleCommand(s, m, args, &command) {
			break
		}
	}

}

func handleCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string, command *db.Command) bool {
	userCommand := args[0]
	if userCommand != command.Command {
		return false
	}
	switch command.Behavior {
	case db.BehaviorMedia:
		if command.Media != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, command.Media.Payload)
		}
	case db.BehaviorSpam:
		if command.Spam == nil {
			return true
		}
		count := 1
		if command.Spam.SpamCount != nil {
			count = *command.Spam.SpamCount
		}
		for range count {
			_, _ = s.ChannelMessageSend(m.ChannelID, command.Spam.Payload)
		}
	}

	return true
}
