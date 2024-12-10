package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AbdoAnss/go-fantasy-pl/endpoints"
)

const (
	baseURL        = "https://fantasy.premierleague.com/api"
	defaultTimeout = 10 * time.Second
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	rateLimit  *rateLimiter

	// services
	Players  *endpoints.PlayerService
	Fixtures *endpoints.FixtureService
	Teams    *endpoints.TeamService
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
			Transport: &http.Transport{
				MaxIdleConns:          10,
				IdleConnTimeout:       30 * time.Second,
				DisableCompression:    false,
				DisableKeepAlives:     false,
				MaxConnsPerHost:       10,
				ResponseHeaderTimeout: 10 * time.Second,
			},
		},
		baseURL:   baseURL,
		rateLimit: newRateLimiter(50, time.Minute),
	}

	for _, opt := range opts {
		opt(c)
	}

	// services

	c.Players = endpoints.NewPlayerService(c)
	c.Fixtures = endpoints.NewFixtureService(c)
	c.Teams = endpoints.NewTeamService(c)

	return c
}

func (c *Client) Get(endpoint string) (*http.Response, error) {
	c.rateLimit.Wait()
	url := c.baseURL + endpoint
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	return resp, nil
}
