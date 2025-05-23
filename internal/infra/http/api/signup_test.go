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

	specs := []spec{
		{
			Name:           "public registration disabled", // Renamed from "ok"
			Method:         http.MethodPost,
			Target:         "/signup",
			ReqBody:        `{"email":"test@example.com","password":"password123","confirm_password":"password123"}`, // Valid SignupCommand
			ExpectedStatus: http.StatusForbidden,
			ExpectedBody:   `{"error":{"code":403,"details":"Forbidden - public registration is disabled"}}`,
			MockFn: func(_ *mocks.ItemService, _ *mocks.FeedService, a *mocks.AuthService) {
				a.EXPECT().Signup(gomock.Any()).Return(fmt.Errorf("public registration is disabled"))
			},
		},
		{
			Name:           "invalid json payload",
			Method:         http.MethodPost,
			Target:         "/signup",
			ReqBody:        `{"email":"test@example.com","password":"password123"`,
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   `{"error":{"code":400,"details":"Bad Request - invalid request body"}}`,
			MockFn:         func(_ *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {},
		},
		{
			Name:           "missing required fields",
			Method:         http.MethodPost,
			Target:         "/signup",
			ReqBody:        `{"email":"","password":"","confirm_password":""}`, // For SignupCommand
			ExpectedStatus: http.StatusBadRequest,
			// The actual error from command.Validate() for SignupCommand would be "email and password are required"
			// or specific field errors. The current test expects a generic "invalid request body".
			// This might be because toSignupCommand sends a generic error if validation fails.
			// For now, keeping the existing ExpectedBody as per the test's original behavior for parsing/validation failures.
			ExpectedBody: `{"error":{"code":400,"details":"Bad Request - invalid request body"}}`,
			MockFn:       func(_ *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {},
		},
		// Removing the old "service returns error" test as the "public registration disabled"
		// case now covers the scenario where the service is called and returns an error.
		// The other error cases (invalid payload, missing fields, empty body) cover scenarios
		// where the service's Signup method is not even called due to request parsing/validation failures.
		{
			Name:           "empty body",
			Method:         http.MethodPost,
			Target:         "/signup",
			ReqBody:        ``,
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   `{"error":{"code":400,"details":"Bad Request - invalid request body"}}`,
			MockFn:         func(_ *mocks.ItemService, _ *mocks.FeedService, _ *mocks.AuthService) {},
		},
	}

	for _, s := range specs {
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService, authService))
	}
}
