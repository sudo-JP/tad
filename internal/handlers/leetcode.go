package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/sudo-JP/tad/internal/db"
)

func handleOwedLeetcode(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, database *db.TadDatabase, args []string) {
	if len(args) != 2 {
		return
	}
	user, err := database.GetUser(ctx, args[1])
	if err != nil {
		return
	}

	owed, err := database.GetLeetcodeOwed(ctx, args[1])
	if errors.Is(err, sql.ErrNoRows) {
		reply(s, m, user.Username + " isn't on leetcode list")
		return
	}
	if err != nil {
		return
	}
	reply(s, m, fmt.Sprintf("%s owes %d leetcode", user.Username, owed.OwedCount))
}

func handleFetchRandomLeetcode(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, database *db.TadDatabase, args []string) {
	http.NewRequestWithContext(ctx, "GET", "", nil)
}
