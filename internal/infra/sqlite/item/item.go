package item

import (
	"context"
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

const itemColumns = "id, is_new, desc, link, rss_id, title, dir, starred, timestamp, full_content, scraped_at, scrape_status"

func (r *DB) GetUnreadItems() ([]*item.Item, error) {
	query := "SELECT " + itemColumns + " FROM item WHERE is_new = 1 ORDER BY timestamp DESC"
	result, err := r.sqliteDB.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	return resultToItems(result)
}

func (r *DB) GetStarredItems() ([]*item.Item, error) {
	query := "SELECT " + itemColumns + " FROM item WHERE starred = 1 ORDER BY timestamp DESC"
	result, err := r.sqliteDB.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	return resultToItems(result)
}

func (r *DB) GetFeedItems(feedID int) ([]*item.Item, error) {
	query := "SELECT " + itemColumns + " FROM item WHERE rss_id = ? ORDER BY timestamp DESC"
	result, err := r.sqliteDB.QueryContext(context.Background(), query, feedID)
	if err != nil {
		return nil, err
	}
	return resultToItems(result)
}

func (r *DB) GetFolderItems(folderID, userID int) ([]*item.Item, error) {
	query := "SELECT i.id, i.is_new, i.desc, i.link, i.rss_id, i.title, i.dir, i.starred, i.timestamp, " +
		"i.full_content, i.scraped_at, i.scrape_status " +
		"FROM item i JOIN rss r ON i.rss_id = r.id " +
		"WHERE r.folder_id = ? AND r.user_id = ? ORDER BY i.timestamp DESC"
	result, err := r.sqliteDB.QueryContext(context.Background(), query, folderID, userID)
	if err != nil {
		return nil, err
	}
	return resultToItems(result)
}

// GetItemForUser returns the item if and only if it belongs to a feed owned
// by userID. Returns (nil, nil) when not found so callers can map to a 404.
func (r *DB) GetItemForUser(id, userID int) (*item.Item, error) {
	query := "SELECT i.id, i.is_new, i.desc, i.link, i.rss_id, i.title, i.dir, i.starred, i.timestamp, " +
		"i.full_content, i.scraped_at, i.scrape_status " +
		"FROM item i JOIN rss r ON i.rss_id = r.id " +
		"WHERE i.id = ? AND r.user_id = ?"
	result, err := r.sqliteDB.QueryContext(context.Background(), query, id, userID)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	if !result.Next() {
		return nil, nil
	}
	return resultToItem(result)
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
	var fullContent sql.NullString
	var scrapedAt sql.NullTime
	var scrapeStatus sql.NullString
	err := result.Scan(&id, &isNew, &desc, &link, &rssID, &title, &dir, &starred, &timestamp, &fullContent, &scrapedAt, &scrapeStatus)
	if err != nil {
		return nil, err
	}
	i := &item.Item{
		ID:           id,
		IsNew:        isNew == 1,
		Desc:         desc,
		Link:         link,
		Title:        title,
		Dir:          dir,
		Starred:      starred == 1,
		Timestamp:    timestamp,
		FullContent:  fullContent.String,
		ScrapeStatus: scrapeStatus.String,
	}
	if scrapedAt.Valid {
		t := scrapedAt.Time
		i.ScrapedAt = &t
	}
	return i, nil
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
	result, err := r.sqliteDB.QueryContext(context.Background(), "SELECT id FROM item WHERE link = ? AND rss_id = ?", i.Link, feedID)
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
	insertResult, execErr := r.sqliteDB.ExecContext(context.Background(), query, i.IsNew, i.Desc, i.Link, feedID, i.Title, i.Dir, i.Timestamp)
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
	_, err := r.sqliteDB.ExecContext(context.Background(), "UPDATE item SET starred = ?, is_new = ? WHERE id = ?", starred, isNew, id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateItemContent persists scraped article content (or an error marker) on
// the item. On error rows fullContent will be empty and status will be
// "error"; the caller decides whether to retry.
func (r *DB) UpdateItemContent(id int, fullContent, status string, scrapedAt time.Time) error {
	_, err := r.sqliteDB.ExecContext(context.Background(),
		"UPDATE item SET full_content = ?, scraped_at = ?, scrape_status = ? WHERE id = ?",
		fullContent, scrapedAt, status, id)
	return err
}

func (r *DB) ReadFeedItems(feedID int) error {
	_, err := r.sqliteDB.ExecContext(context.Background(), "UPDATE item SET is_new = 0 WHERE rss_id = ?", feedID)
	if err != nil {
		return err
	}
	return nil
}

func (r *DB) UnreadFeedItems(feedID int) error {
	_, err := r.sqliteDB.ExecContext(context.Background(), "UPDATE item SET is_new = 1 WHERE rss_id = ?", feedID)
	if err != nil {
		return err
	}
	return nil
}

func (r *DB) ReadFolderItems(folderID, userID int) error {
	q := "UPDATE item SET is_new = 0 WHERE rss_id IN " +
		"(SELECT id FROM rss WHERE folder_id = ? AND user_id = ?)"
	_, err := r.sqliteDB.ExecContext(context.Background(), q, folderID, userID)
	return err
}

func (r *DB) GetUnreadItemsCount() (int, error) {
	result, err := r.sqliteDB.QueryContext(context.Background(), "SELECT COUNT(*) FROM item WHERE is_new = 1")
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
	result, err := r.sqliteDB.QueryContext(context.Background(), "SELECT COUNT(*) FROM item WHERE starred = 1")
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
	_, err := r.sqliteDB.ExecContext(context.Background(), "DELETE FROM item WHERE rss_id = ?", feedID)
	return err
}
