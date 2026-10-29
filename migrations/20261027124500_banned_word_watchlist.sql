-- +goose Up
CREATE TABLE IF NOT EXISTS banned_word_watchlist (
    user_id TEXT PRIMARY KEY,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS banned_word_watchlist;
