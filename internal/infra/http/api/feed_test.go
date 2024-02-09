package api_test

import (
	"encoding/json"
	"github.com/cubny/lite-reader/internal/app/feed"
	"github.com/cubny/lite-reader/internal/app/item"
	"github.com/cubny/lite-reader/internal/infra/http/api"
	mocks "github.com/cubny/lite-reader/internal/mocks/infra/http/api"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type spec struct {
	Name           string
	ReqBody        string
	ExpectedStatus int
	ExpectedBody   string
	Method         string
	Target         string
	MockFn         func(i *mocks.ItemService, f *mocks.FeedService)
}

func (s *spec) execHTTPTestCases(i *mocks.ItemService, f *mocks.FeedService) func(t *testing.T) {
	return func(t *testing.T) {
		s.MockFn(i, f)
		handler, err := api.New(i, f)
		assert.Nil(t, err)
		s.HandlerTest(t, handler)
	}
}

// HandlerTest is a helper method to run http test cases
func (s *spec) HandlerTest(t *testing.T, h *api.Router) {
	t.Helper()

	req := httptest.NewRequest(s.Method, s.Target, strings.NewReader(s.ReqBody))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	switch {
	case s.ExpectedBody != "" && isJSON(s.ExpectedBody):
		assert.JSONEq(t, s.ExpectedBody, string(body))
	case s.ExpectedBody != "" && !isJSON(s.ExpectedBody):
		assert.Equal(t, s.ExpectedBody, strings.TrimSpace(string(body)))
	}

	assert.Equal(t, s.ExpectedStatus, resp.StatusCode)
}

func isJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func TestRouter_addFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	feedService := mocks.NewFeedService(ctrl)
	itemService := mocks.NewItemService(ctrl)
	now := time.Now()

	specs := []spec{
		{
			Name:           "ok",
			Method:         http.MethodPost,
			Target:         "/feeds",
			ReqBody:        `{"url":"http://valid.url"}`,
			ExpectedStatus: http.StatusCreated,
			ExpectedBody:   `{"id":1,"title":"title","desc":"description","link":"link","url":"url","updated_at":"` + now.Format(time.RFC3339Nano) + `","lang":"lang","unread_count":0}`,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				f.EXPECT().AddFeed(gomock.Any()).Return(&feed.Feed{
					Id:          1,
					Title:       "title",
					Description: "description",
					Link:        "link",
					URL:         "url",
					Lang:        "lang",
					UpdatedAt:   now,
					UnreadCount: 0,
				}, nil)
			},
		},
		{
			Name:           "invalid json payload",
			ReqBody:        `{"url":"http://valid.url"`,
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   `{"error":{"code":400,"details":"Bad Request - cannot decode request body"}}`,
			Method:         http.MethodPost,
			Target:         "/feeds",
			MockFn:         func(i *mocks.ItemService, f *mocks.FeedService) {},
		},
		{
			Name:           "invalid url",
			ReqBody:        `{"url":"invalid.url"}`,
			ExpectedStatus: http.StatusUnprocessableEntity,
			ExpectedBody:   `{"error":{"code":422,"details":"Invalid params - invalid params"}}`,
			Method:         http.MethodPost,
			Target:         "/feeds",
			MockFn:         func(i *mocks.ItemService, f *mocks.FeedService) {},
		},
		{
			Name:           "service returns error",
			ReqBody:        `{"url":"http://valid.url"}`,
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - failed to add feed due to server internal error"}}`,
			Method:         http.MethodPost,
			Target:         "/feeds",
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				f.EXPECT().AddFeed(gomock.Any()).Return(nil, assert.AnError)
			},
		},
		{
			Name:           "empty body",
			ReqBody:        ``,
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   `{"error":{"code":400,"details":"Bad Request - cannot decode request body"}}`,
			Method:         http.MethodPost,
			Target:         "/feeds",
			MockFn:         func(i *mocks.ItemService, f *mocks.FeedService) {},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService))
	}
}

func TestRouter_getFeedItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	itemService := mocks.NewItemService(ctrl)
	feedService := mocks.NewFeedService(ctrl)
	now := time.Now()

	specs := []spec{
		{
			Name:           "ok",
			Method:         http.MethodGet,
			Target:         "/feeds/1/items",
			ExpectedStatus: http.StatusOK,
			ExpectedBody:   `[{"id":1,"title":"title","desc":"description","link":"link","is_new":true,"starred":false,"timestamp":"` + now.Format(time.RFC3339Nano) + `"}]`,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				i.EXPECT().GetFeedItems(gomock.Any()).Return([]*item.Item{
					{
						Id:        1,
						Title:     "title",
						Desc:      "description",
						Dir:       "dir",
						Link:      "link",
						IsNew:     true,
						Starred:   false,
						Timestamp: now,
					},
				}, nil)

			},
		},
		{
			Name:           "invalid feed id",
			Method:         http.MethodGet,
			Target:         "/feeds/invalid/items",
			ExpectedStatus: http.StatusUnprocessableEntity,
			ExpectedBody:   `{"error":{"code":422,"details":"Invalid params - invalid feed id"}}`,
			MockFn:         func(i *mocks.ItemService, f *mocks.FeedService) {},
		},
		{
			Name:           "service returns error",
			Method:         http.MethodGet,
			Target:         "/feeds/1/items",
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - cannot get feed items"}}`,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				i.EXPECT().GetFeedItems(gomock.Any()).Return(nil, assert.AnError)
			},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService))
	}
}

func TestRouter_fetchFeedNewItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	itemService := mocks.NewItemService(ctrl)
	feedService := mocks.NewFeedService(ctrl)
	now := time.Now()

	specs := []spec{
		{
			Name:           "ok",
			Method:         http.MethodPut,
			Target:         "/feeds/1/fetch",
			ExpectedStatus: http.StatusOK,
			ExpectedBody:   `[{"id":1,"title":"title","desc":"description","link":"link","is_new":false,"starred":true,"timestamp":"` + now.Format(time.RFC3339Nano) + `"}]`,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				f.EXPECT().FetchItems(1).Return([]*item.Item{
					{
						Id:        1,
						Title:     "title",
						Desc:      "description",
						Dir:       "dir",
						Link:      "link",
						IsNew:     true,
						Starred:   false,
						Timestamp: now,
					},
				}, nil)
				i.EXPECT().UpsertItems(gomock.Any()).Return(nil)
				i.EXPECT().GetFeedItems(gomock.Any()).Return([]*item.Item{
					{
						Id:        1,
						Title:     "title",
						Desc:      "description",
						Dir:       "dir",
						Link:      "link",
						IsNew:     false,
						Starred:   true,
						Timestamp: now,
					},
				}, nil)
			},
		},
		{
			Name:           "invalid feed id",
			Method:         http.MethodPut,
			Target:         "/feeds/invalid/fetch",
			ExpectedStatus: http.StatusUnprocessableEntity,
			ExpectedBody:   `{"error":{"code":422,"details":"Invalid params - invalid feed id"}}`,
			MockFn:         func(i *mocks.ItemService, f *mocks.FeedService) {},
		},
		{
			Name:           "service returns error",
			Method:         http.MethodPut,
			Target:         "/feeds/1/fetch",
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - cannot fetch feed items"}}`,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				f.EXPECT().FetchItems(1).Return(nil, assert.AnError)
			},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService))
	}
}
