package internal

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sudo-JP/tad/internal/db"
	"github.com/sudo-JP/tad/internal/handlers"
)

type Bot struct {
	session  *discordgo.Session
	database *db.TadDatabase
}

func NewBot(config *Env, database *db.TadDatabase) (*Bot, error) {
	session, err := discordgo.New(config.BotToken)
	if err != nil {
		return nil, err
	}
	bot := Bot{
		session: session,
		database: database,
	}
	session.AddHandler(bot.handleMessage)
	return &bot, nil
}

func (b *Bot) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
    handlers.HandleMessage(s, m, b.database)
}
