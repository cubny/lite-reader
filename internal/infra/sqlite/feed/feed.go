package feed

import (
	"database/sql"
	"time"

	"github.com/cubny/lite-reader/internal/app/feed"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	sqliteDB *sql.DB
}

func NewDB(client *sql.DB) *DB {
	return &DB{sqliteDB: client}
}

func (r *DB) AddFeed(feed *feed.Feed) (int, error) {
	result, err := r.sqliteDB.Exec("INSERT INTO rss (title, desc, link, url, lang, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		feed.Title, feed.Description, feed.Link, feed.URL, feed.Lang, feed.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastInsertId), nil
}

func (r *DB) GetFeed(id int) (*feed.Feed, error) {
	rows, err := r.sqliteDB.Query("SELECT id, title, desc, link, url, lang, updated_at FROM rss WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			panic(err)
		}
	}()
	for rows.Next() {
		return resultToFeed(rows)
	}
	return nil, nil
}

func (r *DB) ListFeeds() ([]*feed.Feed, error) {
	query := "SELECT " +
		"id, title, desc, link, url, lang, updated_at, " +
		"(SELECT COUNT(*) FROM item WHERE rss_id = rss.id AND is_new = 1) AS unread_count FROM rss"
	result, err := r.sqliteDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = result.Close(); err != nil {
			panic(err)
		}
	}()
	return resultToFeeds(result)
}

func (r *DB) DeleteFeed(id int) error {
	_, err := r.sqliteDB.Exec("DELETE FROM rss WHERE id = ?", id)
	return err
}

func resultToFeeds(result *sql.Rows) ([]*feed.Feed, error) {
	feeds := make([]*feed.Feed, 0)
	for result.Next() {
		f, err := resultToFeed(result)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, f)
	}
	return feeds, nil
}

func resultToFeed(result *sql.Rows) (*feed.Feed, error) {
	var id, unreadCount int
	var title, description, link, url, lang, updatedAt string
	err := result.Scan(&id, &title, &description, &link, &url, &lang, &updatedAt, &unreadCount)
	if err != nil {
		return nil, err
	}
	updatedAtTime, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return nil, err
	}
	return &feed.Feed{
		Id:          id,
		Title:       title,
		Description: description,
		Link:        link,
		URL:         url,
		Lang:        lang,
		UpdatedAt:   updatedAtTime,
		UnreadCount: unreadCount,
	}, nil
}
