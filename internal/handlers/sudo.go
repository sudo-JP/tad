package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/bwmarrin/discordgo"
	"github.com/sudo-JP/tad/internal/db"
)

const sudoUsage = "```\n" +
	"$sudo ban|unban <word>\n" +
	"$sudo watch|unwatch @user\n" +
	"$sudo add|remove @user\n" +
	"$sudo new media <name> <url>\n" +
	"$sudo new spam <name> <count> <text...>\n" +
	"$sudo rm <name>\n" +
	"$sudo forgive @user <n>\n" +
	"```"

// reservedCommands are built-in names that DB commands must not shadow.
var reservedCommands = map[string]struct{}{"sudo": {}}

// sudoError is a message meant to be shown to the user as-is.
type sudoError string

func (e sudoError) Error() string { return string(e) }

var errUsage = sudoError(sudoUsage)

// handleSudo expects args without the leading "sudo", e.g. ["ban", "word"].
func handleSudo(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, database *db.TadDatabase, args []string) {
	isSudo, err := database.IsSudo(ctx, m.Author.ID)
	if err != nil {
		log.Printf("sudo: check sudo for %s: %v", m.Author.ID, err)
		return
	}
	if !isSudo {
		return
	}
	if len(args) == 0 {
		reply(s, m, sudoUsage)
		return
	}

	msg, err := runSudo(ctx, m, database, args[0], args[1:])
	if err != nil {
		msg = describeSudoError(args[0], err)
	}
	reply(s, m, msg)
}

func runSudo(ctx context.Context, m *discordgo.MessageCreate, database *db.TadDatabase, sub string, args []string) (string, error) {
	switch sub {
	case "ban", "unban":
		if len(args) != 1 {
			return "", errUsage
		}
		word := strings.ToLower(args[0])
		if !isSingleWord(word) {
			return "", sudoError("banned words can only contain letters and digits")
		}
		if sub == "ban" {
			return "banned `" + word + "`", database.CreateBannedWord(ctx, word)
		}
		return "unbanned `" + word + "`", database.DeleteBannedWord(ctx, word)

	case "watch", "unwatch", "add", "remove":
		if len(args) != 1 {
			return "", errUsage
		}
		user, err := mentionedUser(m)
		if err != nil {
			return "", err
		}
		switch sub {
		case "watch":
			return "watching " + user.Mention(), database.AddToWatchlist(ctx, user.ID, user.Username)
		case "unwatch":
			return "stopped watching " + user.Mention(), database.RemoveFromWatchlist(ctx, user.ID)
		case "add":
			return user.Mention() + " is now sudo", database.CreateSudoUser(ctx, user.ID, user.Username)
		default:
			if user.ID == m.Author.ID {
				return "", sudoError("you can't remove yourself")
			}
			return user.Mention() + " is no longer sudo", database.DeleteSudoUser(ctx, user.ID)
		}

	case "new":
		if len(args) == 0 {
			return "", errUsage
		}
		return runSudoNew(ctx, database, args[0], args[1:])

	case "rm":
		if len(args) != 1 {
			return "", errUsage
		}
		return "removed `" + args[0] + "`", database.DeleteCommand(ctx, args[0])

	case "forgive":
		if len(args) != 2 {
			return "", errUsage
		}
		user, err := mentionedUser(m)
		if err != nil {
			return "", err
		}
		n, err := strconv.Atoi(args[1])
		if err != nil || n <= 0 {
			return "", sudoError("amount must be a positive number")
		}
		return fmt.Sprintf("forgave %d leetcode for %s", n, user.Mention()),
			database.DecrementLeetcodeOwed(ctx, user.ID, n)
	}

	return "", errUsage
}

func runSudoNew(ctx context.Context, database *db.TadDatabase, kind string, args []string) (string, error) {
	switch kind {
	case "media":
		if len(args) != 2 {
			return "", errUsage
		}
		if err := checkCommandName(args[0]); err != nil {
			return "", err
		}
		return "added media command `" + args[0] + "`", database.CreateMediaCommand(ctx, args[0], args[1])

	case "spam":
		if len(args) < 3 {
			return "", errUsage
		}
		if err := checkCommandName(args[0]); err != nil {
			return "", err
		}
		count, err := strconv.Atoi(args[1])
		if err != nil || count <= 0 {
			return "", sudoError("count must be a positive number")
		}
		payload := strings.Join(args[2:], " ")
		return "added spam command `" + args[0] + "`", database.CreateSpamCommand(ctx, args[0], payload, &count)
	}

	return "", errUsage
}

func describeSudoError(sub string, err error) string {
	var userErr sudoError
	switch {
	case errors.As(err, &userErr):
		return userErr.Error()
	case errors.Is(err, db.ErrAlreadyExists):
		return "already exists"
	case errors.Is(err, db.ErrNotFound), errors.Is(err, sql.ErrNoRows):
		return "not found"
	case errors.Is(err, db.ErrOwedBelowZero):
		return "they don't owe that many"
	}
	log.Printf("sudo %s: %v", sub, err)
	return "something went wrong"
}

func mentionedUser(m *discordgo.MessageCreate) (*discordgo.User, error) {
	if len(m.Mentions) == 0 {
		return nil, sudoError("mention a user, e.g. @someone")
	}
	return m.Mentions[0], nil
}

func checkCommandName(name string) error {
	if _, ok := reservedCommands[strings.ToLower(name)]; ok {
		return sudoError("`" + name + "` is a built-in command")
	}
	return nil
}

// isSingleWord matches how countBannedWords tokenizes, so anything else could never trigger.
func isSingleWord(word string) bool {
	for _, r := range word {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return word != ""
}

func reply(s *discordgo.Session, m *discordgo.MessageCreate, msg string) {
	if _, err := s.ChannelMessageSend(m.ChannelID, msg); err != nil {
		log.Printf("send message: %v", err)
	}
}
