package main

import (
	"log"

	"github.com/sudo-JP/tad/internal"
	"github.com/sudo-JP/tad/internal/db"
)

func main()  {
	env, err := internal.NewEnv()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	db, err := db.NewTadDB(env.DatabaseURL)
	if err != nil {
		log.Fatal("Error create DB")
	}

	bot, err := internal.NewBot(env, db)
	if err != nil {
		log.Fatal("Error create bot")
	}

	if err := bot.Run(); err != nil {
		log.Fatal("Error running bot")
	}

	if err := bot.CloseOnTermination(); err != nil {
		log.Fatal("Error closing bot")
	}
}


