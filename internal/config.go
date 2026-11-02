package internal

import (
	"errors"
	"io/fs"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DatabaseURL string
	BotToken	string
	OwnerID		string
}

func NewEnv() (*Env, error) {
	// .env is optional; in Docker the variables come from compose instead.
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
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
