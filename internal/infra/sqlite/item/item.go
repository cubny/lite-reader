package item

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"

	"github.com/cubny/lite-reader/internal/app/item"
)

type DB struct {
	sqliteDB *sql.DB
}

func NewDB(client *sql.DB) *DB {
	return &DB{sqliteDB: client}
}

func (r *DB) GetUnreadItems() ([]*item.Item, error) {
	query := "SELECT id, is_new, desc, link, rss_id, title, dir, starred, timestamp FROM item WHERE is_new = 1 ORDER BY timestamp DESC"
	result, err := r.sqliteDB.Query(query)
	if err != nil {
		return nil, err
	}
	return resultToItems(result)
}

func (r *DB) GetStarredItems() ([]*item.Item, error) {
	query := "SELECT id, is_new, desc, link, rss_id, title, dir, starred, timestamp FROM item WHERE starred = 1 ORDER BY timestamp DESC"
	result, err := r.sqliteDB.Query(query)
	if err != nil {
		return nil, err
	}
	return resultToItems(result)
}

func (r *DB) GetFeedItems(feedID int) ([]*item.Item, error) {
	query := "SELECT id, is_new, desc, link, rss_id, title, dir, starred, timestamp FROM item WHERE rss_id = ? ORDER BY timestamp DESC"
	result, err := r.sqliteDB.Query(query, feedID)
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

func (r *DB) UpsertItem(feedID int, i *item.Item) (int, error) {
	// first find if item exists
	result, err := r.sqliteDB.Query("SELECT id FROM item WHERE link = ?", i.Link)
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
	if resultErr := result.Close(); resultErr != nil {
		return 0, resultErr
	}
	// if item exists, ignore it
	if id != 0 {
		return id, nil
	}
	// if item does not exist, insert it
	query := "INSERT INTO item (is_new, desc, link, rss_id, title, dir, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?)"
	insertResult, execErr := r.sqliteDB.Exec(query, i.IsNew, i.Desc, i.Link, feedID, i.Title, i.Dir, i.Timestamp)
	if execErr != nil {
		return 0, execErr
	}
	lastInsertID, lastInsertError := insertResult.LastInsertId()
	if lastInsertError != nil {
		return 0, lastInsertError
	}
	return int(lastInsertID), nil
}

func (r *DB) UpdateItem(id int, starred, isNew bool) error {
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
