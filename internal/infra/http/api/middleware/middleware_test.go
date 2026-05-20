package middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"

	"github.com/cubny/lite-reader/internal/app/auth"
	"github.com/cubny/lite-reader/internal/infra/http/api/cxutil"
)

type fakeAuthService struct {
	session *auth.Session
	err     error
}

func (f *fakeAuthService) GetSession(string) (*auth.Session, error) {
	return f.session, f.err
}

const bearerToken = "Bearer tok"

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		svc        *fakeAuthService
		wantStatus int
		wantNext   bool
	}{
		{
			name:       "missing token",
			authHeader: "",
			svc:        &fakeAuthService{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "session not found returns 401",
			authHeader: bearerToken,
			svc:        &fakeAuthService{err: sql.ErrNoRows},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "nil session returns 401",
			authHeader: bearerToken,
			svc:        &fakeAuthService{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "transient DB error returns 500, not 401",
			authHeader: bearerToken,
			svc:        &fakeAuthService{err: errors.New("database is locked (SQLITE_BUSY)")},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "valid session calls next",
			authHeader: bearerToken,
			svc:        &fakeAuthService{session: &auth.Session{UserID: 42}},
			wantStatus: http.StatusOK,
			wantNext:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false
			next := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				nextCalled = true
				if uid, ok := r.Context().Value(cxutil.UserIDKey).(int); !ok || uid != 42 {
					t.Errorf("expected userID 42 in context, got %v (ok=%v)", uid, ok)
				}
				w.WriteHeader(http.StatusOK)
			}

			handler := AuthMiddleware(tt.svc)(next)

			req := httptest.NewRequest(http.MethodGet, "/feeds", http.NoBody)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()
			handler(rec, req, nil)

			if rec.Code != tt.wantStatus {
				t.Errorf("status: got %d, want %d", rec.Code, tt.wantStatus)
			}
			if nextCalled != tt.wantNext {
				t.Errorf("next called: got %v, want %v", nextCalled, tt.wantNext)
			}
		})
	}
}

func TestAuthMiddlewareSkipsPublicPaths(t *testing.T) {
	for _, path := range []string{"/login", "/signup", "/setup", "/setup/status"} {
		t.Run(path, func(t *testing.T) {
			nextCalled := false
			next := func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			}
			// No auth header and an erroring service: public paths must skip
			// the auth check entirely.
			handler := AuthMiddleware(&fakeAuthService{err: sql.ErrNoRows})(next)

			req := httptest.NewRequest(http.MethodPost, path, http.NoBody)
			rec := httptest.NewRecorder()
			handler(rec, req, nil)

			if !nextCalled {
				t.Errorf("%s: expected next to be called", path)
			}
		})
	}
}
