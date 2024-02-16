package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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
