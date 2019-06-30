package middleware

import (
	"log"
	"net/http"
	"time"
)

type LoggerMiddleware struct {
	Next http.Handler
}

func (l *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	l.Next.ServeHTTP(w, r)
	log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
}

func NewLoggerMiddleware(handlerToWrap http.Handler) *LoggerMiddleware {
	return &LoggerMiddleware{Next: handlerToWrap}
}
