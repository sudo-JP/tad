package db

import (
	"context"
)

func (db *TadDatabase) CreateUser(ctx context.Context, userID string, username string) error {
	query := "INSERT INTO users (user_id, username) VALUES (?, ?);"

	_, err := db.pool.ExecContext(ctx, query, userID, username)

	return err
}

