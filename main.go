package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

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

	// ping redis to ensure connection

	_, err = rdb.Ping().Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	
	}

	limiter := redis_rate.NewLimiter(rdb)

	// Initialize Middleware
	n := negroni.Classic()

	// Rate Limiter

	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		// TODO: make this configurable via env variables as well.
		result, err := limiter.Allow("api", redis_rate.PerMinute(5))

		if !result.Allowed {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		} else if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// TODO:  we need to add the X-headers to inform client about remaining tries and reset time

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
