package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"

	"github.com/cubny/lite-reader/internal/app/item"
)

func (h *Router) scrapeItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	command, err := toScrapeItemCommand(w, r, p)
	if err != nil {
		return
	}

	updated, err := h.itemService.ScrapeItem(r.Context(), command)
	if err != nil {
		if errors.Is(err, item.ErrItemNotFound) {
			_ = NotFound(w, "item not found")
			return
		}
		log.WithError(err).Errorf("scrapeItem: %s", err)
		_ = InternalError(w, "cannot scrape item")
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toItemResponse(updated)); err != nil {
		log.WithError(err).Errorf("scrapeItem: encoder %s", err)
	}
}

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
}

func (h *Router) getStarredItems(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
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

func (h *Router) getUnreadItems(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
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

func (h *Router) getUnreadItemsCount(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
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

func (h *Router) getStarredItemsCount(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
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
