-- +goose Up
SELECT 'up SQL query';

-- +goose Down
SELECT 'down SQL query';

CREATE TABLE IF NOT EXISTS users (
    user_id PRIMARY KEY TEXT,
    username TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sudo_user (
    user_id TEXT NOT NULL,
    FOREIGN KEY
);
