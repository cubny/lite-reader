package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"

	"github.com/cubny/lite-reader/internal/app/auth"
)

func (h *Router) setup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var command auth.SetupCommand
	if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
		log.WithError(err).Error("setup: failed to decode request body")
		_ = BadRequest(w, "invalid request body")
		return
	}

	if err := command.Validate(); err != nil {
		_ = BadRequest(w, err.Error())
		return
	}

	if err := h.authService.Setup(&command); err != nil {
		_ = BadRequest(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Router) needsSetup(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	needs, err := h.authService.NeedsSetup()
	if err != nil {
		_ = InternalError(w, "failed to check setup status")
		return
	}

	allowSignup := h.authService.IsSignupAllowed()

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"needs_setup":  needs,
		"allow_signup": allowSignup,
	})
}

func (h *Router) getSettings(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	allowSignup := h.authService.IsSignupAllowed()

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]bool{"allow_signup": allowSignup})
}

func (h *Router) updateSettings(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req struct {
		AllowSignup bool `json:"allow_signup"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = BadRequest(w, "invalid request body")
		return
	}

	if err := h.authService.SetAllowSignup(req.AllowSignup); err != nil {
		_ = InternalError(w, "failed to update settings")
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]bool{"allow_signup": req.AllowSignup})
}
