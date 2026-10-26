package db

import (
	"context"
	"errors"
	"fmt"
)

var ErrInvalidCommandBehavior = errors.New("invalid command behavior")

func validateBehavior(behavior string) error {
	if behavior != "spam" && behavior != "media" {
		return fmt.Errorf("%w: %q", ErrInvalidCommandBehavior, behavior)
	}
	return nil
}

func (db *TadDatabase) CreateCommand(ctx context.Context, command, behavior, payload string) error {
	if err := validateBehavior(behavior); err != nil {
		return err
	}
	_, err := db.pool.ExecContext(ctx,
		"INSERT INTO commands (command, behavior, payload) VALUES (?, ?, ?)",
		command, behavior, payload)
	return err
}

func (db *TadDatabase) GetCommand(ctx context.Context, command string) (Command, error) {
	var result Command
	err := db.pool.QueryRowContext(ctx,
		"SELECT command, behavior, payload FROM commands WHERE command = ?", command).
		Scan(&result.Command, &result.Behavior, &result.Payload)
	return result, err
}

func (db *TadDatabase) GetCommands(ctx context.Context) ([]Command, error) {
	rows, err := db.pool.QueryContext(ctx, "SELECT command, behavior, payload FROM commands")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []Command
	for rows.Next() {
		var command Command
		if err := rows.Scan(&command.Command, &command.Behavior, &command.Payload); err != nil {
			return nil, err
		}
		commands = append(commands, command)
	}
	return commands, rows.Err()
}

func (db *TadDatabase) UpdateCommand(ctx context.Context, command, behavior, payload string) error {
	if err := validateBehavior(behavior); err != nil {
		return err
	}
	_, err := db.pool.ExecContext(ctx,
		"UPDATE commands SET behavior = ?, payload = ? WHERE command = ?",
		behavior, payload, command)
	return err
}

func (db *TadDatabase) DeleteCommand(ctx context.Context, command string) error {
	_, err := db.pool.ExecContext(ctx, "DELETE FROM commands WHERE command = ?", command)
	return err
}
