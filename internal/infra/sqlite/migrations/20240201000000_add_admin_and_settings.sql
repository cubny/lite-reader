-- +goose Up
ALTER TABLE users ADD COLUMN is_admin INTEGER NOT NULL DEFAULT 0;

CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

INSERT OR IGNORE INTO settings (key, value) VALUES ('allow_signup', 'false');

-- +goose Down
-- SQLite does not support DROP COLUMN, so we recreate the table
CREATE TABLE users_backup (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO users_backup SELECT id, email, password, created_at FROM users;
DROP TABLE users;
ALTER TABLE users_backup RENAME TO users;

DROP TABLE IF EXISTS settings;
