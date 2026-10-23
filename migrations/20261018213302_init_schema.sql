-- +goose Up
SELECT 'up SQL query';

-- +goose Down
SELECT 'down SQL query';

CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY ,
    username TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sudo_users (
    user_id TEXT PRIMARY KEY,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS actions (
    action TEXT PRIMARY KEY,
    payload TEXT
);

CREATE TABLE IF NOT EXISTS banned_words (
    word TEXT PRIMARY KEY,
    action TEXT NOT NULL,
    FOREIGN KEY (action) REFERENCES actions(action)
);

CREATE TABLE IF NOT EXISTS commands (
    command TEXT PRIMARY KEY,
    action TEXT NOT NULL,
    FOREIGN KEY (action) REFERENCES actions(action)
);
