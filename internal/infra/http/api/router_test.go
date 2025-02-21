package api_test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cubny/lite-reader/internal/infra/http/api"
	mocks "github.com/cubny/lite-reader/internal/mocks/infra/http/api"
)

type spec struct {
	Name           string
	ReqBody        string
	ExpectedStatus int
	ExpectedBody   string
	Method         string
	Target         string
	AuthToken      string
	MockFn         func(i *mocks.ItemService, f *mocks.FeedService, a *mocks.AuthService)
}

func (s *spec) execHTTPTestCases(i *mocks.ItemService, f *mocks.FeedService, a *mocks.AuthService) func(t *testing.T) {
	return func(t *testing.T) {
		s.MockFn(i, f, a)
		s.AuthToken = "test"
		handler, err := api.New(i, f, a)
		assert.Nil(t, err)
		s.HandlerTest(t, handler)
	}
}

// HandlerTest is a helper method to run http test cases
func (s *spec) HandlerTest(t *testing.T, h *api.Router) {
	t.Helper()

	req := httptest.NewRequest(s.Method, s.Target, strings.NewReader(s.ReqBody))
	if s.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+s.AuthToken)
	}

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
