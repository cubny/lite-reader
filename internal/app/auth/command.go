package auth

import "errors"

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
