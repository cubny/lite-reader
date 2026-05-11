-- +goose Up
CREATE TABLE folder (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  position INTEGER NOT NULL DEFAULT 0,
  user_id INTEGER NOT NULL
);
CREATE INDEX idx_folder_user ON folder(user_id);

ALTER TABLE rss ADD COLUMN folder_id INTEGER;
ALTER TABLE rss ADD COLUMN position INTEGER NOT NULL DEFAULT 0;
CREATE INDEX idx_rss_folder ON rss(folder_id);

-- +goose Down
DROP INDEX IF EXISTS idx_rss_folder;
DROP INDEX IF EXISTS idx_folder_user;
DROP TABLE folder;
