package auth

type Repository interface {
	GetUserByEmail(email string) (*User, error)
	CreateUser(email, password string) error
	CreateSession(userID int64) (*Session, error)
	GetSessionByToken(token string) (*Session, error)
	GetAllUsers() ([]*User, error)
}
