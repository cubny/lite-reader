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
		_ = BadRequest(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}
