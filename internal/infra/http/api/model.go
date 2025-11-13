package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"

	"github.com/cubny/lite-reader/internal/app/auth"
	"github.com/cubny/lite-reader/internal/app/feed"
	"github.com/cubny/lite-reader/internal/app/item"
	"github.com/cubny/lite-reader/internal/infra/http/api/cxutil"
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
	ID          int       `json:"id"`
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

func toAddFeedResponse(f *feed.Feed) *AddFeedResponse {
	return &AddFeedResponse{
		ID:          f.ID,
		Title:       f.Title,
		Desc:        f.Description,
		Link:        f.Link,
		URL:         f.URL,
		UpdatedAt:   f.UpdatedAt,
		Lang:        f.Lang,
		UnreadCount: f.UnreadCount,
	}
}

type ItemResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Dir       string    `json:"dir"`
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
			ID:        i.ID,
			Title:     i.Title,
			Dir:       i.Dir,
			Desc:      i.Desc,
			Link:      i.Link,
			IsNew:     i.IsNew,
			Starred:   i.Starred,
			Timestamp: i.Timestamp,
		})
	}
	return resp
}

func toAddFeedCommand(w http.ResponseWriter, r *http.Request, _ httprouter.Params) (*feed.AddFeedCommand, error) {
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
		URL:    request.URL,
		UserID: r.Context().Value(cxutil.UserIDKey).(int),
	}, nil
}

func toGetFeedItemsCommand(w http.ResponseWriter, _ *http.Request, p httprouter.Params) (*item.GetFeedItemsCommand, error) {
	feedIDString := p.ByName("id")
	feedID, err := strconv.Atoi(feedIDString)
	if err != nil {
		_ = InvalidParams(w, "invalid feed id")
		return nil, err
	}
	return &item.GetFeedItemsCommand{
		FeedID: feedID,
	}, nil
}

type UpdateItemRequest struct {
	Starred bool `json:"starred"`
	IsNew   bool `json:"is_new"`
}

func toUpdateItemCommand(w http.ResponseWriter, r *http.Request, p httprouter.Params) (*item.UpdateItemCommand, error) {
	itemIDString := p.ByName("id")
	itemID, err := strconv.Atoi(itemIDString)
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
		ID:      itemID,
		Starred: request.Starred,
		IsNew:   request.IsNew,
	}, nil
}

func toFetchFeedNewItemsCommand(w http.ResponseWriter, _ *http.Request, p httprouter.Params) (*item.FetchFeedNewItemsCommand, error) {
	feedIDString := p.ByName("id")
	feedID, err := strconv.Atoi(feedIDString)
	if err != nil {
		_ = InvalidParams(w, "invalid feed id")
		return nil, err
	}

	return &item.FetchFeedNewItemsCommand{
		FeedID: feedID,
	}, nil
}

func toReadFeedItemsCommand(w http.ResponseWriter, _ *http.Request, p httprouter.Params) (*item.ReadFeedItemsCommand, error) {
	feedIDString := p.ByName("id")
	feedID, err := strconv.Atoi(feedIDString)
	if err != nil {
		_ = InvalidParams(w, "invalid feed id")
		return nil, err
	}

	return &item.ReadFeedItemsCommand{
		FeedID: feedID,
	}, nil
}

func toUnreadFeedItemCommand(w http.ResponseWriter, _ *http.Request, p httprouter.Params) (*item.UnreadFeedItemsCommand, error) {
	feedIDString := p.ByName("id")
	feedID, err := strconv.Atoi(feedIDString)
	if err != nil {
		_ = InvalidParams(w, "invalid feed id")
		return nil, err
	}

	return &item.UnreadFeedItemsCommand{
		FeedID: feedID,
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

func toDeleteFeedCommand(w http.ResponseWriter, _ *http.Request, p httprouter.Params) (*feed.DeleteFeedCommand, error) {
	feedIDString := p.ByName("id")
	feedID, err := strconv.Atoi(feedIDString)
	if err != nil {
		_ = InvalidParams(w, "invalid feed id")
		return nil, err
	}

	return &feed.DeleteFeedCommand{
		FeedID: feedID,
	}, nil
}

func toDeleteFeedItemsCommand(w http.ResponseWriter, _ *http.Request, p httprouter.Params) (*item.DeleteFeedItemsCommand, error) {
	feedIDString := p.ByName("id")
	feedID, err := strconv.Atoi(feedIDString)
	if err != nil {
		_ = InvalidParams(w, "invalid feed id")
		return nil, err
	}

	return &item.DeleteFeedItemsCommand{
		FeedID: feedID,
	}, nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if r.Email == "" || r.Password == "" {
		return errors.New("email and password are required")
	}
	return nil
}

func toLoginCommand(w http.ResponseWriter, r *http.Request, _ httprouter.Params) (*auth.LoginCommand, error) {
	request := &LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		log.WithError(err).Error("login: failed to decode request body")
		_ = BadRequest(w, "invalid request body")
		return nil, err
	}

	if err := request.Validate(); err != nil {
		log.WithError(err).Error("login: invalid request body")
		_ = BadRequest(w, "invalid request body")
		return nil, err
	}

	return &auth.LoginCommand{
		Email:    request.Email,
		Password: request.Password,
	}, nil
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *SignupRequest) Validate() error {
	if r.Email == "" || r.Password == "" {
		return errors.New("email and password are required")
	}
	return nil
}

func toSignupCommand(w http.ResponseWriter, r *http.Request, _ httprouter.Params) (*auth.SignupCommand, error) {
	request := &SignupRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		log.WithError(err).Error("signup: failed to decode request body")
		_ = BadRequest(w, "invalid request body")
		return nil, err
	}

	if err := request.Validate(); err != nil {
		log.WithError(err).Error("signup: invalid request body")
		_ = BadRequest(w, "invalid request body")
		return nil, err
	}

	return &auth.SignupCommand{
		Email:    request.Email,
		Password: request.Password,
	}, nil
}
