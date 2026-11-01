package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sudo-JP/tad/internal/db"
)

func HandleMessage(parent context.Context, s *discordgo.Session, m *discordgo.MessageCreate, db *db.TadDatabase) {
	if m.Author.Bot {
		return
	}

	ctx, cancel := context.WithTimeout(parent, 5 * time.Second)
	defer cancel()

	const prefix = "$"
	var args []string
	if strings.HasPrefix(m.Content, prefix) {
		args = strings.Fields(strings.TrimPrefix(m.Content, prefix))
	}

	// Sudo runs before the banned word check so "$sudo unban <word>" doesn't charge the sender.
	switch args[0] {
	case "sudo": 
		handleSudo(ctx, s, m, db, args[1:])
	case "owed":
		handleOwedLeetcode()
	case "leetcode":
		handleFetchRandomLeetcode()
	default:
		if checkBannedWord(ctx, s, m, db) {
			return
		}
		handleArgs(s, m, db, args, ctx)
	}
}

func handleArgs(s *discordgo.Session, m *discordgo.MessageCreate, db *db.TadDatabase, args []string, ctx context.Context) {
	if len(args) == 0 {
		return
	}
	command, err := db.GetCommand(ctx, args[0])
	if errors.Is(err, sql.ErrNoRows) {
		return
	}
	if err != nil {
		log.Printf("get command %q: %v", args[0], err)
		return
	}
	handleCommand(s, m, args, &command)
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
