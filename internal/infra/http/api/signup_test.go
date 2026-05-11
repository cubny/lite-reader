package api_test

import (
	"fmt"
	"net/http"
	"testing"

	"go.uber.org/mock/gomock"

	mocks "github.com/cubny/lite-reader/internal/mocks/infra/http/api"
)

func TestRouter_signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	feedService := mocks.NewFeedService(ctrl)
	itemService := mocks.NewItemService(ctrl)
	authService := mocks.NewAuthService(ctrl)
	folderService := mocks.NewFolderService(ctrl)

	specs := []spec{
		{
			Name:           tcOK,
			Method:         http.MethodPost,
			Target:         signupPath,
			ReqBody:        `{"email":"test@example.com","password":"password123"}`,
			ExpectedStatus: http.StatusCreated,
			ExpectedBody:   ``,
			MockFn: func(_ *mocks.ItemService, _ *mocks.FeedService, a *mocks.AuthService, _ *mocks.FolderService) {
				a.EXPECT().Signup(gomock.Any()).Return(nil)
			},
		},
		{
			Name:           "invalid json payload",
			Method:         http.MethodPost,
			Target:         signupPath,
			ReqBody:        `{"email":"test@example.com","password":"password123"`,
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   respInvalidReqBody,
			MockFn:         func(_ *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService, _ *mocks.FolderService) {},
		},
		{
			Name:           "missing required fields",
			Method:         http.MethodPost,
			Target:         signupPath,
			ReqBody:        `{"email":"","password":""}`,
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   respInvalidReqBody,
			MockFn:         func(_ *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService, _ *mocks.FolderService) {},
		},
		{
			Name:           tcServiceErr,
			Method:         http.MethodPost,
			Target:         signupPath,
			ReqBody:        `{"email":"test@example.com","password":"password123"}`,
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   `{"error":{"code":400,"details":"Bad Request - user already exists"}}`,
			MockFn: func(_ *mocks.ItemService, _ *mocks.FeedService, a *mocks.AuthService, _ *mocks.FolderService) {
				a.EXPECT().Signup(gomock.Any()).Return(fmt.Errorf("user already exists"))
			},
		},
		{
			Name:           "empty body",
			Method:         http.MethodPost,
			Target:         signupPath,
			ReqBody:        ``,
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   respInvalidReqBody,
			MockFn:         func(_ *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService, _ *mocks.FolderService) {},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService, authService, folderService))
	}
}
