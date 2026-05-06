package api_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/cubny/lite-reader/internal/app/auth"
	"github.com/cubny/lite-reader/internal/app/feed"
	"github.com/cubny/lite-reader/internal/app/item"
	mocks "github.com/cubny/lite-reader/internal/mocks/infra/http/api"
)

func TestRouter_addFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	feedService := mocks.NewFeedService(ctrl)
	itemService := mocks.NewItemService(ctrl)
	authService := mocks.NewAuthService(ctrl)
	authService.EXPECT().GetSession(gomock.Any()).Return(&auth.Session{}, nil).AnyTimes()
	now := time.Now()

	specs := []spec{
		{
			Name:           tcOK,
			Method:         http.MethodPost,
			Target:         feedsPath,
			ReqBody:        `{"url":"http://valid.url"}`,
			ExpectedStatus: http.StatusCreated,
			ExpectedBody: `{"id":1,"title":"title","desc":"description","link":"link","url":"url","updated_at":"` +
				now.Format(time.RFC3339Nano) + `","lang":"lang","unread_count":0}`,
			MockFn: func(_ *mocks.ItemService, f *mocks.FeedService, _ *mocks.AuthService) {
				f.EXPECT().AddFeed(gomock.Any()).Return(&feed.Feed{
					ID:          1,
					Title:       feedTitle,
					Description: feedDesc,
					Link:        feedLink,
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
			Target:         feedsPath,
			MockFn:         func(_ *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {},
		},
		{
			Name:           "invalid url",
			ReqBody:        `{"url":"invalid.url"}`,
			ExpectedStatus: http.StatusUnprocessableEntity,
			ExpectedBody:   `{"error":{"code":422,"details":"Invalid params - invalid params"}}`,
			Method:         http.MethodPost,
			Target:         feedsPath,
			MockFn:         func(_ *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {},
		},
		{
			Name:           tcServiceErr,
			ReqBody:        `{"url":"http://valid.url"}`,
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - failed to add feed due to server internal error"}}`,
			Method:         http.MethodPost,
			Target:         feedsPath,
			MockFn: func(_ *mocks.ItemService, f *mocks.FeedService, _ *mocks.AuthService) {
				f.EXPECT().AddFeed(gomock.Any()).Return(nil, assert.AnError)
			},
		},
		{
			Name:           "empty body",
			ReqBody:        ``,
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   `{"error":{"code":400,"details":"Bad Request - cannot decode request body"}}`,
			Method:         http.MethodPost,
			Target:         feedsPath,
			MockFn:         func(_ *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService, authService))
	}
}

func TestRouter_getFeedItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	itemService := mocks.NewItemService(ctrl)
	feedService := mocks.NewFeedService(ctrl)
	authService := mocks.NewAuthService(ctrl)
	authService.EXPECT().GetSession(gomock.Any()).Return(&auth.Session{}, nil).AnyTimes()
	now := time.Now()

	specs := []spec{
		{
			Name:           tcOK,
			Method:         http.MethodGet,
			Target:         "/feeds/1/items",
			ExpectedStatus: http.StatusOK,
			ExpectedBody: `[{"id":1,"title":"title","dir":"dir","desc":"description","link":"link","is_new":true,"starred":false,"timestamp":"` +
				now.Format(time.RFC3339Nano) + `"}]`,
			MockFn: func(i *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {
				i.EXPECT().GetFeedItems(gomock.Any()).Return([]*item.Item{
					{
						ID:        1,
						Title:     feedTitle,
						Desc:      feedDesc,
						Dir:       feedDir,
						Link:      feedLink,
						IsNew:     true,
						Starred:   false,
						Timestamp: now,
					},
				}, nil)
			},
		},
		{
			Name:           tcInvalidFeedID,
			Method:         http.MethodGet,
			Target:         "/feeds/invalid/items",
			ExpectedStatus: http.StatusUnprocessableEntity,
			ExpectedBody:   respInvalidFeedID,
			MockFn:         func(_ *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {},
		},
		{
			Name:           tcServiceErr,
			Method:         http.MethodGet,
			Target:         "/feeds/1/items",
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - cannot get feed items"}}`,
			MockFn: func(i *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {
				i.EXPECT().GetFeedItems(gomock.Any()).Return(nil, assert.AnError)
			},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService, authService))
	}
}

func TestRouter_fetchFeedNewItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	itemService := mocks.NewItemService(ctrl)
	feedService := mocks.NewFeedService(ctrl)
	authService := mocks.NewAuthService(ctrl)
	authService.EXPECT().GetSession(gomock.Any()).Return(&auth.Session{}, nil).AnyTimes()
	now := time.Now()

	specs := []spec{
		{
			Name:           tcOK,
			Method:         http.MethodPut,
			Target:         feedFetchPath,
			ExpectedStatus: http.StatusOK,
			ExpectedBody: `[{"id":1,"title":"title","dir":"dir","desc":"description","link":"link","is_new":false,"starred":true,"timestamp":"` +
				now.Format(time.RFC3339Nano) + `"}]`,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService, _ *mocks.AuthService) {
				f.EXPECT().FetchItems(1).Return([]*item.Item{
					{
						ID:        1,
						Title:     feedTitle,
						Desc:      feedDesc,
						Dir:       feedDir,
						Link:      feedLink,
						IsNew:     true,
						Starred:   false,
						Timestamp: now,
					},
				}, nil)
				i.EXPECT().UpsertItems(gomock.Any()).Return(nil)
				i.EXPECT().GetFeedItems(gomock.Any()).Return([]*item.Item{
					{
						ID:        1,
						Title:     feedTitle,
						Desc:      feedDesc,
						Dir:       feedDir,
						Link:      feedLink,
						IsNew:     false,
						Starred:   true,
						Timestamp: now,
					},
				}, nil)
			},
		},
		{
			Name:           tcInvalidFeedID,
			Method:         http.MethodPut,
			Target:         "/feeds/invalid/fetch",
			ExpectedStatus: http.StatusUnprocessableEntity,
			ExpectedBody:   respInvalidFeedID,
			MockFn:         func(_ *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {},
		},
		{
			Name:           "feed service fetch items returns error",
			Method:         http.MethodPut,
			Target:         feedFetchPath,
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - cannot fetch feed items"}}`,
			MockFn: func(_ *mocks.ItemService, f *mocks.FeedService, _ *mocks.AuthService) {
				f.EXPECT().FetchItems(1).Return(nil, assert.AnError)
			},
		},
		{
			Name:           "item service returns error",
			Method:         http.MethodPut,
			Target:         feedFetchPath,
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - cannot store feed items"}}`,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService, _ *mocks.AuthService) {
				f.EXPECT().FetchItems(1).Return([]*item.Item{
					{
						ID:        1,
						Title:     feedTitle,
						Desc:      feedDesc,
						Dir:       feedDir,
						Link:      feedLink,
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
			Target:         feedFetchPath,
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - cannot get feed items"}}`,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService, _ *mocks.AuthService) {
				f.EXPECT().FetchItems(1).Return([]*item.Item{
					{
						ID:        1,
						Title:     feedTitle,
						Desc:      feedDesc,
						Dir:       feedDir,
						Link:      feedLink,
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
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService, authService))
	}
}

func TestRouter_readFeedItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	itemService := mocks.NewItemService(ctrl)
	feedService := mocks.NewFeedService(ctrl)
	authService := mocks.NewAuthService(ctrl)
	authService.EXPECT().GetSession(gomock.Any()).Return(&auth.Session{}, nil).AnyTimes()

	specs := []spec{
		{
			Name:           tcOK,
			Method:         http.MethodPost,
			Target:         "/feeds/1/read",
			ExpectedStatus: http.StatusNoContent,
			ExpectedBody:   ``,
			MockFn: func(i *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {
				i.EXPECT().ReadFeedItems(gomock.Any()).Return(nil)
			},
		},
		{
			Name:           tcInvalidFeedID,
			Method:         http.MethodPost,
			Target:         "/feeds/invalid/read",
			ExpectedStatus: http.StatusUnprocessableEntity,
			ExpectedBody:   respInvalidFeedID,
			MockFn:         func(_ *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {},
		},
		{
			Name:           tcServiceErr,
			Method:         http.MethodPost,
			Target:         "/feeds/1/read",
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - cannot read feed items"}}`,
			MockFn: func(i *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {
				i.EXPECT().ReadFeedItems(gomock.Any()).Return(assert.AnError)
			},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService, authService))
	}
}

func TestRouter_unreadFeedItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	itemService := mocks.NewItemService(ctrl)
	feedService := mocks.NewFeedService(ctrl)
	authService := mocks.NewAuthService(ctrl)
	authService.EXPECT().GetSession(gomock.Any()).Return(&auth.Session{}, nil).AnyTimes()

	specs := []spec{
		{
			Name:           tcOK,
			Method:         http.MethodPost,
			Target:         "/feeds/1/unread",
			ExpectedStatus: http.StatusNoContent,
			ExpectedBody:   ``,
			MockFn: func(i *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {
				i.EXPECT().UnreadFeedItems(gomock.Any()).Return(nil)
			},
		},
		{
			Name:           tcInvalidFeedID,
			Method:         http.MethodPost,
			Target:         "/feeds/invalid/unread",
			ExpectedStatus: http.StatusUnprocessableEntity,
			ExpectedBody:   respInvalidFeedID,
			MockFn:         func(_ *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {},
		},
		{
			Name:           tcServiceErr,
			Method:         http.MethodPost,
			Target:         "/feeds/1/unread",
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - cannot unread feed items"}}`,
			MockFn: func(i *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {
				i.EXPECT().UnreadFeedItems(gomock.Any()).Return(assert.AnError)
			},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService, authService))
	}
}

func TestRouter_DeleteFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	itemService := mocks.NewItemService(ctrl)
	feedService := mocks.NewFeedService(ctrl)
	authService := mocks.NewAuthService(ctrl)
	authService.EXPECT().GetSession(gomock.Any()).Return(&auth.Session{}, nil).AnyTimes()

	specs := []spec{
		{
			Name:           tcOK,
			Method:         http.MethodDelete,
			Target:         feed1Path,
			ExpectedStatus: http.StatusNoContent,
			ExpectedBody:   ``,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService, _ *mocks.AuthService) {
				f.EXPECT().DeleteFeed(gomock.Any()).Return(nil)
				i.EXPECT().DeleteFeedItems(gomock.Any()).Return(nil)
			},
		},
		{
			Name:           tcInvalidFeedID,
			Method:         http.MethodDelete,
			Target:         "/feeds/invalid",
			ExpectedStatus: http.StatusUnprocessableEntity,
			ExpectedBody:   respInvalidFeedID,
			MockFn:         func(_ *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {},
		},
		{
			Name:           "item service returns error",
			Method:         http.MethodDelete,
			Target:         feed1Path,
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - cannot delete feed"}}`,
			MockFn: func(i *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {
				i.EXPECT().DeleteFeedItems(gomock.Any()).Return(assert.AnError)
			},
		},
		{
			Name:           "feed service returns error",
			Method:         http.MethodDelete,
			Target:         feed1Path,
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500,"details":"Internal error - cannot delete feed"}}`,
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService, _ *mocks.AuthService) {
				i.EXPECT().DeleteFeedItems(gomock.Any()).Return(nil)
				f.EXPECT().DeleteFeed(gomock.Any()).Return(assert.AnError)
			},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService, authService))
	}
}
