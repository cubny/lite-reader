package feed

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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
	const q = "INSERT INTO rss (title, desc, link, url, lang, user_id, updated_at, folder_id, position) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := r.sqliteDB.ExecContext(context.Background(), q,
		f.Title, f.Description, f.Link, f.URL, f.Lang, f.UserID,
		f.UpdatedAt.Format(time.RFC3339), nullableInt(f.FolderID), f.Position)
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
		"id, title, desc, link, url, lang, updated_at, folder_id, position, " +
		"(SELECT COUNT(*) FROM item WHERE rss_id = rss.id AND is_new = 1) AS unread_count FROM rss where id = ?"
	rows, err := r.sqliteDB.QueryContext(context.Background(), query, id)
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
		"id, title, desc, link, url, lang, updated_at, folder_id, position, " +
		"(SELECT COUNT(*) FROM item WHERE rss_id = rss.id AND is_new = 1) AS unread_count FROM rss " +
		"WHERE user_id = ? ORDER BY folder_id, position, id"
	rows, err := r.sqliteDB.QueryContext(context.Background(), query, userID)
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
	_, err := r.sqliteDB.ExecContext(context.Background(), "DELETE FROM rss WHERE id = ?", id)
	return err
}

func (r *DB) MoveFeed(feedID int, folderID *int) error {
	_, err := r.sqliteDB.ExecContext(context.Background(),
		"UPDATE rss SET folder_id = ? WHERE id = ?", nullableInt(folderID), feedID)
	return err
}

func (r *DB) ReorderFeed(feedID int, position int) error {
	_, err := r.sqliteDB.ExecContext(context.Background(),
		"UPDATE rss SET position = ? WHERE id = ?", position, feedID)
	return err
}

func (r *DB) BulkMoveFeeds(feedIDs []int, folderID *int) error {
	if len(feedIDs) == 0 {
		return nil
	}
	placeholders := make([]string, len(feedIDs))
	args := make([]interface{}, 0, len(feedIDs)+1)
	args = append(args, nullableInt(folderID))
	for i, id := range feedIDs {
		placeholders[i] = "?"
		args = append(args, id)
	}
	q := fmt.Sprintf("UPDATE rss SET folder_id = ? WHERE id IN (%s)", strings.Join(placeholders, ","))
	_, err := r.sqliteDB.ExecContext(context.Background(), q, args...)
	return err
}

func nullableInt(p *int) interface{} {
	if p == nil {
		return nil
	}
	return *p
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
	var id, unreadCount, position int
	var title, description, link, url, lang, updatedAt string
	var folderID sql.NullInt64
	err := result.Scan(&id, &title, &description, &link, &url, &lang, &updatedAt,
		&folderID, &position, &unreadCount)
	if err != nil {
		return nil, err
	}
	updatedAtTime, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return nil, err
	}
	var fid *int
	if folderID.Valid {
		v := int(folderID.Int64)
		fid = &v
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
		FolderID:    fid,
		Position:    position,
	}, nil
}
