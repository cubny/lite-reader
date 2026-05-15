-- +goose Up
ALTER TABLE item ADD COLUMN full_content TEXT;
ALTER TABLE item ADD COLUMN scraped_at DATETIME;
ALTER TABLE item ADD COLUMN scrape_status TEXT;

-- +goose Down
ALTER TABLE item DROP COLUMN full_content;
ALTER TABLE item DROP COLUMN scraped_at;
ALTER TABLE item DROP COLUMN scrape_status;
