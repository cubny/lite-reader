package api

import (
	"bytes"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func (h *Router) createUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Read the body for logging and then restore it for parsing
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("createUser: failed to read request body")
		_ = BadRequest(w, "invalid request body")
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body)) // Restore the body for toCreateUserCommand

	log.WithField("body", string(body)).Info("createUser: received request")

	command, err := toCreateUserCommand(w, r, ps) // This is defined in command.go
	if err != nil {
		// Error response is already sent by toCreateUserCommand
		return
	}

	if err := h.authService.CreateUser(command); err != nil {
		// In a real app, you might check for specific error types from authService
		// and return different HTTP status codes, e.g., http.StatusConflict for duplicate email.
		// For now, we'll use http.StatusBadRequest for all creation errors.
		log.WithError(err).Error("createUser: failed to create user")
		_ = BadRequest(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Info("createUser: user created successfully")
}
