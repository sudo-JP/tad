package db

import "context"

func (db *TadDatabase) CreateBannedWord(ctx context.Context, word string) error {
	_, err := db.pool.ExecContext(ctx, "INSERT INTO banned_words (word) VALUES (?)", word)
	return translateErr(err)
}

func (db *TadDatabase) GetBannedWords(ctx context.Context) ([]string, error) {
	rows, err := db.pool.QueryContext(ctx, "SELECT word FROM banned_words")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []string
	for rows.Next() {
		var word string
		if err := rows.Scan(&word); err != nil {
			return nil, err
		}
		words = append(words, word)
	}
	return words, rows.Err()
}

func (db *TadDatabase) DeleteBannedWord(ctx context.Context, word string) error {
	return requireAffected(db.pool.ExecContext(ctx, "DELETE FROM banned_words WHERE word = ?", word))
}
