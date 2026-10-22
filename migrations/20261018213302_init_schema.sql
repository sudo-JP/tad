-- +goose Up
SELECT 'up SQL query';

-- +goose Down
SELECT 'down SQL query';

CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY ,
    username TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sudo_user (
    user_id TEXT PRIMARY KEY,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS action (
    action TEXT PRIMARY KEY,
    payload TEXT
);

CREATE TABLE IF NOT EXISTS banned_words (
    word TEXT PRIMARY KEY,
    action TEXT NOT NULL,
    FOREIGN KEY (action) REFERENCES action(action)
);

CREATE TABLE IF NOT EXISTS commands (
    command TEXT PRIMARY KEY,
    action TEXT NOT NULL,
    FOREIGN KEY (action) REFERENCES action(action)
);
