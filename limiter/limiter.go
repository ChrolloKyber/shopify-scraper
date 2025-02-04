package limiter

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/ChrolloKryber/shopify-scraper/models"
	"golang.org/x/time/rate"
)

func PerClientRateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {

	type Client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*Client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &Client{limiter: rate.NewLimiter(2, 4)}
		}
		clients[ip].lastSeen = time.Now()
		if !clients[ip].limiter.Allow() {
			mu.Unlock()

			msg := models.Message{
				Status: "Request failed",
				Body:   "API max capacity reached. Try again!",
			}
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(&msg)
			return
		}
		mu.Unlock()
		next(w, r)

	})
}
