package api

import (
	"bytes"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func (h *Router) signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("signup: failed to read request body")
		_ = BadRequest(w, "invalid request body")
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	log.WithField("body", string(body)).Info("signup: received request")
	command, err := toSignupCommand(w, r, nil)
	if err != nil {
		return
	}

	if err := h.authService.Signup(command); err != nil {
		// As authService.Signup now always returns an error,
		// we expect this path to be taken.
		// The error message is "public registration is disabled".
		// We'll return a StatusForbidden.
		_ = Forbidden(w, err.Error())
		return
	}

	// This part should ideally not be reached if authService.Signup always errors.
	// However, to be safe and handle any unexpected success,
	// we can log a warning and return a generic server error.
	log.Warn("signup: authService.Signup unexpectedly succeeded")
	_ = InternalServerError(w, "unexpected server behavior")
}
