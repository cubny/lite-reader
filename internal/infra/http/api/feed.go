package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"

	"github.com/cubny/lite-reader/internal/app/item"
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
		return
	}

	log.Infof("addFeed: command %v", command)
	// define t as a new uuid
	t, err := h.feedService.AddFeed(command)
	if err != nil {
		log.WithError(err).Errorf("addFeed: service %s", err)
		_ = InternalError(w, "failed to add feed due to server internal error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(toAddFeedResponse(t)); err != nil {
		log.WithError(err).Errorf("setFeed: encoder %s", err)
		_ = InternalError(w, "cannot encode response")
		return
	}
}
func (h *Router) listFeeds(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	resp, err := h.feedService.ListFeeds()
	if err != nil {
		log.WithError(err).Errorf("listFeeds: service %s", err)
		_ = InternalError(w, "cannot list feeds")
		return
	}

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

	items, err := h.feedService.FetchItems(command.FeedId)
	if err != nil {
		log.WithError(err).Errorf("fetchFeedNewItems: FetchItems")
		_ = InternalError(w, "cannot fetch feed items")
		return
	}

	upsertItemsCommand := &item.UpsertItemsCommand{FeedId: command.FeedId, Items: items}
	if err := h.itemService.UpsertItems(upsertItemsCommand); err != nil {
		log.WithError(err).Errorf("fetchFeedNewItems: UpsertItems")
		_ = InternalError(w, "cannot store feed items")
		return
	}

	getFeedItemsCommand := &item.GetFeedItemsCommand{FeedId: command.FeedId}
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
		_ = InternalError(w, "cannot read feed items")
		return
	}

	if err := h.itemService.ReadFeedItems(command); err != nil {
		_ = InternalError(w, "cannot read feed items")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

func (h *Router) unreadFeedItems(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	command, err := toUnreadFeedItemCommand(w, r, p)
	if err != nil {
		_ = InternalError(w, "cannot unread feed items")
		return
	}

	if err := h.itemService.UnreadFeedItems(command); err != nil {
		_ = InternalError(w, "cannot unread feed items")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

func (h *Router) deleteFeed(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	command, err := toDeleteFeedCommand(w, r, p)
	if err != nil {
		_ = InternalError(w, "cannot delete feed")
		return
	}

	cmdDeleteFeedItems, err := toDeleteFeedItemsCommand(w, r, p)
	if err := h.itemService.DeleteFeedItems(cmdDeleteFeedItems); err != nil {
		_ = InternalError(w, "cannot delete feed")
		return
	}

	if err := h.feedService.DeleteFeed(command); err != nil {
		_ = InternalError(w, "cannot delete feed")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
