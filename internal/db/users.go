package db

import (
	"context"
)

func (db *TadDatabase) CreateUser(ctx context.Context, userID string, username string) error {
	query := "INSERT INTO users (user_id, username) VALUES (?, ?)"
	_, err := db.pool.ExecContext(ctx, query, userID, username)
	return err
}

func (db *TadDatabase) GetUsernames(ctx context.Context) ([]string, error) {
	query := "SELECT username FROM users"
	rows, err := db.pool.Query(query)
	if err != nil {
		return []string{}, err
	}

	var usernames []string
	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			return []string{}, err
		}
		usernames = append(usernames, username)
	}

	if err = rows.Err(); err != nil {
		return []string{}, err
	}

	return usernames, err
}

func (db *TadDatabase) DeleteUser(ctx context.Context, userID string) error {
	query := "DELETE FROM users WHERE user_id = ?"
	_, err := db.pool.ExecContext(ctx, query, userID)
	return err
}

func (db *TadDatabase) UpdateUserName(ctx context.Context, userID string, newName string) error {
	query := "UPDATE users SET username = ? WHERE user_id = ?"
	_, err := db.pool.ExecContext(ctx, query, userID, newName)
	return err
}

func (db *TadDatabase) CreateSudoUser(ctx context.Context, userID string, username string) error {
	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Users query
	usersQuery := "INSERT INTO users (user_id, username) VALUES (?, ?)"
	_, err = tx.ExecContext(ctx, usersQuery, userID, username)

	sudoUsersQuery := "INSERT INTO sudo_users (user_id) VALUES (?)"
	_, err = tx.ExecContext(ctx, sudoUsersQuery, userID)

	return tx.Commit()
}
