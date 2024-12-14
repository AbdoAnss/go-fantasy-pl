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

	// core service
	Bootstrap *endpoints.BootstrapService

	// services
	Players  *endpoints.PlayerService
	Fixtures *endpoints.FixtureService
	Teams    *endpoints.TeamService
	Managers *endpoints.ManagerService
	Leagues  *endpoints.LeagueService
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

	// Bootstrap service
	c.Bootstrap = endpoints.NewBootstrapService(c)

	// services dependant on bootstrap:
	c.Players = endpoints.NewPlayerService(c, c.Bootstrap)
	c.Teams = endpoints.NewTeamService(c, c.Bootstrap)
	c.Managers = endpoints.NewManagerService(c, c.Bootstrap)
	// standalone services
	c.Fixtures = endpoints.NewFixtureService(c)
	c.Leagues = endpoints.NewLeagueService(c)

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
