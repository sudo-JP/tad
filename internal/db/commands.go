package db

import "context"

func (db *TadDatabase) CreateCommand(ctx context.Context, command string, action string, payload string) error {
	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Command query
	commandQuery := "INSERT INTO commands (command, action) VALUES (?, ?)"
	_, err = tx.ExecContext(ctx, commandQuery, command, action)

	actionQuery := "INSERT INTO action (action, payload) VALUES (?, ?)"
	_, err = tx.ExecContext(ctx, actionQuery, action, payload)

	return tx.Commit()
}
