package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrInvalidCommandBehavior = errors.New("invalid command behavior")

func validateSpamCount(count *int) error {
	if count != nil && *count <= 0 {
		return fmt.Errorf("spam count must be positive")
	}
	return nil
}

func (db *TadDatabase) CreateSpamCommand(ctx context.Context, command, payload string, spamCount *int) error {
	if err := validateSpamCount(spamCount); err != nil {
		return err
	}
	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, "INSERT INTO commands (command, behavior) VALUES (?, ?)", command, BehaviorSpam); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO spam_commands (command, payload, spam_count) VALUES (?, ?, ?)", command, payload, spamCount); err != nil {
		return err
	}
	return tx.Commit()
}

func (db *TadDatabase) CreateMediaCommand(ctx context.Context, command, payload string) error {
	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, "INSERT INTO commands (command, behavior) VALUES (?, ?)", command, BehaviorMedia); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO media_commands (command, payload) VALUES (?, ?)", command, payload); err != nil {
		return err
	}
	return tx.Commit()
}

func (db *TadDatabase) GetCommand(ctx context.Context, command string) (Command, error) {
	var result Command
	var behavior string
	err := db.pool.QueryRowContext(ctx, "SELECT command, behavior FROM commands WHERE command = ?", command).
		Scan(&result.Command, &behavior)
	if err != nil {
		return result, err
	}
	result.Behavior = CommandBehavior(behavior)
	switch result.Behavior {
	case BehaviorSpam:
		var item SpamCommand
		var count sql.NullInt64
		err = db.pool.QueryRowContext(ctx, "SELECT command, payload, spam_count FROM spam_commands WHERE command = ?", command).
			Scan(&item.Command, &item.Payload, &count)
		if err == nil && count.Valid {
			value := int(count.Int64)
			item.SpamCount = &value
		}
		result.Spam = &item
	case BehaviorMedia:
		item := MediaCommand{Command: command}
		err = db.pool.QueryRowContext(ctx, "SELECT payload FROM media_commands WHERE command = ?", command).Scan(&item.Payload)
		result.Media = &item
	default:
		return result, fmt.Errorf("%w: %q", ErrInvalidCommandBehavior, behavior)
	}
	return result, err
}

func (db *TadDatabase) GetCommands(ctx context.Context) ([]Command, error) {
	rows, err := db.pool.QueryContext(ctx, "SELECT command FROM commands")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var commands []Command
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		command, err := db.GetCommand(ctx, name)
		if err != nil {
			return nil, err
		}
		commands = append(commands, command)
	}
	return commands, rows.Err()
}

func (db *TadDatabase) UpdateSpamCommand(ctx context.Context, command, payload string, spamCount *int) error {
	if err := validateSpamCount(spamCount); err != nil {
		return err
	}
	_, err := db.pool.ExecContext(ctx,
		"UPDATE spam_commands SET payload = ?, spam_count = ? WHERE command = ?",
		payload, spamCount, command)
	return err
}

func (db *TadDatabase) UpdateMediaCommand(ctx context.Context, command, payload string) error {
	_, err := db.pool.ExecContext(ctx,
		"UPDATE media_commands SET payload = ? WHERE command = ?", payload, command)
	return err
}

func (db *TadDatabase) DeleteCommand(ctx context.Context, command string) error {
	_, err := db.pool.ExecContext(ctx, "DELETE FROM commands WHERE command = ?", command)
	return err
}
