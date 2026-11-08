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
	if after, found := strings.CutPrefix(m.Content, prefix); found {
		args = strings.Fields(after)
	}

	// Sudo skips the banned word check so "$sudo unban <word>" doesn't charge the sender.
	isSudo := len(args) > 0 && args[0] == "sudo"
	if !isSudo && checkBannedWord(ctx, s, m, db) {
		return
	}
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case m.Author.ID:
		handleOllamaMessage(ctx, s, m, args[1:])
	case "help":
		handleHelp(ctx, s, m, db)
	case "sudo": 
		handleSudo(ctx, s, m, db, args[1:])
	case "owed":
		handleOwedLeetcode(ctx, s, m, db, args[1:])
	case "leetcode":
		handleFetchRandomLeetcode(ctx, s, m, db, args[1:])
	default:
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
		handleSpam(s, m, args, command)
	}

	return true
}
