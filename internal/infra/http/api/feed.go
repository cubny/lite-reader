package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// addFeed is the handler for
// swagger:route POST /feeds setFeedsRequest
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
		_ = InternalError(w, "failed to parse request body")
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
	command, err := toListFeedsCommand(w, r, p)
	if err != nil {
		_ = InternalError(w, "cannot list feeds")
		return
	}

	resp, err := h.feedService.ListFeeds(command)
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

// getFeed is the handler for
// swagger:route GET /feeds/{feed_id} getFeedRequest
//
// Responds how much time remains until the feed's webhook is shot.
//
// Responses:
//
//	200: getFeed
//	404: notFoundError
//	422: invalidParams
//	500: serverError
func (h *Router) getFeed(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	feedID := p.ByName("id")
	if feedID == "" {
		_ = InvalidParams(w, "invalid param: id is empty")
		return
	}

	resp := "1 hour"

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.WithError(err).Errorf("getFeed: encoder %s", err)
		_ = InternalError(w, "cannot encode response")
		return
	}
}

func (h *Router) getUnreadItems(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	command, err := toGetUnreadItemsCommand(w, r, p)
	if err != nil {
		_ = InternalError(w, "cannot get unread items")
		return
	}

	items, err := h.itemService.GetUnreadItems(command)
	if err != nil {
		_ = InternalError(w, "cannot get unread items")
		return
	}

	resp, err := toGetItemsResponse(items)
	if err != nil {
		_ = InternalError(w, "cannot get unread items")
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.WithError(err).Errorf("getUnreadItems: encoder %s", err)
		_ = InternalError(w, "cannot encode response")
		return
	}
}

func (h *Router) getStarredItems(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	command, err := toGetStarredItemsCommand(w, r, p)
	if err != nil {
		_ = InternalError(w, "cannot get unread items")
		return
	}

	items, err := h.itemService.GetStarredItems(command)
	if err != nil {
		_ = InternalError(w, "cannot get unread items")
		return
	}

	resp, err := toGetItemsResponse(items)
	if err != nil {
		_ = InternalError(w, "cannot get unread items")
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.WithError(err).Errorf("getStarredItems: encoder %s", err)
		_ = InternalError(w, "cannot encode response")
		return
	}
}

func (h *Router) getFeedItems(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	command, err := toGetFeedItemsCommand(w, r, p)
	if err != nil {
		_ = InternalError(w, "cannot get feed items")
		return
	}

	items, err := h.itemService.GetFeedItems(command)
	if err != nil {
		_ = InternalError(w, "cannot get feed items")
		return
	}

	resp, err := toGetItemsResponse(items)
	if err != nil {
		_ = InternalError(w, "cannot get unread items")
		return
	}
	log.Infof("getFeedItems: resp %v", resp)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.WithError(err).Errorf("getFeedItems: encoder %s", err)
		_ = InternalError(w, "cannot encode response")
		return
	}
}
