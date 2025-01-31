package auth

type User struct {
	Email    string
	Password string // bcrypt hashed
}

type Repository interface {
	GetUserByEmail(email string) (*User, error)
	CreateUser(email, password string) error
}
