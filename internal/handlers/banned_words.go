package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/bwmarrin/discordgo"
	"github.com/sudo-JP/tad/internal/db"
)

// countBannedWords splits on anything that isn't a letter or digit, so "ass!" matches
// but "class" does not match "ass".
func countBannedWords(content string, bannedWords []string) int {
	banned := make(map[string]struct{}, len(bannedWords))
	for _, word := range bannedWords {
		banned[strings.ToLower(word)] = struct{}{}
	}

	words := strings.FieldsFunc(strings.ToLower(content), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	count := 0
	for _, word := range words {
		if _, ok := banned[word]; ok {
			count++
		}
	}
	return count
}

// checkBannedWord reports whether the message contained a banned word and was handled.
func checkBannedWord(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, database *db.TadDatabase) bool {
	watched, err := database.IsWatched(ctx, m.Author.ID)
	if err != nil {
		log.Printf("banned words: check watchlist for %s: %v", m.Author.ID, err)
		return false
	}
	if !watched {
		return false
	}

	bannedWords, err := database.GetBannedWords(ctx)
	if err != nil {
		log.Printf("banned words: load banned words: %v", err)
		return false
	}

	hits := countBannedWords(m.Content, bannedWords)
	if hits == 0 {
		return false
	}

	if err := database.IncrementLeetcodeOwed(ctx, m.Author.ID, hits); err != nil {
		log.Printf("banned words: increment owed for %s: %v", m.Author.ID, err)
		return true
	}

	msg := fmt.Sprintf("%s said a banned word, +%d leetcode owed", m.Author.Mention(), hits)
	if _, err := s.ChannelMessageSend(m.ChannelID, msg); err != nil {
		log.Printf("banned words: send message: %v", err)
	}
	return true
}
