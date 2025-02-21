package feed

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"

	"github.com/cubny/lite-reader/internal/app/feed"
)

type DB struct {
	sqliteDB *sql.DB
}

func NewDB(client *sql.DB) *DB {
	return &DB{sqliteDB: client}
}

func (r *DB) AddFeed(f *feed.Feed) (int, error) {
	result, err := r.sqliteDB.Exec("INSERT INTO rss (title, desc, link, url, lang, user_id, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		f.Title, f.Description, f.Link, f.URL, f.Lang, f.UserID, f.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastInsertID), nil
}

func (r *DB) GetFeed(id int) (*feed.Feed, error) {
	query := "SELECT " +
		"id, title, desc, link, url, lang, updated_at, " +
		"(SELECT COUNT(*) FROM item WHERE rss_id = rss.id AND is_new = 1) AS unread_count FROM rss where id = ?"
	rows, err := r.sqliteDB.Query(query, id)
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

func (r *DB) ListFeeds(userID int) ([]*feed.Feed, error) {
	query := "SELECT " +
		"id, title, desc, link, url, lang, updated_at, " +
		"(SELECT COUNT(*) FROM item WHERE rss_id = rss.id AND is_new = 1) AS unread_count FROM rss " +
		"WHERE user_id = ?"
	rows, err := r.sqliteDB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			panic(err)
		}
	}()
	return resultToFeeds(rows)
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
		ID:          id,
		Title:       title,
		Description: description,
		Link:        link,
		URL:         url,
		Lang:        lang,
		UpdatedAt:   updatedAtTime,
		UnreadCount: unreadCount,
	}, nil
}
