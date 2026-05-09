package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const settingAllowSignup = "allow_signup"

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

	if bErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(command.Password)); bErr != nil {
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
	// Check if signup is allowed
	if !s.IsSignupAllowed() {
		return errors.New("registration is currently disabled")
	}

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

	err = s.repo.CreateUser(command.Email, string(hashedPassword))
	switch {
	case err == nil:
		return nil
	case strings.Contains(err.Error(), "UNIQUE constraint failed"):
		return errors.New("email already registered")
	default:
		return err
	}
}

// Setup creates the first admin user. It only works when there are zero users
// in the database.
func (s *Service) Setup(command *SetupCommand) error {
	count, err := s.repo.CountUsers()
	if err != nil {
		return fmt.Errorf("failed to check user count: %w", err)
	}

	if count > 0 {
		return errors.New("setup has already been completed")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(command.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.repo.CreateAdmin(command.Email, string(hashedPassword)); err != nil {
		return err
	}

	allowSignup := "false"
	if command.AllowSignup {
		allowSignup = "true"
	}

	return s.repo.SetSetting(settingAllowSignup, allowSignup)
}

// NeedsSetup returns true if no users exist yet (first-run state).
func (s *Service) NeedsSetup() (bool, error) {
	count, err := s.repo.CountUsers()
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// IsSignupAllowed checks the allow_signup setting. Defaults to false on error.
func (s *Service) IsSignupAllowed() bool {
	val, err := s.repo.GetSetting(settingAllowSignup)
	if err != nil {
		return false
	}
	return val == "true"
}

// SetAllowSignup toggles the allow_signup setting.
func (s *Service) SetAllowSignup(allow bool) error {
	val := "false"
	if allow {
		val = "true"
	}
	return s.repo.SetSetting(settingAllowSignup, val)
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
