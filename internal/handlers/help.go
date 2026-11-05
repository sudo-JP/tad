package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sudo-JP/tad/internal/db"
)

const builtinHelp = "$help              this menu\n" +
	"$owed @user        show leetcode owed\n" +
	"$leetcode          random leetcode problem\n" +
	"$sudo help         admin commands (sudo only)"

// Discord rejects embed field values over 1024 characters.
const maxFieldLen = 1024

func handleHelp(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, database *db.TadDatabase) {
	commands, err := database.GetCommands(ctx)
	if err != nil {
		log.Printf("help: get commands: %v", err)
		reply(s, m, "something went wrong")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Help",
		Description: "Prefix: `$`",
		Color:       0xFFFFFF,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Built-in", Value: codeBlock(builtinHelp)},
			{Name: "Custom", Value: codeBlock(customHelp(commands))},
		},
	}
	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, embed); err != nil {
		log.Printf("help: send embed: %v", err)
	}
}

func customHelp(commands []db.Command) string {
	if len(commands) == 0 {
		return "none yet, ask a sudo user to add some"
	}

	// Leave room for the code block fences and the "..." line.
	limit := maxFieldLen - len("```\n\n```") - len("\n...")
	var b strings.Builder
	for _, command := range commands {
		line := fmt.Sprintf("$%-17s %s\n", command.Command, describeCommand(command))
		if b.Len()+len(line) > limit {
			b.WriteString("...")
			break
		}
		b.WriteString(line)
	}
	return strings.TrimRight(b.String(), "\n")
}

func describeCommand(command db.Command) string {
	if command.Behavior == db.BehaviorSpam && command.Spam != nil && command.Spam.SpamCount != nil {
		return fmt.Sprintf("spam x%d", *command.Spam.SpamCount)
	}
	return string(command.Behavior)
}

func codeBlock(text string) string {
	return "```\n" + text + "\n```"
}
