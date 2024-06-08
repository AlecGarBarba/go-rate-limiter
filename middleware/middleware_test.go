package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/go-redis/redis_rate/v8"
	"github.com/stretchr/testify/assert"

	"github.com/alecGarBarba/go-rate-limiter/config"
)

func TestSetRateLimitHeaders(t *testing.T) {
	res := httptest.NewRecorder()

	redisRateResult := &redis_rate.Limit{
		Rate:   10,
		Period: 1,
		Burst:  10,
	}
	// Create a mock rate limit result and configuration
	result := &redis_rate.Result{
		Limit:     redisRateResult,
		Remaining: 5,
		Allowed:   true,
	} // Replace with your mock rate limit result
	rateLimitConfig := &config.RateLimitConfig{
		Limit: 10,
	} // Replace with your mock rate limit configuration

	// Call the setRateLimitHeaders function
	setRateLimitHeaders(res, result, rateLimitConfig)

	// Assert that the response headers are set correctly

	assert.Equal(t, "10", res.Header().Get("X-RateLimit-Limit"))
	assert.Equal(t, "5", res.Header().Get("X-RateLimit-Remaining"))
	assert.Equal(t, "0s", res.Header().Get("X-RateLimit-Reset"))

}
