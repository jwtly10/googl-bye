package middleware

import (
	"net/http"
)

type Middleware interface {
	BeforeNext(next http.Handler) http.Handler
}

// Chain applies multiple middleware to a http.Handler.
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware.BeforeNext(h)
	}
	return h
}
