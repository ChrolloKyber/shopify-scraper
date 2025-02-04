package limiter

import (
	"encoding/json"
	"net/http"

	"github.com/ChrolloKryber/shopify-scraper/models"
	"golang.org/x/time/rate"
)

func RateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	limiter := rate.NewLimiter(2, 4)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			msg := models.Message{
				Status: "Request failed",
				Body:   "API max capacity reached. Try again!",
			}
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(&msg)
			return
		} else {
			next(w, r)
		}
	})
}
