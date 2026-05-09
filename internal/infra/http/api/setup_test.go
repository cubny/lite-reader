package api_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/cubny/lite-reader/internal/app/auth"
	mocks "github.com/cubny/lite-reader/internal/mocks/infra/http/api"
)

func testSession() *auth.Session {
	return &auth.Session{
		ID:           "test-session",
		UserID:       1,
		AccessToken:  "test",
		RefreshToken: "test-refresh",
		ExpiresAt:    time.Now().Add(time.Hour),
		CreatedAt:    time.Now(),
	}
}

func TestRouter_setup(t *testing.T) {
	ctrl := gomock.NewController(t)
	feedService := mocks.NewFeedService(ctrl)
	itemService := mocks.NewItemService(ctrl)
	authService := mocks.NewAuthService(ctrl)

	specs := []spec{
		{
			Name:           tcOK,
			Method:         http.MethodPost,
			Target:         setupPath,
			ReqBody:        `{"email":"admin@example.com","password":"password123","confirm_password":"password123","allow_signup":false}`,
			ExpectedStatus: http.StatusCreated,
			MockFn: func(_ *mocks.ItemService, _ *mocks.FeedService, a *mocks.AuthService) {
				a.EXPECT().Setup(gomock.Any()).Return(nil)
			},
		},
		{
			Name:           "setup already completed",
			Method:         http.MethodPost,
			Target:         setupPath,
			ReqBody:        `{"email":"admin@example.com","password":"password123","confirm_password":"password123"}`,
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   `{"error":{"code":400,"details":"Bad Request - setup has already been completed"}}`,
			MockFn: func(_ *mocks.ItemService, _ *mocks.FeedService, a *mocks.AuthService) {
				a.EXPECT().Setup(gomock.Any()).Return(errors.New("setup has already been completed"))
			},
		},
		{
			Name:           "invalid json",
			Method:         http.MethodPost,
			Target:         setupPath,
			ReqBody:        `{invalid`,
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   respInvalidReqBody,
			MockFn:         func(_ *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {},
		},
		{
			Name:           "validation error - passwords don't match",
			Method:         http.MethodPost,
			Target:         setupPath,
			ReqBody:        `{"email":"admin@example.com","password":"password123","confirm_password":"different"}`,
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   `{"error":{"code":400,"details":"Bad Request - password and confirm password do not match"}}`,
			MockFn:         func(_ *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {},
		},
		{
			Name:           "validation error - missing fields",
			Method:         http.MethodPost,
			Target:         setupPath,
			ReqBody:        `{"email":"","password":"","confirm_password":""}`,
			ExpectedStatus: http.StatusBadRequest,
			MockFn:         func(_ *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService, authService))
	}
}

func TestRouter_needsSetup(t *testing.T) {
	ctrl := gomock.NewController(t)
	feedService := mocks.NewFeedService(ctrl)
	itemService := mocks.NewItemService(ctrl)
	authService := mocks.NewAuthService(ctrl)

	specs := []spec{
		{
			Name:           "needs setup",
			Method:         http.MethodGet,
			Target:         setupStatusPath,
			ExpectedStatus: http.StatusOK,
			ExpectedBody:   `{"allow_signup":false,"needs_setup":true}`,
			MockFn: func(_ *mocks.ItemService, _ *mocks.FeedService, a *mocks.AuthService) {
				a.EXPECT().NeedsSetup().Return(true, nil)
				a.EXPECT().IsSignupAllowed().Return(false)
			},
		},
		{
			Name:           "setup already done",
			Method:         http.MethodGet,
			Target:         setupStatusPath,
			ExpectedStatus: http.StatusOK,
			ExpectedBody:   `{"allow_signup":true,"needs_setup":false}`,
			MockFn: func(_ *mocks.ItemService, _ *mocks.FeedService, a *mocks.AuthService) {
				a.EXPECT().NeedsSetup().Return(false, nil)
				a.EXPECT().IsSignupAllowed().Return(true)
			},
		},
		{
			Name:           "error checking setup",
			Method:         http.MethodGet,
			Target:         setupStatusPath,
			ExpectedStatus: http.StatusInternalServerError,
			MockFn: func(_ *mocks.ItemService, _ *mocks.FeedService, a *mocks.AuthService) {
				a.EXPECT().NeedsSetup().Return(false, errors.New("db error"))
			},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService, authService))
	}
}

func TestRouter_getSettings(t *testing.T) {
	ctrl := gomock.NewController(t)
	feedService := mocks.NewFeedService(ctrl)
	itemService := mocks.NewItemService(ctrl)
	authService := mocks.NewAuthService(ctrl)

	specs := []spec{
		{
			Name:           "get settings",
			Method:         http.MethodGet,
			Target:         settingsPath,
			ExpectedStatus: http.StatusOK,
			ExpectedBody:   respAllowSignupTrue,
			MockFn: func(_ *mocks.ItemService, _ *mocks.FeedService, a *mocks.AuthService) {
				a.EXPECT().GetSession("test").Return(testSession(), nil)
				a.EXPECT().IsSignupAllowed().Return(true)
			},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService, authService))
	}
}

func TestRouter_updateSettings(t *testing.T) {
	ctrl := gomock.NewController(t)
	feedService := mocks.NewFeedService(ctrl)
	itemService := mocks.NewItemService(ctrl)
	authService := mocks.NewAuthService(ctrl)

	specs := []spec{
		{
			Name:           "update settings - allow signup",
			Method:         http.MethodPut,
			Target:         settingsPath,
			ReqBody:        respAllowSignupTrue,
			ExpectedStatus: http.StatusOK,
			ExpectedBody:   respAllowSignupTrue,
			MockFn: func(_ *mocks.ItemService, _ *mocks.FeedService, a *mocks.AuthService) {
				a.EXPECT().GetSession("test").Return(testSession(), nil)
				a.EXPECT().SetAllowSignup(true).Return(nil)
			},
		},
		{
			Name:           "update settings - disable signup",
			Method:         http.MethodPut,
			Target:         settingsPath,
			ReqBody:        respAllowSignupFalse,
			ExpectedStatus: http.StatusOK,
			ExpectedBody:   respAllowSignupFalse,
			MockFn: func(_ *mocks.ItemService, _ *mocks.FeedService, a *mocks.AuthService) {
				a.EXPECT().GetSession("test").Return(testSession(), nil)
				a.EXPECT().SetAllowSignup(false).Return(nil)
			},
		},
		{
			Name:           "update settings - error",
			Method:         http.MethodPut,
			Target:         settingsPath,
			ReqBody:        respAllowSignupTrue,
			ExpectedStatus: http.StatusInternalServerError,
			MockFn: func(_ *mocks.ItemService, _ *mocks.FeedService, a *mocks.AuthService) {
				a.EXPECT().GetSession("test").Return(testSession(), nil)
				a.EXPECT().SetAllowSignup(true).Return(errors.New("db error"))
			},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService, authService))
	}
}
