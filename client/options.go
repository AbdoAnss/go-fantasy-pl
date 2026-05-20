package client

import (
	"net/http"
	"time"

	"github.com/AbdoAnss/go-fantasy-pl/endpoints"
	"github.com/AbdoAnss/go-fantasy-pl/internal/cache"
)

// Option is a functional option for configuring the Client.
type Option func(*Client)

// RedisOptions configures the SDK's Redis-backed cache.
type RedisOptions = cache.RedisOptions

// WithHTTPClient sets a custom http.Client for the SDK to use.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout sets the timeout for all API requests.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = timeout
	}
}

// WithBaseURL overrides the default FPL API base URL.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithRateLimit configures the client's internal rate limiter.
func WithRateLimit(requests int, interval time.Duration) Option {
	return func(c *Client) {
		c.rateLimit = newRateLimiter(requests, interval)
	}
}

// WithRedisCache configures the client to use a Redis-backed distributed cache.
// This overrides the default cache selection, enabling shared state across
// multiple instances (e.g., in a horizontally-scaled microservice deployment).
// NewClient will return an error if the Redis server is unreachable.
func WithRedisCache(opts RedisOptions) Option {
	return func(c *Client) {
		c.cacheSet = true
		rc, err := cache.NewRedisCache(opts)
		if err != nil {
			c.cacheErr = err
			return
		}
		endpoints.SetSharedCache(rc)
	}
}

// WithMemoryCache forces the SDK to use the in-memory cache backend.
func WithMemoryCache() Option {
	return func(c *Client) {
		c.cacheSet = true
		endpoints.SetSharedCache(cache.NewMemoryCache())
	}
}
