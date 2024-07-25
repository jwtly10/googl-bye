package middleware

import (
	"net/http"
	"time"

	"github.com/jwtly10/googl-bye/internal/common"
)

type RequestLoggerMiddleware struct {
	log common.Logger
}

func NewRequestLoggerMiddleware(log common.Logger) *RequestLoggerMiddleware {
	return &RequestLoggerMiddleware{
		log: log,
	}
}

// BeforeNext is a method that implements the Middleware interface using a pointer receiver.
func (rmw *RequestLoggerMiddleware) BeforeNext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		rmw.log.Infof("Method: %s, Path: %s, Duration: %s", r.Method, r.URL.Path, duration)
	})
}
