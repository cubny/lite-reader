package api

import (
	"encoding/json"
	"github.com/cubny/lite-reader/internal/app/feed"
	"github.com/cubny/lite-reader/internal/app/item"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type AddFeedRequest struct {
	URL string `json:"url"`
}

func (r *AddFeedRequest) Validate() error {
	if _, err := url.ParseRequestURI(r.URL); err != nil {
		return err
	}
	return nil
}

type AddFeedResponse struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Desc        string    `json:"desc"`
	Link        string    `json:"link"`
	URL         string    `json:"url"`
	UpdatedAt   time.Time `json:"updated_at"`
	Lang        string    `json:"lang"`
	UnreadCount int       `json:"unread_count"`
}

type ListFeedResponse []*AddFeedResponse

func toListFeedResponse(feeds []*feed.Feed) *ListFeedResponse {
	resp := make(ListFeedResponse, 0)
	for _, f := range feeds {
		resp = append(resp, toAddFeedResponse(f))
	}
	return &resp
}

func toAddFeedResponse(feed *feed.Feed) *AddFeedResponse {
	return &AddFeedResponse{
		Id:          feed.Id,
		Title:       feed.Title,
		Desc:        feed.Description,
		Link:        feed.Link,
		URL:         feed.URL,
		UpdatedAt:   feed.UpdatedAt,
		Lang:        feed.Lang,
		UnreadCount: feed.UnreadCount,
	}
}

type ItemResponse struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Desc      string    `json:"desc"`
	Link      string    `json:"link"`
	IsNew     bool      `json:"is_new"`
	Starred   bool      `json:"starred"`
	Timestamp time.Time `json:"timestamp"`
}

type GetUnreadItemsResponse struct {
	Items []*ItemResponse `json:"items"`
}

func toGetItemsResponse(items []*item.Item) []*ItemResponse {
	resp := make([]*ItemResponse, 0)
	for _, i := range items {
		resp = append(resp, &ItemResponse{
			Id:        i.Id,
			Title:     i.Title,
			Desc:      i.Desc,
			Link:      i.Link,
			IsNew:     i.IsNew,
			Starred:   i.Starred,
			Timestamp: i.Timestamp,
		})
	}
	return resp
}

func toAddFeedCommand(w http.ResponseWriter, r *http.Request, p httprouter.Params) (*feed.AddFeedCommand, error) {
	request := &AddFeedRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		log.WithError(err).Errorf("toAddFeedCommand: decoder %s", err)
		_ = BadRequest(w, "cannot decode request body")
		return nil, err
	}

	if err := request.Validate(); err != nil {
		log.WithError(err).Errorf("toAddFeedCommand: validator %s", err)
		_ = InvalidParams(w, "invalid params")
		return nil, err
	}

	return &feed.AddFeedCommand{
		URL: request.URL,
	}, nil
}

func toGetFeedItemsCommand(w http.ResponseWriter, r *http.Request, p httprouter.Params) (*item.GetFeedItemsCommand, error) {
	feedIdString := p.ByName("id")
	feedId, err := strconv.Atoi(feedIdString)
	if err != nil {
		_ = InvalidParams(w, "invalid feed id")
		return nil, err
	}
	return &item.GetFeedItemsCommand{
		FeedId: feedId,
	}, nil
}

type UpdateItemRequest struct {
	Starred bool `json:"starred"`
	IsNew   bool `json:"is_new"`
}

func toUpdateItemCommand(w http.ResponseWriter, r *http.Request, p httprouter.Params) (*item.UpdateItemCommand, error) {
	itemIdString := p.ByName("id")
	itemId, err := strconv.Atoi(itemIdString)
	if err != nil {
		_ = InvalidParams(w, "invalid item id")
		return nil, err
	}

	request := &UpdateItemRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		log.WithError(err).Errorf("toUpdateItemCommand: decoder %s", err)
		_ = BadRequest(w, "cannot decode request body")
		return nil, err
	}

	return &item.UpdateItemCommand{
		Id:      itemId,
		Starred: request.Starred,
		IsNew:   request.IsNew,
	}, nil
}

func toFetchFeedNewItemsCommand(w http.ResponseWriter, r *http.Request, p httprouter.Params) (*item.FetchFeedNewItemsCommand, error) {
	feedIdString := p.ByName("id")
	feedId, err := strconv.Atoi(feedIdString)
	if err != nil {
		_ = InvalidParams(w, "invalid feed id")
		return nil, err
	}

	return &item.FetchFeedNewItemsCommand{
		FeedId: feedId,
	}, nil
}

func toReadFeedItemsCommand(w http.ResponseWriter, r *http.Request, p httprouter.Params) (*item.ReadFeedItemsCommand, error) {
	feedIdString := p.ByName("id")
	feedId, err := strconv.Atoi(feedIdString)
	if err != nil {
		_ = InvalidParams(w, "invalid feed id")
		return nil, err
	}

	return &item.ReadFeedItemsCommand{
		FeedId: feedId,
	}, nil
}

func toUnreadFeedItemCommand(w http.ResponseWriter, r *http.Request, p httprouter.Params) (*item.UnreadFeedItemsCommand, error) {
	feedIdString := p.ByName("id")
	feedId, err := strconv.Atoi(feedIdString)
	if err != nil {
		_ = InvalidParams(w, "invalid feed id")
		return nil, err
	}

	return &item.UnreadFeedItemsCommand{
		FeedId: feedId,
	}, nil
}

type GetItemsCountResponse struct {
	Count int `json:"count"`
}

func toGetItemsCountResponse(count int) *GetItemsCountResponse {
	return &GetItemsCountResponse{
		Count: count,
	}
}

func toDeleteFeedCommand(w http.ResponseWriter, r *http.Request, p httprouter.Params) (*feed.DeleteFeedCommand, error) {
	feedIdString := p.ByName("id")
	feedId, err := strconv.Atoi(feedIdString)
	if err != nil {
		_ = InvalidParams(w, "invalid feed id")
		return nil, err
	}

	return &feed.DeleteFeedCommand{
		FeedId: feedId,
	}, nil
}

func toDeleteFeedItemsCommand(w http.ResponseWriter, r *http.Request, p httprouter.Params) (*item.DeleteFeedItemsCommand, error) {
	feedIdString := p.ByName("id")
	feedId, err := strconv.Atoi(feedIdString)
	if err != nil {
		_ = InvalidParams(w, "invalid feed id")
		return nil, err
	}

	return &item.DeleteFeedItemsCommand{
		FeedId: feedId,
	}, nil
}
