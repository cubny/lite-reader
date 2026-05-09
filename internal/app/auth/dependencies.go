package auth

type Repository interface {
	GetUserByEmail(email string) (*User, error)
	CreateUser(email, password string) error
	CreateAdmin(email, password string) error
	CreateSession(userID int) (*Session, error)
	GetSessionByToken(token string) (*Session, error)
	GetAllUsers() ([]*User, error)
	CountUsers() (int, error)
	GetSetting(key string) (string, error)
	SetSetting(key, value string) error
}
