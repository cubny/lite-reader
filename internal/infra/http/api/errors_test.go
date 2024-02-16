package api_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cubny/lite-reader/internal/infra/http/api"
)

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()

	err := api.BadRequest(w, "test")
	require.NoError(t, err)

	assert.Equal(t, w.Code, http.StatusBadRequest)

	expectedBody := `{"error":{"code":400, "details":"Bad Request - test"}}`
	assertBody(t, expectedBody, w.Body)
}

func TestInternalError(t *testing.T) {
	w := httptest.NewRecorder()

	err := api.InternalError(w, "test")
	require.NoError(t, err)

	assert.Equal(t, w.Code, http.StatusInternalServerError)

	expectedBody := `{"error":{"code":500, "details":"Internal error - test"}}`
	assertBody(t, expectedBody, w.Body)
}

func TestInvalidParams(t *testing.T) {
	w := httptest.NewRecorder()

	err := api.InvalidParams(w, "test")
	require.NoError(t, err)

	assert.Equal(t, w.Code, http.StatusUnprocessableEntity)

	expectedBody := `{"error":{"code":422, "details":"Invalid params - test"}}`
	assertBody(t, expectedBody, w.Body)
}

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()

	err := api.NotFound(w, "test")
	require.NoError(t, err)

	assert.Equal(t, w.Code, http.StatusNotFound)

	expectedBody := `{"error":{"code":404, "details":"Not found - test"}}`
	assertBody(t, expectedBody, w.Body)
}

func assertBody(t *testing.T, expectedBody string, actualBody *bytes.Buffer) {
	t.Helper()

	body, err := io.ReadAll(actualBody)
	assert.Nil(t, err)
	assert.JSONEq(t, expectedBody, string(body))
}
