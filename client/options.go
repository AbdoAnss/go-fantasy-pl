package client

import (
	"net/http"
	"time"

	"github.com/AbdoAnss/go-fantasy-pl/endpoints"
	"github.com/AbdoAnss/go-fantasy-pl/internal/cache"
)

type Option func(*Client)

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = timeout
	}
}

func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

func WithRateLimit(requests int, interval time.Duration) Option {
	return func(c *Client) {
		c.rateLimit = newRateLimiter(requests, interval)
	}
}

// WithRedisCache configures the client to use a Redis-backed distributed cache.
// This replaces the default in-memory cache, enabling safe use across multiple
// instances (e.g. in a horizontally-scaled microservice deployment).
// NewClient returns an error if the Redis server is unreachable.
//
// Example:
//
//	c, err := client.NewClient(
//	    client.WithRedisCache(cache.RedisOptions{
//	        Addr:      "redis:6379",
//	        Password:  "secret",
//	        DB:        0,
//	        KeyPrefix: "fpl",
//	    }),
//	)
func WithRedisCache(opts cache.RedisOptions) Option {
	return func(c *Client) {
		rc, err := cache.NewRedisCache(opts)
		if err != nil {
			// Surface the error at build-time by storing it; callers can check
			// via a build-time validate or accept the default memory cache.
			// We do not panic to keep the library non-fatal on mis-configuration.
			c.cacheErr = err
			return
		}
		endpoints.SetSharedCache(rc)
	}
}
