package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Login(command *LoginCommand) (string, error) {
	user, err := s.repo.GetUserByEmail(command.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(command.Password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// For now just return email as the session token
	// TODO: Implement proper JWT or session tokens
	return user.Email, nil
}
