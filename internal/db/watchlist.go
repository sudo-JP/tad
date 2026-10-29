package db

import (
	"context"
	"database/sql"
	"errors"
)

// AddToWatchlist upserts the user (refreshing their username) and enrolls them.
// Adding an already watched user is a no-op.
func (db *TadDatabase) AddToWatchlist(ctx context.Context, userID, username string) error {
	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, `
		INSERT INTO users (user_id, username) VALUES (?, ?)
		ON CONFLICT(user_id) DO UPDATE SET username = excluded.username`,
		userID, username); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx,
		"INSERT OR IGNORE INTO banned_word_watchlist (user_id) VALUES (?)", userID); err != nil {
		return err
	}
	return tx.Commit()
}

func (db *TadDatabase) IsWatched(ctx context.Context, userID string) (bool, error) {
	var id string
	err := db.pool.QueryRowContext(ctx,
		"SELECT user_id FROM banned_word_watchlist WHERE user_id = ?", userID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (db *TadDatabase) GetWatchlist(ctx context.Context) ([]User, error) {
	rows, err := db.pool.QueryContext(ctx, `
		SELECT u.user_id, u.username
		FROM banned_word_watchlist w
		JOIN users u ON u.user_id = w.user_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UserID, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (db *TadDatabase) RemoveFromWatchlist(ctx context.Context, userID string) error {
	_, err := db.pool.ExecContext(ctx,
		"DELETE FROM banned_word_watchlist WHERE user_id = ?", userID)
	return err
}
