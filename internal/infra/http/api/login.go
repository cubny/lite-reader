package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// isHTMXRequest checks if the request is an HTMX request
func isHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == htmxRequestHeader
}

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
		if isHTMXRequest(r) {
			// For HTMX requests, return HTML error message
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`<div class="error-message">Invalid email or password</div>`))
		} else {
			_ = BadRequest(w, err.Error())
		}
		return
	}

	if isHTMXRequest(r) {
		// For HTMX requests, set auth token in a cookie and redirect
		http.SetCookie(w, &http.Cookie{
			Name:     "authToken",
			Value:    response.AccessToken,
			Path:     "/",
			HttpOnly: false, // Allow JavaScript to read for localStorage
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	} else {
		// For JSON requests, return JSON response
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}
}
