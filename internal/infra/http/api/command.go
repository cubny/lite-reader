package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kgedala/go-web-app/internal/app/auth" // Assuming kgedala/go-web-app is the correct base path
	log "github.com/sirupsen/logrus"
)

// toLoginCommand parses the request body into an auth.LoginCommand.
func toLoginCommand(w http.ResponseWriter, r *http.Request, _ httprouter.Params) (*auth.LoginCommand, error) {
	var command auth.LoginCommand
	if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
		log.WithError(err).Error("toLoginCommand: failed to decode request body")
		_ = BadRequest(w, "invalid request body: "+err.Error())
		return nil, err
	}

	if err := command.Validate(); err != nil {
		log.WithError(err).Error("toLoginCommand: validation failed")
		_ = BadRequest(w, err.Error())
		return nil, err
	}
	return &command, nil
}

// toSignupCommand parses the request body into an auth.SignupCommand.
func toSignupCommand(w http.ResponseWriter, r *http.Request, _ httprouter.Params) (*auth.SignupCommand, error) {
	var command auth.SignupCommand
	if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
		log.WithError(err).Error("toSignupCommand: failed to decode request body")
		_ = BadRequest(w, "invalid request body: "+err.Error())
		return nil, err
	}

	if err := command.Validate(); err != nil {
		log.WithError(err).Error("toSignupCommand: validation failed")
		_ = BadRequest(w, err.Error())
		return nil, err
	}
	return &command, nil
}

// toCreateUserCommand parses the request body into an auth.CreateUserCommand.
func toCreateUserCommand(w http.ResponseWriter, r *http.Request, _ httprouter.Params) (*auth.CreateUserCommand, error) {
	var command auth.CreateUserCommand
	if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
		log.WithError(err).Error("toCreateUserCommand: failed to decode request body")
		_ = BadRequest(w, "invalid request body: "+err.Error())
		return nil, err
	}

	if err := command.Validate(); err != nil {
		log.WithError(err).Error("toCreateUserCommand: validation failed")
		_ = BadRequest(w, err.Error())
		return nil, err
	}
	return &command, nil
}
