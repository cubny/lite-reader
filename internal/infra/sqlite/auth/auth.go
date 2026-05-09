package auth

import (
	"context"

	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/cubny/lite-reader/internal/app/auth"
)

const (
	dayInHours        = 24 * time.Hour
	secureTokenLength = 32
)

type DB struct {
	sqliteDB *sql.DB
}

func NewDB(client *sql.DB) *DB {
	return &DB{sqliteDB: client}
}

func (d *DB) GetUserByEmail(email string) (*auth.User, error) {
	var user auth.User
	err := d.sqliteDB.
		QueryRowContext(context.Background(), "SELECT id, email, password, is_admin FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Email, &user.Password, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *DB) CreateUser(email, password string) error {
	_, err := d.sqliteDB.ExecContext(context.Background(), "INSERT INTO users (email, password) VALUES (?, ?)", email, password)
	return err
}

func (d *DB) CreateAdmin(email, password string) error {
	_, err := d.sqliteDB.ExecContext(context.Background(), "INSERT INTO users (email, password, is_admin) VALUES (?, ?, 1)", email, password)
	return err
}

func (d *DB) GetAllUsers() ([]*auth.User, error) {
	rows, err := d.sqliteDB.QueryContext(context.Background(), "SELECT id, email, password, is_admin FROM users")
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	users := make([]*auth.User, 0)
	for rows.Next() {
		var user auth.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.IsAdmin); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (d *DB) CountUsers() (int, error) {
	var count int
	err := d.sqliteDB.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (d *DB) GetSetting(key string) (string, error) {
	var value string
	err := d.sqliteDB.QueryRowContext(context.Background(), "SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (d *DB) SetSetting(key, value string) error {
	_, err := d.sqliteDB.ExecContext(context.Background(),
		"INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value",
		key, value)
	return err
}

func (d *DB) Login(f *auth.LoginCommand) error {
	user, err := d.GetUserByEmail(f.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("invalid credentials")
		}
		return err
	}

	// Compare passwords using constant-time comparison
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(f.Password)); err != nil {
		return errors.New("invalid credentials")
	}

	return nil
}

func (d *DB) CreateSession(userID int) (*auth.Session, error) {
	session := &auth.Session{
		ID:           uuid.New().String(),
		UserID:       userID,
		AccessToken:  generateSecureToken(),
		RefreshToken: generateSecureToken(),
		ExpiresAt:    time.Now().Add(dayInHours), // Access token expires in 24 hours
		CreatedAt:    time.Now(),
	}

	// Store session in database
	_, err := d.sqliteDB.ExecContext(context.Background(), `
        INSERT INTO sessions (id, user_id, access_token, refresh_token, expires_at, created_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `, session.ID, session.UserID, session.AccessToken, session.RefreshToken, session.ExpiresAt, session.CreatedAt)

	if err != nil {
		return nil, err
	}

	return session, nil
}

func (d *DB) GetSessionByToken(token string) (*auth.Session, error) {
	var session auth.Session
	err := d.sqliteDB.QueryRowContext(context.Background(), `
	    SELECT id, user_id, access_token, refresh_token, expires_at, created_at 
	    FROM sessions 
	    WHERE access_token = ?`, token).Scan(
		&session.ID, &session.UserID, &session.AccessToken,
		&session.RefreshToken, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func generateSecureToken() string {
	b := make([]byte, secureTokenLength)
	if _, err := rand.Read(b); err != nil {
		return "" // Return empty string in case of error instead of panic
	}
	return hex.EncodeToString(b)
}
