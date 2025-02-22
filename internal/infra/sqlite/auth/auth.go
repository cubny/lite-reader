package auth

import (
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
	err := d.sqliteDB.QueryRow("SELECT id, email, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *DB) CreateUser(email, password string) error {
	_, err := d.sqliteDB.Exec("INSERT INTO users (email, password) VALUES (?, ?)", email, password)
	return err
}

func (d *DB) GetAllUsers() ([]*auth.User, error) {
	rows, err := d.sqliteDB.Query("SELECT id, email, password FROM users")
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
		if err := rows.Scan(&user.ID, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
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
	_, err := d.sqliteDB.Exec(`
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
	err := d.sqliteDB.QueryRow(`
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
