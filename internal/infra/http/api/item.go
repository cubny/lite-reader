package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func (h *Router) updateItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	command, err := toUpdateItemCommand(w, r, p)
	if err != nil {
		return
	}

	if err := h.itemService.UpdateItem(command); err != nil {
		_ = InternalError(w, "cannot update item")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

func (h *Router) getStarredItems(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	items, err := h.itemService.GetStarredItems()
	if err != nil {
		_ = InternalError(w, "cannot get unread items")
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toGetItemsResponse(items)); err != nil {
		log.WithError(err).Errorf("getStarredItems: encoder %s", err)
		_ = InternalError(w, "cannot encode response")
		return
	}
}

func (h *Router) getUnreadItems(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	items, err := h.itemService.GetUnreadItems()
	if err != nil {
		_ = InternalError(w, "cannot get unread items")
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toGetItemsResponse(items)); err != nil {
		log.WithError(err).Errorf("getUnreadItems: encoder %s", err)
		_ = InternalError(w, "cannot encode response")
		return
	}
}

func (h *Router) getUnreadItemsCount(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	items, err := h.itemService.GetUnreadItemsCount()
	if err != nil {
		_ = InternalError(w, "cannot get unread items")
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toGetItemsCountResponse(items)); err != nil {
		log.WithError(err).Errorf("getUnreadItemsCount: encoder %s", err)
		_ = InternalError(w, "cannot encode response")
		return
	}
}

func (h *Router) getStarredItemsCount(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	items, err := h.itemService.GetStarredItemsCount()
	if err != nil {
		_ = InternalError(w, "cannot get unread items")
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toGetItemsCountResponse(items)); err != nil {
		log.WithError(err).Errorf("getStarredItemsCount: encoder %s", err)
		_ = InternalError(w, "cannot encode response")
		return
	}
}
