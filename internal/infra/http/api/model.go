package api

import (
	"encoding/json"
	"github.com/cubny/lite-reader/internal/app/item"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/url"
	"time"

	"github.com/cubny/lite-reader/internal/app/feed"
	log "github.com/sirupsen/logrus"
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
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Desc      string    `json:"desc"`
	Link      string    `json:"link"`
	URL       string    `json:"url"`
	UpdatedAt time.Time `json:"updated_at"`
	Lang      string    `json:"lang"`
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
		Id:        feed.Id,
		Title:     feed.Title,
		Desc:      feed.Description,
		Link:      feed.Link,
		URL:       feed.URL,
		UpdatedAt: feed.UpdatedAt,
		Lang:      feed.Lang,
	}
}

type ItemResponse struct {
	Id        string    `json:"id"`
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

func toGetItemsResponse(items []*item.Item) ([]*ItemResponse, error) {
	var resp []*ItemResponse
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
	return resp, nil
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

func toGetUnreadItemsCommand(w http.ResponseWriter, r *http.Request, p httprouter.Params) (*item.GetUnreadItemsCommand, error) {
	return &item.GetUnreadItemsCommand{}, nil
}

func toGetStarredItemsCommand(w http.ResponseWriter, r *http.Request, p httprouter.Params) (*item.GetStarredItemsCommand, error) {
	return &item.GetStarredItemsCommand{}, nil
}

func toGetFeedItemsCommand(w http.ResponseWriter, r *http.Request, p httprouter.Params) (*item.GetFeedItemsCommand, error) {
	return &item.GetFeedItemsCommand{
		FeedId: p.ByName("id"),
	}, nil
}

func toListFeedsCommand(w http.ResponseWriter, r *http.Request, p httprouter.Params) (*feed.ListFeedsCommand, error) {
	return &feed.ListFeedsCommand{}, nil
}
