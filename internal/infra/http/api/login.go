package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// isHTMXRequest checks if the request is from HTMX
func isHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func (h *Router) login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	command, err := toLoginCommand(w, r, nil)
	if err != nil {
		return
	}

	response, err := h.authService.Login(command)
	if err != nil {
		if isHTMXRequest(r) {
			// Return HTML error for HTMX requests
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`<div class="error-message">Invalid email or password</div>`))
		} else {
			_ = BadRequest(w, err.Error())
		}
		return
	}

	if isHTMXRequest(r) {
		// For HTMX: Store token in cookie and redirect
		http.SetCookie(w, &http.Cookie{
			Name:     "authToken",
			Value:    response.AccessToken,
			Path:     "/",
			HttpOnly: false, // JavaScript needs to read it
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   3600 * 24 * 7, // 7 days
		})
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	} else {
		// JSON API response
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}
}
