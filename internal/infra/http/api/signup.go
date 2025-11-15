package api

import (
	"bytes"
	"io"
	"net/http"

	"github.com/cubny/lite-reader/internal/app/auth"
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

	// Check if this is a form submission (HTMX) or JSON request
	contentType := r.Header.Get("Content-Type")
	isHTMX := r.Header.Get("HX-Request") == "true"

	var command *auth.SignupCommand
	var parseErr error

	if contentType == "application/x-www-form-urlencoded" || isHTMX {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			log.WithError(err).Error("signup: failed to parse form")
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

		command = &auth.SignupCommand{
			Email:    email,
			Password: password,
		}
	} else {
		// Parse JSON
		command, parseErr = toSignupCommand(w, r, nil)
		if parseErr != nil {
			return
		}
	}

	if err := h.authService.Signup(command); err != nil {
		if isHTMX {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`<div class="error-message">Error creating account. Email may already be registered.</div>`))
		} else {
			_ = BadRequest(w, err.Error())
		}
		return
	}

	if isHTMX {
		// Return HTML fragment with success indicator
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`<div data-success="true" class="success-message">Account created successfully</div>`))
	} else {
		// Return empty response for JSON (backward compatibility)
		w.WriteHeader(http.StatusCreated)
	}
}
