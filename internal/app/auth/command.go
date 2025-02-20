package auth

import (
	"errors"
)

type LoginCommand struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *LoginCommand) Validate() error {
	if c.Email == "" || c.Password == "" {
		return errors.New("email and password are required")
	}
	return nil
}

type SignupCommand struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (c *SignupCommand) Validate() error {
	if c.Email == "" || c.Password == "" || c.ConfirmPassword == "" {
		return errors.New("email and password are required")
	}
	if c.Password != c.ConfirmPassword {
		return errors.New("password and confirm password do not match")
	}
	if len(c.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}
