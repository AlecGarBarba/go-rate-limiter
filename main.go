package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/alecGarBarba/go-rate-limiter/config"
	"github.com/go-redis/redis/v7"
	"github.com/go-redis/redis_rate/v8"
	"github.com/urfave/negroni"
)

func main() {

	config, err := config.LoadConfig()

	if err != nil {

		log.Fatalf("Error loading config: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	// ping redis to ensure connection is healthy

	_, err = rdb.Ping().Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	limiter := redis_rate.NewLimiter(rdb)

	// Initialize Middleware
	n := negroni.Classic()

	// Rate Limiter

	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		clientKey := r.Header.Get("X-Client-Id")

		if clientKey == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		result, err := limiter.Allow(clientKey, redis_rate.PerSecond(config.RateLimit.Limit))

		setRateLimitHeaders(w, result, &config.RateLimit)
		if !result.Allowed {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		} else if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		next(w, r)
	})

	fmt.Println("API URL: ", config.APIUrl)

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(config.APIUrl)

	// Apply middleware to proxy
	n.UseHandler(proxy)

	fmt.Println("Server is running on port 8080")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", n))
}

func setRateLimitHeaders(w http.ResponseWriter, result *redis_rate.Result, rateLimitConfig *config.RateLimitConfig) {
	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rateLimitConfig.Limit))
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
	w.Header().Set("X-RateLimit-Reset", result.ResetAfter.String())
}
