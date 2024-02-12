package item

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/cubny/lite-reader/internal/app/item"
)

type DB struct {
	sqliteDB *sql.DB
}

func NewDB(client *sql.DB) *DB {
	return &DB{sqliteDB: client}
}

// CREATE TABLE IF NOT EXISTS item (is_new NUMERIC, desc BLOB, id INTEGER PRIMARY KEY, link TEXT, rss_id NUMERIC, title TEXT, dir TEXT, starred NUMBERIC DEFAULT 0, timestamp DATETIME);
func (r *DB) GetUnreadItems() ([]*item.Item, error) {
	result, err := r.sqliteDB.Query("SELECT id, is_new, desc, link, rss_id, title, dir, starred, timestamp FROM item WHERE is_new = 1 ORDER BY timestamp DESC")
	if err != nil {
		return nil, err
	}
	return resultToItems(result)
}

func (r *DB) GetStarredItems() ([]*item.Item, error) {
	result, err := r.sqliteDB.Query("SELECT id, is_new, desc, link, rss_id, title, dir, starred, timestamp FROM item WHERE starred = 1 ORDER BY timestamp DESC")
	if err != nil {
		return nil, err
	}
	return resultToItems(result)
}

func (r *DB) GetFeedItems(feedID int) ([]*item.Item, error) {
	result, err := r.sqliteDB.Query("SELECT id, is_new, desc, link, rss_id, title, dir, starred, timestamp FROM item WHERE rss_id = ? ORDER BY timestamp DESC", feedID)
	if err != nil {
		return nil, err
	}
	return resultToItems(result)
}

func resultToItem(result *sql.Rows) (*item.Item, error) {
	var id int
	var isNew int
	var desc string
	var link string
	var rssID int
	var title string
	var dir string
	var starred int
	var timestamp time.Time
	err := result.Scan(&id, &isNew, &desc, &link, &rssID, &title, &dir, &starred, &timestamp)
	if err != nil {
		return nil, err
	}
	return &item.Item{
		ID:        id,
		IsNew:     isNew == 1,
		Desc:      desc,
		Link:      link,
		Title:     title,
		Dir:       dir,
		Starred:   starred == 1,
		Timestamp: timestamp,
	}, nil
}

func resultToItems(result *sql.Rows) ([]*item.Item, error) {
	items := make([]*item.Item, 0)
	for result.Next() {
		i, err := resultToItem(result)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := result.Close(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *DB) UpsertItems(feedID int, items []*item.Item) error {
	for _, t := range items {
		_, err := r.UpsertItem(feedID, t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *DB) UpsertItem(feedID int, item *item.Item) (int, error) {
	// first find if item exists
	result, err := r.sqliteDB.Query("SELECT id FROM item WHERE link = ?", item.Link)
	if err != nil {
		return 0, err
	}
	var id int
	if result.Next() {
		err = result.Scan(&id)
		if err != nil {
			return 0, err
		}
	}
	if err := result.Close(); err != nil {
		return 0, err
	}
	// if item exists, ignore it
	if id != 0 {
		return id, nil
	}
	// if item does not exist, insert it
	insertResult, err := r.sqliteDB.Exec("INSERT INTO item (is_new, desc, link, rss_id, title, dir, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?)",
		item.IsNew, item.Desc, item.Link, feedID, item.Title, item.Dir, item.Timestamp)
	if err != nil {
		return 0, err
	}
	lastInsertID, err := insertResult.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastInsertID), nil
}

func (r *DB) UpdateItem(id int, starred bool, isNew bool) error {
	_, err := r.sqliteDB.Exec("UPDATE item SET starred = ?, is_new = ? WHERE id = ?", starred, isNew, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *DB) ReadFeedItems(feedID int) error {
	_, err := r.sqliteDB.Exec("UPDATE item SET is_new = 0 WHERE rss_id = ?", feedID)
	if err != nil {
		return err
	}
	return nil
}

func (r *DB) UnreadFeedItems(feedID int) error {
	_, err := r.sqliteDB.Exec("UPDATE item SET is_new = 1 WHERE rss_id = ?", feedID)
	if err != nil {
		return err
	}
	return nil
}

func (r *DB) GetUnreadItemsCount() (int, error) {
	result, err := r.sqliteDB.Query("SELECT COUNT(*) FROM item WHERE is_new = 1")
	if err != nil {
		return 0, err
	}
	var count int
	if result.Next() {
		err = result.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	if err := result.Close(); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *DB) GetStarredItemsCount() (int, error) {
	result, err := r.sqliteDB.Query("SELECT COUNT(*) FROM item WHERE starred = 1")
	if err != nil {
		return 0, err
	}
	var count int
	if result.Next() {
		err = result.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	if err := result.Close(); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *DB) DeleteFeedItems(feedID int) error {
	_, err := r.sqliteDB.Exec("DELETE FROM item WHERE rss_id = ?", feedID)
	return err
}
