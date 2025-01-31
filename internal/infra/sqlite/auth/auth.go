package auth

import (
	"database/sql"

	"github.com/cubny/lite-reader/internal/app/auth"
)

type DB struct {
	sqliteDB *sql.DB
}

func NewDB(client *sql.DB) *DB {
	return &DB{sqliteDB: client}
}

func (d *DB) GetUserByEmail(email string) (*auth.User, error) {
	var user auth.User
	err := d.sqliteDB.QueryRow("SELECT email, password FROM users WHERE email = ?", email).Scan(&user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *DB) CreateUser(email, password string) error {
	_, err := d.sqliteDB.Exec("INSERT INTO users (email, password) VALUES (?, ?)", email, password)
	return err
}

func (r *DB) Login(f *auth.LoginCommand) error {
	return nil
}
