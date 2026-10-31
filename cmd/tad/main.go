package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sudo-JP/tad/internal"
	"github.com/sudo-JP/tad/internal/db"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	env, err := internal.NewEnv()
	if err != nil {
		log.Fatalf("load env: %v", err)
	}
	if env.OwnerID == "" {
		log.Fatal("OWNER_ID is not set")
	}

	database, err := db.NewTadDB(env.DatabaseURL)
	if err != nil {
		log.Fatalf("open DB: %v", err)
	}
	defer database.Close()

	seedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	err = database.Seed(seedCtx, env.OwnerID)
	cancel()
	if err != nil {
		log.Fatalf("seed: %v", err)
	}

	bot, err := internal.NewBot(ctx, env, database)
	if err != nil {
		log.Fatalf("create bot: %v", err)
	}
	if err := bot.Run(); err != nil {
		log.Fatalf("run bot: %v", err)
	}
	defer bot.Close()

	log.Println("bot running, Ctrl+C to stop")
	<-ctx.Done()
}
