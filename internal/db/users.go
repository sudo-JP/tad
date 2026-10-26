package db

import "context"

func (db *TadDatabase) CreateUser(ctx context.Context, userID, username string) error {
	_, err := db.pool.ExecContext(ctx, "INSERT INTO users (user_id, username) VALUES (?, ?)", userID, username)
	return err
}

func (db *TadDatabase) GetUser(ctx context.Context, userID string) (User, error) {
	var user User
	err := db.pool.QueryRowContext(ctx, "SELECT user_id, username FROM users WHERE user_id = ?", userID).
		Scan(&user.UserID, &user.Username)
	return user, err
}

func (db *TadDatabase) GetUsers(ctx context.Context) ([]User, error) {
	rows, err := db.pool.QueryContext(ctx, "SELECT user_id, username FROM users")
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

func (db *TadDatabase) GetUsernames(ctx context.Context) ([]string, error) {
	users, err := db.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	usernames := make([]string, 0, len(users))
	for _, user := range users {
		usernames = append(usernames, user.Username)
	}
	return usernames, nil
}

func (db *TadDatabase) UpdateUserName(ctx context.Context, userID, username string) error {
	_, err := db.pool.ExecContext(ctx, "UPDATE users SET username = ? WHERE user_id = ?", username, userID)
	return err
}

func (db *TadDatabase) DeleteUser(ctx context.Context, userID string) error {
	_, err := db.pool.ExecContext(ctx, "DELETE FROM users WHERE user_id = ?", userID)
	return err
}

func (db *TadDatabase) CreateSudoUser(ctx context.Context, userID, username string) error {
	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, "INSERT INTO users (user_id, username) VALUES (?, ?)", userID, username); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO sudo_users (user_id) VALUES (?)", userID); err != nil {
		return err
	}
	return tx.Commit()
}

func (db *TadDatabase) GetSudoUser(ctx context.Context, userID string) (SudoUser, error) {
	var user SudoUser
	err := db.pool.QueryRowContext(ctx, "SELECT user_id FROM sudo_users WHERE user_id = ?", userID).Scan(&user.UserID)
	return user, err
}

func (db *TadDatabase) GetSudoUsers(ctx context.Context) ([]SudoUser, error) {
	rows, err := db.pool.QueryContext(ctx, "SELECT user_id FROM sudo_users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []SudoUser
	for rows.Next() {
		var user SudoUser
		if err := rows.Scan(&user.UserID); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (db *TadDatabase) DeleteSudoUser(ctx context.Context, userID string) error {
	_, err := db.pool.ExecContext(ctx, "DELETE FROM sudo_users WHERE user_id = ?", userID)
	return err
}
