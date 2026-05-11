package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"

	"github.com/cubny/lite-reader/internal/app/feed"
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
	}
}
func (h *Router) listFeeds(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := r.Context().Value(cxutil.UserIDKey).(int)
	log.Infof("listFeeds: userID %d", userID)
	resp, err := h.feedService.ListFeeds(userID)
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

func (h *Router) patchFeed(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	feedIDStr := p.ByName("id")
	feedID, err := strconv.Atoi(feedIDStr)
	if err != nil {
		_ = InvalidParams(w, "invalid feed id")
		return
	}
	userID := r.Context().Value(cxutil.UserIDKey).(int)

	request := &UpdateFeedRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		_ = BadRequest(w, "cannot decode request body")
		return
	}

	if request.FolderID != nil || request.UnsetFolder {
		var fid *int
		if !request.UnsetFolder {
			fid = request.FolderID
		}
		if err := h.feedService.MoveFeed(&feed.MoveFeedCommand{
			FeedID: feedID, UserID: userID, FolderID: fid,
		}); err != nil {
			log.WithError(err).Error("patchFeed: move")
			_ = InternalError(w, "cannot move feed")
			return
		}
	}

	if request.Position != nil {
		if err := h.feedService.ReorderFeed(&feed.ReorderFeedCommand{
			FeedID: feedID, UserID: userID, Position: *request.Position,
		}); err != nil {
			log.WithError(err).Error("patchFeed: reorder")
			_ = InternalError(w, "cannot reorder feed")
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Router) bulkMoveFeeds(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := r.Context().Value(cxutil.UserIDKey).(int)
	request := &BulkMoveFeedsRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		_ = BadRequest(w, "cannot decode request body")
		return
	}
	var fid *int
	if !request.UnsetFolder {
		fid = request.FolderID
	}
	if err := h.feedService.BulkMoveFeeds(&feed.BulkMoveFeedsCommand{
		FeedIDs: request.FeedIDs, UserID: userID, FolderID: fid,
	}); err != nil {
		log.WithError(err).Error("bulkMoveFeeds: service")
		_ = InternalError(w, "cannot move feeds")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Router) deleteFeed(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	command, err := toDeleteFeedCommand(w, r, p)
	if err != nil {
		return
	}

	cmdDeleteFeedItems, err := toDeleteFeedItemsCommand(w, r, p)
	if err != nil {
		return
	}

	if err := h.itemService.DeleteFeedItems(cmdDeleteFeedItems); err != nil {
		_ = InternalError(w, "cannot delete feed")
		return
	}

	if err := h.feedService.DeleteFeed(command); err != nil {
		_ = InternalError(w, "cannot delete feed")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
