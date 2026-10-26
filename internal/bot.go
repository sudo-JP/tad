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

	session.AddHandler(handlers.HandleMessage)

	return &Bot{
		session: session,
		database: database,
	}, nil
}
