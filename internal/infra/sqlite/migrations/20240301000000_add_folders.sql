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

-- SQLite has limited ALTER TABLE support, so we recreate rss without the
-- folder_id/position columns rather than dropping them in place.
CREATE TABLE rss_new (
  id INTEGER PRIMARY KEY,
  desc TEXT,
  title TEXT,
  link TEXT,
  url TEXT,
  updated_at TEXT,
  lang TEXT,
  user_id INTEGER
);
INSERT INTO rss_new (id, desc, title, link, url, updated_at, lang, user_id)
SELECT id, desc, title, link, url, updated_at, lang, user_id FROM rss;
DROP TABLE rss;
ALTER TABLE rss_new RENAME TO rss;

DROP INDEX IF EXISTS idx_folder_user;
DROP TABLE folder;
