package handlers

import (
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sudo-JP/tad/internal/db"
)

func handleSpam(s *discordgo.Session, m *discordgo.MessageCreate, args []string, command *db.Command) {
	if command.Spam == nil {
		return
	}
	count := 0
	if command.Spam.SpamCount != nil {
		count = *command.Spam.SpamCount
	}
	isRandom := strings.Contains(command.Spam.Payload, "{random}")
	isArgsable := strings.Contains(command.Spam.Payload, "{args}")
	text := strings.ReplaceAll(command.Spam.Payload, "{random}", "")
	text = strings.ReplaceAll(text, "{args}", strings.Join(args[1:], " "))

	if isRandom && isArgsable {
		reply(s, m, "Too many configs for random and args-able")
		return
	} else if isRandom {
		for range count {
			randomIdx := rand.Intn(len(args) - 1)
			reply(s, m, args[randomIdx + 1])
		}	
	} else {
		for range count {
			reply(s, m, text)
		}

	}
}
