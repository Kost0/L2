package middleware

import (
	"log"
	"net/http"
	"time"
)

func Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		handler.ServeHTTP(w, r)

		duration := time.Since(start)

		log.Printf("%s %s %s", r.Method, r.URL.Path, duration)
	})
}
