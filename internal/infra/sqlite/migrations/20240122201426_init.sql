-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS item (is_new NUMERIC, desc BLOB, id INTEGER PRIMARY KEY, link TEXT, rss_id NUMERIC, title TEXT, dir TEXT, starred NUMBERIC DEFAULT 0, timestamp DATETIME);
CREATE TABLE IF NOT EXISTS rss (id INTEGER PRIMARY KEY, desc TEXT, title TEXT, link TEXT, url TEXT, updated_at TEXT, lang TEXT);
CREATE TABLE IF NOT EXISTS config (key TEXT, value TEXT);
INSERT INTO config VALUES("version",1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS item;
DROP TABLE IF EXISTS rss;
DROP TABLE IF EXISTS config;
-- +goose StatementEnd
