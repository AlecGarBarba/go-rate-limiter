package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/alecGarBarba/go-rate-limiter/config"
	"github.com/urfave/negroni"
	"golang.org/x/time/rate"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Initialize Middleware
	n := negroni.Classic()

	// Rate Limiter
	limiter := rate.NewLimiter(1, 5) // adjust as needed

	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
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
