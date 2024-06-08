package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/alecGarBarba/go-rate-limiter/config"
	"github.com/alecGarBarba/go-rate-limiter/middleware"
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

	middleware := middleware.NewMiddleware(limiter, config)

	// Initialize Middleware
	n := negroni.Classic()

	// Rate Limiter

	n.UseFunc(middleware.RateLimit)

	fmt.Println("API URL: ", config.APIUrl)

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(config.APIUrl)

	// Apply middleware to proxy
	n.UseHandler(proxy)

	fmt.Println("Server is running on port 8080")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", n))
}
