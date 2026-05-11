package folder

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/cubny/lite-reader/internal/app/folder"
)

type DB struct {
	sqliteDB *sql.DB
}

func NewDB(client *sql.DB) *DB {
	return &DB{sqliteDB: client}
}

func (r *DB) AddFolder(f *folder.Folder) (int, error) {
	const q = "INSERT INTO folder (name, position, user_id) VALUES (?, ?, ?)"
	res, err := r.sqliteDB.ExecContext(context.Background(), q, f.Name, f.Position, f.UserID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *DB) GetFolder(id int) (*folder.Folder, error) {
	const q = "SELECT id, name, position, user_id FROM folder WHERE id = ?"
	row := r.sqliteDB.QueryRowContext(context.Background(), q, id)
	f := &folder.Folder{}
	err := row.Scan(&f.ID, &f.Name, &f.Position, &f.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return f, nil
}

func (r *DB) ListFolders(userID int) ([]*folder.Folder, error) {
	query := "SELECT f.id, f.name, f.position, f.user_id, " +
		"(SELECT COUNT(*) FROM item i JOIN rss r ON i.rss_id = r.id " +
		"WHERE r.folder_id = f.id AND i.is_new = 1) AS unread_count " +
		"FROM folder f WHERE f.user_id = ? ORDER BY f.position, f.id"
	rows, err := r.sqliteDB.QueryContext(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			panic(cerr)
		}
	}()
	folders := make([]*folder.Folder, 0)
	for rows.Next() {
		f := &folder.Folder{}
		if err := rows.Scan(&f.ID, &f.Name, &f.Position, &f.UserID, &f.UnreadCount); err != nil {
			return nil, err
		}
		folders = append(folders, f)
	}
	return folders, nil
}

func (r *DB) RenameFolder(id int, name string) error {
	_, err := r.sqliteDB.ExecContext(context.Background(),
		"UPDATE folder SET name = ? WHERE id = ?", name, id)
	return err
}

func (r *DB) ReorderFolder(id int, position int) error {
	_, err := r.sqliteDB.ExecContext(context.Background(),
		"UPDATE folder SET position = ? WHERE id = ?", position, id)
	return err
}

func (r *DB) DeleteFolder(id int) error {
	tx, err := r.sqliteDB.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(context.Background(),
		"UPDATE rss SET folder_id = NULL WHERE folder_id = ?", id); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err := tx.ExecContext(context.Background(),
		"DELETE FROM folder WHERE id = ?", id); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
