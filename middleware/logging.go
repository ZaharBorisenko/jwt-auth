package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(`{
			"timestamp": "%s"
			"method": "%s"
			"path": "%s"
			"duration": "%s"
			"ip": "%s"
			"userAgent": "%s"
			"referer": "%s"
		}`,
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			time.Since(start),
			r.RemoteAddr,
			r.UserAgent(),
			r.Referer(),
		)
	})
}
