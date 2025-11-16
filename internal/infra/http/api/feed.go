package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"

	"github.com/cubny/lite-reader/internal/app/item"
	"github.com/cubny/lite-reader/internal/infra/http/api/cxutil"
)

// addFeed is the handler for
// swagger:route POST /feeds AddFeedResponse
//
// Schedule a new feed.
//
// Responses:
//
//	201: setFeeds
//	400: invalidRequestBody
//	404: notFoundError
//	422: invalidParams
//	500: serverError
func (h *Router) addFeed(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	command, err := toAddFeedCommand(w, r, p)
	if err != nil {
		if isHTMXRequest(r) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(renderFeedError("Invalid feed URL")))
		}
		return
	}

	log.Infof("addFeed: command %v", command)
	t, err := h.feedService.AddFeed(command)
	if err != nil {
		log.WithError(err).Errorf("addFeed: service %s", err)
		if isHTMXRequest(r) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(renderFeedError("Failed to add feed")))
			return
		}
		_ = InternalError(w, "failed to add feed due to server internal error")
		return
	}

	// For HTMX requests, get all feeds and return the complete list
	if isHTMXRequest(r) {
		userID := r.Context().Value(cxutil.UserIDKey).(int)
		feeds, err := h.feedService.ListFeeds(userID)
		if err != nil {
			log.WithError(err).Errorf("addFeed: listFeeds %s", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(renderFeedError("Failed to refresh feed list")))
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(renderFeedList(feeds)))
		return
	}

	// JSON response for non-HTMX requests
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(toAddFeedResponse(t)); err != nil {
		log.WithError(err).Errorf("setFeed: encoder %s", err)
		_ = InternalError(w, "cannot encode response")
	}
}
func (h *Router) listFeeds(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := r.Context().Value(cxutil.UserIDKey).(int)
	log.Infof("listFeeds: userID %d", userID)
	resp, err := h.feedService.ListFeeds(userID)
	if err != nil {
		log.WithError(err).Errorf("listFeeds: service %s", err)
		if isHTMXRequest(r) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(renderFeedError("Cannot list feeds")))
			return
		}
		_ = InternalError(w, "cannot list feeds")
		return
	}

	// For HTMX requests, return HTML fragment
	if isHTMXRequest(r) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(renderFeedList(resp)))
		return
	}

	// JSON response for non-HTMX requests
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toListFeedResponse(resp)); err != nil {
		log.WithError(err).Errorf("listFeeds: encoder %s", err)
		_ = InternalError(w, "cannot encode response")
		return
	}
}

func (h *Router) getFeedItems(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	command, err := toGetFeedItemsCommand(w, r, p)
	if err != nil {
		return
	}

	items, err := h.itemService.GetFeedItems(command)
	if err != nil {
		_ = InternalError(w, "cannot get feed items")
		return
	}

	resp := toGetItemsResponse(items)
	log.Infof("getFeedItems: resp %v", resp)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.WithError(err).Errorf("getFeedItems: encoder %s", err)
		_ = InternalError(w, "cannot encode response")
		return
	}
}

func (h *Router) fetchFeedNewItems(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	command, err := toFetchFeedNewItemsCommand(w, r, p)
	if err != nil {
		log.WithError(err).Errorf("fetchFeedNewItems: toFetchFeedNewItemsCommand")
		return
	}

	items, err := h.feedService.FetchItems(command.FeedID)
	if err != nil {
		log.WithError(err).Errorf("fetchFeedNewItems: FetchItems")
		_ = InternalError(w, "cannot fetch feed items")
		return
	}

	upsertItemsCommand := &item.UpsertItemsCommand{FeedID: command.FeedID, Items: items}
	if upsertErr := h.itemService.UpsertItems(upsertItemsCommand); upsertErr != nil {
		log.WithError(upsertErr).Errorf("fetchFeedNewItems: UpsertItems")
		_ = InternalError(w, "cannot store feed items")
		return
	}

	getFeedItemsCommand := &item.GetFeedItemsCommand{FeedID: command.FeedID}
	items, err = h.itemService.GetFeedItems(getFeedItemsCommand)
	if err != nil {
		log.WithError(err).Errorf("fetchFeedNewItems: GetFeedItems")
		_ = InternalError(w, "cannot get feed items")
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toGetItemsResponse(items)); err != nil {
		log.WithError(err).Errorf("fetchFeedNewItems: encoder %s", err)
		_ = InternalError(w, "cannot encode response")
		return
	}
}

func (h *Router) readFeedItems(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	command, err := toReadFeedItemsCommand(w, r, p)
	if err != nil {
		return
	}

	if err := h.itemService.ReadFeedItems(command); err != nil {
		_ = InternalError(w, "cannot read feed items")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Router) unreadFeedItems(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	command, err := toUnreadFeedItemCommand(w, r, p)
	if err != nil {
		return
	}

	if err := h.itemService.UnreadFeedItems(command); err != nil {
		_ = InternalError(w, "cannot unread feed items")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Router) deleteFeed(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	command, err := toDeleteFeedCommand(w, r, p)
	if err != nil {
		if isHTMXRequest(r) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(renderFeedError("Invalid feed ID")))
		}
		return
	}

	cmdDeleteFeedItems, err := toDeleteFeedItemsCommand(w, r, p)
	if err != nil {
		if isHTMXRequest(r) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(renderFeedError("Invalid feed ID")))
		}
		return
	}

	if err := h.itemService.DeleteFeedItems(cmdDeleteFeedItems); err != nil {
		if isHTMXRequest(r) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(renderFeedError("Cannot delete feed")))
			return
		}
		_ = InternalError(w, "cannot delete feed")
		return
	}

	if err := h.feedService.DeleteFeed(command); err != nil {
		if isHTMXRequest(r) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(renderFeedError("Cannot delete feed")))
			return
		}
		_ = InternalError(w, "cannot delete feed")
		return
	}

	// For HTMX requests, return updated feed list
	if isHTMXRequest(r) {
		userID := r.Context().Value(cxutil.UserIDKey).(int)
		feeds, err := h.feedService.ListFeeds(userID)
		if err != nil {
			log.WithError(err).Errorf("deleteFeed: listFeeds %s", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(renderFeedError("Failed to refresh feed list")))
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(renderFeedList(feeds)))
		return
	}

	// JSON response for non-HTMX requests
	w.WriteHeader(http.StatusNoContent)
}
