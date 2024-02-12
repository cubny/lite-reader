package api_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/cubny/lite-reader/internal/app/feed"
	"github.com/cubny/lite-reader/internal/app/item"
	mocks "github.com/cubny/lite-reader/internal/mocks/infra/http/api"
)

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
			MockFn: func(_ *mocks.ItemService, f *mocks.FeedService) {
				f.EXPECT().AddFeed(gomock.Any()).Return(&feed.Feed{
					ID:          1,
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
			MockFn:         func(_ *mocks.ItemService, _ *mocks.FeedService) {},
		},
		{
			Name:           "invalid url",
			ReqBody:        `{"url":"invalid.url"}`,
			ExpectedStatus: http.StatusUnprocessableEntity,
			ExpectedBody:   `{"error":{"code":422,"details":"Invalid params - invalid params"}}`,
			Method:         http.MethodPost,
			Target:         "/feeds",
			MockFn:         func(_ *mocks.ItemService, _ *mocks.FeedService) {},
		},
		{
			Name:           "service returns error",
			ReqBody:        `{"url":"http://valid.url"}`,
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - failed to add feed due to server internal error"}}`,
			Method:         http.MethodPost,
			Target:         "/feeds",
			MockFn: func(_ *mocks.ItemService, f *mocks.FeedService) {
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
			MockFn:         func(_ *mocks.ItemService, _ *mocks.FeedService) {},
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
						ID:        1,
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
						ID:        1,
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
						ID:        1,
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
			Name:           "feed service fetch items returns error",
			Method:         http.MethodPut,
			Target:         "/feeds/1/fetch",
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - cannot fetch feed items"}}`,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				f.EXPECT().FetchItems(1).Return(nil, assert.AnError)
			},
		},
		{
			Name:           "item service returns error",
			Method:         http.MethodPut,
			Target:         "/feeds/1/fetch",
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - cannot store feed items"}}`,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				f.EXPECT().FetchItems(1).Return([]*item.Item{
					{
						ID:        1,
						Title:     "title",
						Desc:      "description",
						Dir:       "dir",
						Link:      "link",
						IsNew:     true,
						Starred:   false,
						Timestamp: now,
					},
				}, nil)
				i.EXPECT().UpsertItems(gomock.Any()).Return(assert.AnError)
			},
		},
		{
			Name:           "item service get feed items returns error",
			Method:         http.MethodPut,
			Target:         "/feeds/1/fetch",
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - cannot get feed items"}}`,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				f.EXPECT().FetchItems(1).Return([]*item.Item{
					{
						ID:        1,
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
				i.EXPECT().GetFeedItems(gomock.Any()).Return(nil, assert.AnError)
			},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService))
	}
}

func TestRouter_readFeedItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	itemService := mocks.NewItemService(ctrl)
	feedService := mocks.NewFeedService(ctrl)

	specs := []spec{
		{
			Name:           "ok",
			Method:         http.MethodPost,
			Target:         "/feeds/1/read",
			ExpectedStatus: http.StatusNoContent,
			ExpectedBody:   ``,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				i.EXPECT().ReadFeedItems(gomock.Any()).Return(nil)
			},
		},
		{
			Name:           "invalid feed id",
			Method:         http.MethodPost,
			Target:         "/feeds/invalid/read",
			ExpectedStatus: http.StatusUnprocessableEntity,
			ExpectedBody:   `{"error":{"code":422,"details":"Invalid params - invalid feed id"}}`,
			MockFn:         func(i *mocks.ItemService, f *mocks.FeedService) {},
		},
		{
			Name:           "service returns error",
			Method:         http.MethodPost,
			Target:         "/feeds/1/read",
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - cannot read feed items"}}`,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				i.EXPECT().ReadFeedItems(gomock.Any()).Return(assert.AnError)
			},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService))
	}
}

func TestRouter_unreadFeedItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	itemService := mocks.NewItemService(ctrl)
	feedService := mocks.NewFeedService(ctrl)

	specs := []spec{
		{
			Name:           "ok",
			Method:         http.MethodPost,
			Target:         "/feeds/1/unread",
			ExpectedStatus: http.StatusNoContent,
			ExpectedBody:   ``,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				i.EXPECT().UnreadFeedItems(gomock.Any()).Return(nil)
			},
		},
		{
			Name:           "invalid feed id",
			Method:         http.MethodPost,
			Target:         "/feeds/invalid/unread",
			ExpectedStatus: http.StatusUnprocessableEntity,
			ExpectedBody:   `{"error":{"code":422,"details":"Invalid params - invalid feed id"}}`,
			MockFn:         func(i *mocks.ItemService, f *mocks.FeedService) {},
		},
		{
			Name:           "service returns error",
			Method:         http.MethodPost,
			Target:         "/feeds/1/unread",
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - cannot unread feed items"}}`,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				i.EXPECT().UnreadFeedItems(gomock.Any()).Return(assert.AnError)
			},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService))
	}
}

func TestRouter_DeleteFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	itemService := mocks.NewItemService(ctrl)
	feedService := mocks.NewFeedService(ctrl)

	specs := []spec{
		{
			Name:           "ok",
			Method:         http.MethodDelete,
			Target:         "/feeds/1",
			ExpectedStatus: http.StatusNoContent,
			ExpectedBody:   ``,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				f.EXPECT().DeleteFeed(gomock.Any()).Return(nil)
				i.EXPECT().DeleteFeedItems(gomock.Any()).Return(nil)
			},
		},
		{
			Name:           "invalid feed id",
			Method:         http.MethodDelete,
			Target:         "/feeds/invalid",
			ExpectedStatus: http.StatusUnprocessableEntity,
			ExpectedBody:   `{"error":{"code":422,"details":"Invalid params - invalid feed id"}}`,
			MockFn:         func(i *mocks.ItemService, f *mocks.FeedService) {},
		},
		{
			Name:           "item service returns error",
			Method:         http.MethodDelete,
			Target:         "/feeds/1",
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - cannot delete feed"}}`,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				i.EXPECT().DeleteFeedItems(gomock.Any()).Return(assert.AnError)
			},
		},
		{
			Name:           "feed service returns error",
			Method:         http.MethodDelete,
			Target:         "/feeds/1",
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - cannot delete feed"}}`,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				i.EXPECT().DeleteFeedItems(gomock.Any()).Return(nil)
				f.EXPECT().DeleteFeed(gomock.Any()).Return(assert.AnError)
			},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService))
	}
}
