package middleware

import (
	"fmt"
	"net/http"
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

type LimiterMock struct {
	allowResult *redis_rate.Result
	allowError  bool
}

func (l *LimiterMock) Allow(key string, limit *redis_rate.Limit) (*redis_rate.Result, error) {
	if l.allowError {
		return nil, fmt.Errorf("assertion error")
	}
	return &redis_rate.Result{
		Remaining: 5,
		Allowed:   l.allowResult.Allowed,
	}, nil
}

func TestRateLimit(t *testing.T) {
	var limiter *LimiterMock
	var res *httptest.ResponseRecorder
	var req *http.Request
	var c config.Configuration

	noop := func(w http.ResponseWriter, r *http.Request) {}

	setup := func() {
		limiter = &LimiterMock{
			allowResult: &redis_rate.Result{
				Allowed:   true,
				Remaining: 5,
			},
			allowError: false,
		}
		c = config.Configuration{
			RateLimit: config.RateLimitConfig{
				Limit: 10,
			},
		}

		res = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "http://example.com/foo", nil)
		req.Header.Set("X-Client-Id", "test")
	}

	t.Run("Should return 400 when no headers are provided", func(t *testing.T) {
		setup()
		// resetting the request so it deletes the header
		req := httptest.NewRequest("GET", "http://example.com/foo", nil)

		middleware := NewMiddleware(limiter, c)

		middleware.RateLimit(res, req, noop)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("Should return 429 when rate is exceeded", func(t *testing.T) {
		setup()

		limiter.allowResult = &redis_rate.Result{
			Allowed:   false,
			Remaining: 0,
		}

		middleware := NewMiddleware(limiter, c)

		req.Header.Set("X-Client-Id", "test")
		middleware.RateLimit(res, req, noop)
		assert.Equal(t, http.StatusTooManyRequests, res.Code)
	})

	t.Run("Should return a 500 if our rate limiter configuration is busted", func(t *testing.T) {
		setup()
		// this will return an error from the limiter allow function.
		limiter.allowError = true

		middleware := NewMiddleware(limiter, c)
		req.Header.Set("X-Client-Id", "test")
		// test
		middleware.RateLimit(res, req, noop)
		// assert a 500 intead of panic.
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})

	t.Run("Should allow the request to go through if no errors happen and is within limit", func(t *testing.T) {
		setup()

		middleware := NewMiddleware(limiter, c)

		req.Header.Set("X-Client-Id", "test")
		middleware.RateLimit(res, req, noop)
		assert.Equal(t, http.StatusOK, res.Code)
	})

}
