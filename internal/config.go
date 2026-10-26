package internal

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DatabaseURL string
	BotToken	string
}

func NewEnv() (*Env, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	dbURL := os.Getenv("DATABASE_URL")
	botToken := os.Getenv("BOT_TOKEN")

	env := Env{
		DatabaseURL: dbURL,
		BotToken: botToken,
	}

	return &env, nil
}
