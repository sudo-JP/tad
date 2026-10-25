package db

import "context"

func (db *TadDatabase) CreateBannedWord(ctx context.Context, bannedWord string, action string, payload string) error {
	tx, err := db.pool.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Banned Word query
	bannedWordQuery := "INSERT INTO banned_words (word, action) VALUES (?, ?)"
	_, err = tx.ExecContext(ctx, bannedWordQuery, bannedWord, action)

	actionQuery := "INSERT INTO action (action, payload) VALUES (?, ?)"
	_, err = tx.ExecContext(ctx, actionQuery, action, payload)

	return tx.Commit()
}

func (db *TadDatabase) GetBannedWords(ctx context.Context) ([]string, error) {
	query := "SELECT word FROM banned_words"
	rows, err := db.pool.Query(query)
	if err != nil {
		return []string{}, err
	}

	var words []string
	for rows.Next() {
		var word string
		err := rows.Scan(&word)
		if err != nil {
			return []string{}, err
		}
		words = append(words, word)
	}

	if err = rows.Err(); err != nil {
		return []string{}, err
	}

	return words, err
}
