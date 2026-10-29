package handlers

import (
	"context"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sudo-JP/tad/internal/db"
)

func checkBannedWord(s *discordgo.Session, m *discordgo.MessageCreate, db *db.TadDatabase, ctx context.Context) bool {
	bannedWords, err := db.GetBannedWords(ctx)
	if err != nil {
		return false
	}

	args := strings.Fields(m.Content)
	if len(args) != 1 {
		return false
	}

	word := args[0]
	if slices.Contains(bannedWords, word) {
		db.GetLeetcodeOwedForUser(ctx, m.Author.ID)
	}

	return false
}
