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
		if isHTMXRequest(r) {
			// For HTMX requests, return HTML error message
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`<div class="error-message">Error creating account. Email may already be registered.</div>`))
		} else {
			_ = BadRequest(w, err.Error())
		}
		return
	}

	if isHTMXRequest(r) {
		// For HTMX requests, redirect to login page
		w.Header().Set("HX-Redirect", "/login.html")
		w.WriteHeader(http.StatusCreated)
	} else {
		// For JSON requests, return 201 Created
		w.WriteHeader(http.StatusCreated)
	}
}
