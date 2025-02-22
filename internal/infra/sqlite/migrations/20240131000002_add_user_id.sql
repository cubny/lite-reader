-- +goose Up
ALTER TABLE item ADD COLUMN user_id INTEGER;
ALTER TABLE rss ADD COLUMN user_id INTEGER;

-- +goose Down
ALTER TABLE item DROP COLUMN user_id;
ALTER TABLE rss DROP COLUMN user_id;
