-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY,
    username TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sudo_users (
    user_id TEXT PRIMARY KEY,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS commands (
    command TEXT PRIMARY KEY,
    behavior TEXT NOT NULL CHECK (behavior IN ('spam', 'media'))
);

CREATE TABLE IF NOT EXISTS spam_commands (
    command TEXT PRIMARY KEY,
    payload TEXT NOT NULL,
    spam_count INTEGER CHECK (spam_count IS NULL OR spam_count > 0),
    FOREIGN KEY (command) REFERENCES commands(command) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS media_commands (
    command TEXT PRIMARY KEY,
    payload TEXT NOT NULL,
    FOREIGN KEY (command) REFERENCES commands(command) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS banned_words (
    word TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS leetcode_owed (
    user_id TEXT PRIMARY KEY,
    owed_count INTEGER NOT NULL CHECK (owed_count >= 0),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS banned_words;
DROP TABLE IF EXISTS leetcode_owed;
DROP TABLE IF EXISTS media_commands;
DROP TABLE IF EXISTS spam_commands;
DROP TABLE IF EXISTS commands;
DROP TABLE IF EXISTS sudo_users;
DROP TABLE IF EXISTS users;
