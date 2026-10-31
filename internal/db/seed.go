package db

import "context"

// Seed makes the owner a sudo user. Safe to run on every startup.
// The user ID doubles as a placeholder username until the user is seen in Discord.
func (db *TadDatabase) Seed(ctx context.Context, ownerID string) error {
	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx,
		"INSERT INTO users (user_id, username) VALUES (?, ?) ON CONFLICT(user_id) DO NOTHING",
		ownerID, ownerID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx,
		"INSERT OR IGNORE INTO sudo_users (user_id) VALUES (?)", ownerID); err != nil {
		return err
	}
	return tx.Commit()
}
