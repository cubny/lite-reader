package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (h *Router) logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// For HTMX: Clear cookie and redirect
	if isHTMXRequest(r) {
		http.SetCookie(w, &http.Cookie{
			Name:  "authToken",
			Value: "",
			Path:  "/",
		})
		w.Header().Set("HX-Redirect", "/login.html")
		w.WriteHeader(http.StatusOK)
	} else {
		// For JSON API: just return success
		w.WriteHeader(http.StatusOK)
	}
}
