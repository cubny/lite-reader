package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/cubny/lite-reader/internal/app/auth"
	"github.com/cubny/lite-reader/internal/infra/http/api/cxutil"

	"github.com/julienschmidt/httprouter"
)

type AuthService interface {
	GetSession(token string) (*auth.Session, error)
}

// HandleFunc is a method type that represents Middleware HandleFunc function.
type HandleFunc func(httprouter.Handle) httprouter.Handle

// Chain represents helper struct that is able to wrap httprouter.Handle with chain.
type Chain struct {
	middlewares []HandleFunc
}

// NewChain create a chain out of middlewares.
func NewChain(middlewares ...HandleFunc) Chain {
	return Chain{middlewares: middlewares}
}

// Wrap the handler with Chain.
func (chain Chain) Wrap(handler httprouter.Handle) httprouter.Handle {
	result := handler

	for i := len(chain.middlewares) - 1; i >= 0; i-- {
		middleware := chain.middlewares[i]

		if result == nil {
			result = middleware(handler)
			continue
		}

		result = middleware(result)
	}

	return result
}

// ContentTypeJSON for the HTTP response.
func ContentTypeJSON(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r, ps)
	}
}

func AuthMiddleware(authService AuthService) HandleFunc {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

			w.Header().Set("Content-Type", "application/json")
			log.Printf("Request path: %s", r.URL.Path)

			// Skip auth check for login and signup endpoints
			if r.URL.Path == "/login" || r.URL.Path == "/signup" {
				next(w, r, ps)
				return
			}

			// For all other endpoints, require auth
			token, err := extractBearerToken(r)
			if err != nil {
				writeUnauthorizedResponse(w, err.Error())
				return
			}

			session, err := authService.GetSession(token)
			if err != nil || session == nil {
				writeUnauthorizedResponse(w, "Unauthorized - Invalid token")
				return
			}

			// Add userID to context
			ctx := context.WithValue(r.Context(), cxutil.UserIDKey, session.UserID)
			r = r.WithContext(ctx)

			next(w, r, ps)
		}
	}
}

func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no token provided")
	}

	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", errors.New("invalid token format")
	}

	return authHeader[7:], nil
}

func writeUnauthorizedResponse(w http.ResponseWriter, details string) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"details": "unauthorized - " + details,
		},
	})
}
