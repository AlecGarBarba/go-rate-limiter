package middleware

import (
	"net/http"
	"strconv"

	"github.com/go-redis/redis_rate/v8"

	"github.com/alecGarBarba/go-rate-limiter/config"
)

type Middleware struct {
	limiter *redis_rate.Limiter
	config  config.Configuration
}

func NewMiddleware(limiter *redis_rate.Limiter, config config.Configuration) *Middleware {
	return &Middleware{
		limiter: limiter,
		config:  config,
	}
}

func (m *Middleware) RateLimit(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	clientKey := r.Header.Get("X-Client-Id")

	if clientKey == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, err := m.limiter.Allow(clientKey, redis_rate.PerSecond(m.config.RateLimit.Limit))

	setRateLimitHeaders(w, result, &m.config.RateLimit)
	if !result.Allowed {
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	next(w, r)

}

func setRateLimitHeaders(w http.ResponseWriter, result *redis_rate.Result, rateLimitConfig *config.RateLimitConfig) {
	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rateLimitConfig.Limit))
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
	w.Header().Set("X-RateLimit-Reset", result.ResetAfter.String())
}
