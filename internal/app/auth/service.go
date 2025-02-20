package auth

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Login(command *LoginCommand) (*LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(command.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(command.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate session token
	session, err := s.repo.CreateSession(user.ID)
	if err != nil {
		return nil, errors.New("failed to create session")
	}

	return &LoginResponse{
		User:         *user,
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		ExpiresIn:    time.Until(session.ExpiresAt).Seconds(),
	}, nil

}

func (s *Service) Signup(command *SignupCommand) error {
	// Check if user already exists
	_, err := s.repo.GetUserByEmail(command.Email)
	if err == nil {
		return errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(command.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.CreateUser(command.Email, string(hashedPassword))
}

func (s *Service) GetSession(token string) (*Session, error) {
	session, err := s.repo.GetSessionByToken(token)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Service) GetAllUsers() ([]*User, error) {
	return s.repo.GetAllUsers()
}
