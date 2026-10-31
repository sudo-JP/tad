package internal

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/sudo-JP/tad/internal/db"
	"github.com/sudo-JP/tad/internal/handlers"
)

type Bot struct {
	// ctx is the app's root context, cancelled on shutdown.
	ctx      context.Context
	session  *discordgo.Session
	database *db.TadDatabase
}

func NewBot(ctx context.Context, config *Env, database *db.TadDatabase) (*Bot, error) {
	session, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		return nil, err
	}
	bot := Bot{
		ctx:     ctx,
		session: session,
		database: database,
	}
	session.AddHandler(bot.handleMessage)
	return &bot, nil
}

func (b *Bot) Run() error {
	err := b.session.Open()
	return err
}

func (b *Bot) Close() error {
	return b.session.Close()
}

func (b *Bot) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
    handlers.HandleMessage(b.ctx, s, m, b.database)
}
