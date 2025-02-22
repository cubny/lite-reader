package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func (h *Router) login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("login: failed to read request body")
		_ = BadRequest(w, "invalid request body")
		return
	}
	// Restore the body for later use by toLoginCommand
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	log.WithField("body", string(body)).Info("login: received request")
	command, err := toLoginCommand(w, r, nil)
	if err != nil {
		return
	}

	response, err := h.authService.Login(command)
	if err != nil {
		_ = BadRequest(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
