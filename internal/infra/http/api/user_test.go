package api_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cubny/lite-reader/internal/app/auth"
	mocks "github.com/cubny/lite-reader/internal/mocks/infra/http/api"
	"go.uber.org/mock/gomock"
)

// Assuming 'spec' struct and 'execHTTPTestCases' are defined in a common test_helper.go or similar.
// If not, I'll need to define them or adapt. For now, I'll write the tests as if they exist,
// similar to signup_test.go.
// If these are not in a shared helper, the `TestRouter_createUser` function would be more verbose,
// setting up the router and request/response recorder for each sub-test.

func TestRouter_createUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	feedService := mocks.NewFeedService(ctrl)  // Mocked, but not used directly by /api/users
	itemService := mocks.NewItemService(ctrl)  // Mocked, but not used directly by /api/users
	authService := mocks.NewAuthService(ctrl)

	validToken := "valid-test-token"
	mockSession := &auth.Session{
		ID:           1,
		UserID:       1,
		AccessToken:  validToken,
		RefreshToken: "valid-refresh-token",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
		CreatedAt:    time.Now(),
	}

	specs := []spec{
		{
			Name:           "success - user created",
			Method:         http.MethodPost,
			Target:         "/api/users",
			ReqBody:        `{"email":"newuser@example.com","password":"password123"}`,
			AuthToken:      validToken, // Add this field to spec or handle in MockFn
			ExpectedStatus: http.StatusCreated,
			ExpectedBody:   ``, // Expect empty body on 201
			MockFn: func(_ *mocks.ItemService, _ *mocks.FeedService, a *mocks.AuthService) {
				// Mock for middleware
				a.EXPECT().GetSession(validToken).Return(mockSession, nil).AnyTimes()
				// Mock for handler
				a.EXPECT().CreateUser(gomock.Any()).DoAndReturn(func(cmd *auth.CreateUserCommand) error {
					if cmd.Email == "newuser@example.com" && cmd.Password == "password123" {
						return nil
					}
					return fmt.Errorf("unexpected command for CreateUser mock: %+v", cmd)
				})
			},
		},
		{
			Name:           "unauthenticated - no token",
			Method:         http.MethodPost,
			Target:         "/api/users",
			ReqBody:        `{"email":"newuser@example.com","password":"password123"}`,
			AuthToken:      "", // No token
			ExpectedStatus: http.StatusUnauthorized,
			ExpectedBody:   `{"error":{"code":401,"details":"unauthorized - no token provided"}}`,
			MockFn: func(_ *mocks.ItemService, _ *mocks.FeedService, a *mocks.AuthService) {
				// GetSession will be called by middleware, but CreateUser should not be.
				// The actual GetSession mock for an empty token might not be needed if the middleware handles it first.
				// If AuthMiddleware calls GetSession(""), then:
				// a.EXPECT().GetSession("").Return(nil, fmt.Errorf("no token provided")).AnyTimes()
			},
		},
		{
			Name:           "unauthenticated - invalid token",
			Method:         http.MethodPost,
			Target:         "/api/users",
			ReqBody:        `{"email":"newuser@example.com","password":"password123"}`,
			AuthToken:      "invalid-token",
			ExpectedStatus: http.StatusUnauthorized,
			ExpectedBody:   `{"error":{"code":401,"details":"unauthorized - Unauthorized - Invalid token"}}`,
			MockFn: func(_ *mocks.ItemService, _ *mocks.FeedService, a *mocks.AuthService) {
				a.EXPECT().GetSession("invalid-token").Return(nil, fmt.Errorf("session not found")).AnyTimes()
			},
		},
		{
			Name:           "invalid input - email already exists",
			Method:         http.MethodPost,
			Target:         "/api/users",
			ReqBody:        `{"email":"existing@example.com","password":"password123"}`,
			AuthToken:      validToken,
			ExpectedStatus: http.StatusBadRequest, // As per createUser handler
			ExpectedBody:   `{"error":{"code":400,"details":"Bad Request - email already registered"}}`,
			MockFn: func(_ *mocks.ItemService, _ *mocks.FeedService, a *mocks.AuthService) {
				a.EXPECT().GetSession(validToken).Return(mockSession, nil).AnyTimes()
				a.EXPECT().CreateUser(gomock.Any()).DoAndReturn(func(cmd *auth.CreateUserCommand) error {
					if cmd.Email == "existing@example.com" {
						return fmt.Errorf("email already registered")
					}
					return fmt.Errorf("unexpected command for CreateUser mock: %+v", cmd)
				})
			},
		},
		{
			Name:           "bad request - malformed json",
			Method:         http.MethodPost,
			Target:         "/api/users",
			ReqBody:        `{"email":"newuser@example.com","password":"password123"`, // Malformed
			AuthToken:      validToken,
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   `{"error":{"code":400,"details":"Bad Request - invalid request body: unexpected EOF"}}`, // Error from json.Decode
			MockFn: func(_ *mocks.ItemService, _ *mocks.FeedService, a *mocks.AuthService) {
				// GetSession will be called by middleware. CreateUser should not be called.
				a.EXPECT().GetSession(validToken).Return(mockSession, nil).AnyTimes()
			},
		},
		{
			Name:           "bad request - missing required field (e.g. email)",
			Method:         http.MethodPost,
			Target:         "/api/users",
			ReqBody:        `{"password":"password123"}`, // Missing email
			AuthToken:      validToken,
			ExpectedStatus: http.StatusBadRequest,
			// This error comes from command.Validate() inside toCreateUserCommand
			ExpectedBody:   `{"error":{"code":400,"details":"Bad Request - email and password are required"}}`,
			MockFn: func(_ *mocks.ItemService, _ *mocks.FeedService, a *mocks.AuthService) {
				a.EXPECT().GetSession(validToken).Return(mockSession, nil).AnyTimes()
			},
		},
	}

	for _, s := range specs {
		// This assumes execHTTPTestCases is defined in another file in the same package (api_test)
		// e.g. in a test_helper.go or in signup_test.go if it's not refactored out.
		// If not, this part needs to be implemented here.
		t.Run(s.Name, s.execHTTPTestCases(itemService, feedService, authService))
	}
}
