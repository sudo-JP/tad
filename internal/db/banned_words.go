package db

import "context"

func (db *TadDatabase) CreateBannedWord(ctx context.Context, word, command string) error {
	_, err := db.pool.ExecContext(ctx,
		"INSERT INTO banned_words (word, command) VALUES (?, ?)", word, command)
	return err
}

func (db *TadDatabase) GetBannedWord(ctx context.Context, word string) (BannedWord, error) {
	var result BannedWord
	err := db.pool.QueryRowContext(ctx, `
		SELECT word, command FROM banned_words WHERE word = ?`, word).
		Scan(&result.Word, &result.Command)
	return result, err
}

func (db *TadDatabase) GetBannedWordsRecords(ctx context.Context) ([]BannedWord, error) {
	rows, err := db.pool.QueryContext(ctx, `
		SELECT word, command FROM banned_words`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []BannedWord
	for rows.Next() {
		var word BannedWord
		if err := rows.Scan(&word.Word, &word.Command); err != nil {
			return nil, err
		}
		words = append(words, word)
	}
	return words, rows.Err()
}

func (db *TadDatabase) GetBannedWords(ctx context.Context) ([]string, error) {
	words, err := db.GetBannedWordsRecords(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(words))
	for _, word := range words {
		result = append(result, word.Word)
	}
	return result, nil
}

func (db *TadDatabase) UpdateBannedWord(ctx context.Context, word, command string) error {
	_, err := db.pool.ExecContext(ctx,
		"UPDATE banned_words SET command = ? WHERE word = ?", command, word)
	return err
}

func (db *TadDatabase) DeleteBannedWord(ctx context.Context, word string) error {
	_, err := db.pool.ExecContext(ctx, "DELETE FROM banned_words WHERE word = ?", word)
	return err
}
