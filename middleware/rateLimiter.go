package middleware

import (
	"github.com/ZaharBorisenko/jwt-auth/helpers/JSON"
	"golang.org/x/time/rate"
	"log"
	"net/http"
)

func RateLimitMiddleware(rps int, burst int) func(handler http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				w.Header().Set("Retry-After", "60")
				log.Printf("Rate limit exceeded, Path: %s", r.URL.Path)
				JSON.WriteERROR(w, http.StatusTooManyRequests, "request too many, try again in 60 second.")
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
