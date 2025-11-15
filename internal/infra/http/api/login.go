package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cubny/lite-reader/internal/app/auth"
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

	// Check if this is a form submission (HTMX) or JSON request
	contentType := r.Header.Get("Content-Type")
	isHTMX := r.Header.Get("HX-Request") == "true"

	var command *auth.LoginCommand
	var parseErr error

	if contentType == "application/x-www-form-urlencoded" || isHTMX {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			log.WithError(err).Error("login: failed to parse form")
			if isHTMX {
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`<div class="error-message">Invalid request</div>`))
			} else {
				_ = BadRequest(w, "invalid request body")
			}
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			if isHTMX {
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`<div class="error-message">Email and password are required</div>`))
			} else {
				_ = BadRequest(w, "email and password are required")
			}
			return
		}

		command = &auth.LoginCommand{
			Email:    email,
			Password: password,
		}
	} else {
		// Parse JSON
		command, parseErr = toLoginCommand(w, r, nil)
		if parseErr != nil {
			return
		}
	}

	response, err := h.authService.Login(command)
	if err != nil {
		if isHTMX {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`<div class="error-message">Invalid email or password</div>`))
		} else {
			_ = BadRequest(w, err.Error())
		}
		return
	}

	if isHTMX {
		// Return HTML fragment with token embedded
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		html := fmt.Sprintf(`<div data-token="%s" class="success-message">Login successful</div>`, response.AccessToken)
		_, _ = w.Write([]byte(html))
	} else {
		// Return JSON for backward compatibility
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}
}
