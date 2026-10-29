package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrOwedBelowZero = errors.New("cannot decrement owed count below zero")

func (db *TadDatabase) GetLeetcodeOwed(ctx context.Context, userID string) (LeetcodeOwed, error) {
	var owed LeetcodeOwed
	err := db.pool.QueryRowContext(ctx,
		"SELECT user_id, owed_count FROM leetcode_owed WHERE user_id = ?", userID).
		Scan(&owed.UserID, &owed.OwedCount)
	return owed, err
}

func (db *TadDatabase) GetAllLeetcodeOwed(ctx context.Context) ([]LeetcodeOwed, error) {
	rows, err := db.pool.QueryContext(ctx, "SELECT user_id, owed_count FROM leetcode_owed")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var owed []LeetcodeOwed
	for rows.Next() {
		var item LeetcodeOwed
		if err := rows.Scan(&item.UserID, &item.OwedCount); err != nil {
			return nil, err
		}
		owed = append(owed, item)
	}
	return owed, rows.Err()
}

func (db *TadDatabase) SetLeetcodeOwed(ctx context.Context, userID string, count int) error {
	if count < 0 {
		return fmt.Errorf("owed count cannot be negative")
	}
	_, err := db.pool.ExecContext(ctx, `
		INSERT INTO leetcode_owed (user_id, owed_count) VALUES (?, ?)
		ON CONFLICT(user_id) DO UPDATE SET owed_count = excluded.owed_count`,
		userID, count)
	return err
}

func (db *TadDatabase) IncrementLeetcodeOwed(ctx context.Context, userID string, amount int) error {
	if amount < 0 {
		return fmt.Errorf("increment amount cannot be negative")
	}
	_, err := db.pool.ExecContext(ctx, `
		INSERT INTO leetcode_owed (user_id, owed_count) VALUES (?, ?)
		ON CONFLICT(user_id) DO UPDATE SET owed_count = owed_count + excluded.owed_count`,
		userID, amount)
	return err
}

func (db *TadDatabase) DecrementLeetcodeOwed(ctx context.Context, userID string, amount int) error {
	if amount < 0 {
		return fmt.Errorf("decrement amount cannot be negative")
	}
	result, err := db.pool.ExecContext(ctx,
		"UPDATE leetcode_owed SET owed_count = owed_count - ? WHERE user_id = ? AND owed_count >= ?",
		amount, userID, amount)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		var count int
		err := db.pool.QueryRowContext(ctx,
			"SELECT owed_count FROM leetcode_owed WHERE user_id = ?", userID).Scan(&count)
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		if err != nil {
			return err
		}
		return ErrOwedBelowZero
	}
	return nil
}

func (db *TadDatabase) DeleteLeetcodeOwed(ctx context.Context, userID string) error {
	_, err := db.pool.ExecContext(ctx, "DELETE FROM leetcode_owed WHERE user_id = ?", userID)
	return err
}
