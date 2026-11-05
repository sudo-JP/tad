package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/sudo-JP/tad/internal/db"
)

// handleOwedLeetcode expects args without the leading "owed", e.g. ["<@123>"].
func handleOwedLeetcode(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, database *db.TadDatabase, args []string) {
	if len(args) != 1 {
		reply(s, m, "usage: $owed @user")
		return
	}
	user, err := mentionedUser(m)
	if err != nil {
		reply(s, m, err.Error())
		return
	}

	owed, err := database.GetLeetcodeOwed(ctx, user.ID)
	if errors.Is(err, sql.ErrNoRows) {
		reply(s, m, user.DisplayName()+" doesn't owe any leetcode")
		return
	}
	if err != nil {
		log.Printf("owed: get owed for %s: %v", user.ID, err)
		reply(s, m, "something went wrong")
		return
	}
	reply(s, m, fmt.Sprintf("%s owes %d leetcode", user.DisplayName(), owed.OwedCount))
}

func handleFetchRandomLeetcode(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, database *db.TadDatabase, args []string) {
	http.NewRequestWithContext(ctx, "GET", "", nil)
}
