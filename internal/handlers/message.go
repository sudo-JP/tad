package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sudo-JP/tad/internal/db"
)

func HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate, db *db.TadDatabase) {
	if m.Author.Bot {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	const prefix = "$"
	var args []string
	if strings.HasPrefix(m.Content, prefix) {
		args = strings.Fields(strings.TrimPrefix(m.Content, prefix))
	}

	// Sudo runs before the banned word check so "$sudo unban <word>" doesn't charge the sender.
	if len(args) > 0 && args[0] == "sudo" {
		handleSudo(ctx, s, m, db, args[1:])
		return
	}

	if checkBannedWord(ctx, s, m, db) {
		return
	}

	handleArgs(s, m, db, args, ctx)
}

func handleArgs(s *discordgo.Session, m *discordgo.MessageCreate, db *db.TadDatabase, args []string, ctx context.Context) {
	if len(args) == 0 {
		return
	}
	commands, err := db.GetCommands(ctx)
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
