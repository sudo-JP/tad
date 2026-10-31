package internal

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DatabaseURL string
	BotToken	string
	OwnerID		string
}

func NewEnv() (*Env, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	dbURL := os.Getenv("DATABASE_URL")
	botToken := os.Getenv("BOT_TOKEN")
	ownerID := os.Getenv("OWNER_ID")

	env := Env{
		DatabaseURL: dbURL,
		BotToken: botToken,
		OwnerID: ownerID,
	}

	return &env, nil
}
