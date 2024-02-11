package api_test

import (
	"github.com/cubny/lite-reader/internal/app/item"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"

	"go.uber.org/mock/gomock"

	mocks "github.com/cubny/lite-reader/internal/mocks/infra/http/api"
)

func TestRouter_updateItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	feedService := mocks.NewFeedService(ctrl)
	itemService := mocks.NewItemService(ctrl)

	specs := []spec{
		{
			Name:           "update item",
			ReqBody:        `{"starred": true, "is_new": true}`,
			ExpectedStatus: http.StatusNoContent,
			ExpectedBody:   ``,
			Method:         http.MethodPut,
			Target:         "/items/1",
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				i.EXPECT().UpdateItem(gomock.Any()).Return(nil)
			},
		},
		{
			Name:           "invalid params",
			ReqBody:        `{"starred": true, "is_new": true}`,
			ExpectedStatus: http.StatusUnprocessableEntity,
			ExpectedBody:   `{"error":{"code":422, "details":"Invalid params - invalid item id"}}`,
			Method:         http.MethodPut,
			Target:         "/items/a",
			MockFn:         func(i *mocks.ItemService, f *mocks.FeedService) {},
		},
		{
			Name:           "bad request",
			MockFn:         func(i *mocks.ItemService, f *mocks.FeedService) {},
			ReqBody:        ``,
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   `{"error":{"code":400, "details":"Bad Request - cannot decode request body"}}`,
			Method:         http.MethodPut,
			Target:         "/items/1",
		},
		{
			Name:           "service returns error",
			ReqBody:        `{"id": "1", "starred": true, "is_new": true}`,
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500, "details":"Internal error - cannot update item"}}`,
			Method:         http.MethodPut,
			Target:         "/items/1",
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				i.EXPECT().UpdateItem(gomock.Any()).Return(assert.AnError)
			},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService))
	}
}

func TestRouter_getStarredItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	feedService := mocks.NewFeedService(ctrl)
	itemService := mocks.NewItemService(ctrl)

	specs := []spec{
		{
			Name:           "get starred items",
			ExpectedStatus: http.StatusOK,
			ExpectedBody:   `[]`,
			Method:         http.MethodGet,
			Target:         "/items/starred",
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				i.EXPECT().GetStarredItems().Return([]*item.Item{}, nil)
			},
		},
		{
			Name:           "service returns error",
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500, "details":"Internal error - cannot get unread items"}}`,
			Method:         http.MethodGet,
			Target:         "/items/starred",
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				i.EXPECT().GetStarredItems().Return(nil, assert.AnError)
			},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService))
	}
}

func TestRouter_getUnreadItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	feedService := mocks.NewFeedService(ctrl)
	itemService := mocks.NewItemService(ctrl)

	specs := []spec{
		{
			Name:           "get unread items",
			ExpectedStatus: http.StatusOK,
			ExpectedBody:   `[]`,
			Method:         http.MethodGet,
			Target:         "/items/unread",
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				i.EXPECT().GetUnreadItems().Return([]*item.Item{}, nil)
			},
		},
		{
			Name:           "service returns error",
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500, "details":"Internal error - cannot get unread items"}}`,
			Method:         http.MethodGet,
			Target:         "/items/unread",
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				i.EXPECT().GetUnreadItems().Return(nil, assert.AnError)
			},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService))
	}
}

func TestRouter_getUnreadItemsCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	feedService := mocks.NewFeedService(ctrl)
	itemService := mocks.NewItemService(ctrl)

	specs := []spec{
		{
			Name:           "get unread items count",
			ExpectedStatus: http.StatusOK,
			ExpectedBody:   `{"count":0}`,
			Method:         http.MethodGet,
			Target:         "/items/unread/count",
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				i.EXPECT().GetUnreadItemsCount().Return(0, nil)
			},
		},
		{
			Name:           "service returns error",
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500, "details":"Internal error - cannot get unread items"}}`,
			Method:         http.MethodGet,
			Target:         "/items/unread/count",
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				i.EXPECT().GetUnreadItemsCount().Return(0, assert.AnError)
			},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService))
	}
}

func TestRouter_getStarredItemsCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	feedService := mocks.NewFeedService(ctrl)
	itemService := mocks.NewItemService(ctrl)

	specs := []spec{
		{
			Name:           "get starred items count",
			ExpectedStatus: http.StatusOK,
			ExpectedBody:   `{"count":0}`,
			Method:         http.MethodGet,
			Target:         "/items/starred/count",
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				i.EXPECT().GetStarredItemsCount().Return(0, nil)
			},
		},
		{
			Name:           "service returns error",
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody:   `{"error":{"code":500, "details":"Internal error - cannot get unread items"}}`,
			Method:         http.MethodGet,
			Target:         "/items/starred/count",
			MockFn: func(i *mocks.ItemService, f *mocks.FeedService) {
				i.EXPECT().GetStarredItemsCount().Return(0, assert.AnError)
			},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService))
	}
}
