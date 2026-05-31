// Package client provides the main entry point for the Fantasy Premier League SDK.
// It includes the Client struct which handles authentication, rate limiting, and
// initializes all domain-specific services.
package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/AbdoAnss/go-fantasy-pl/endpoints"
)

const (
	// baseURL is the primary entry point for the official FPL API.
	baseURL = "https://fantasy.premierleague.com/api"
	// defaultTimeout is the default HTTP timeout used if none is specified.
	defaultTimeout = 10 * time.Second
)

// Client is the main SDK client used to interact with the FPL API.
// It coordinates services, manages rate limiting, and handles HTTP communication.
type Client struct {
	httpClient *http.Client
	baseURL    string
	rateLimit  *rateLimiter
	cacheErr   error // stores errors from cache configuration to be returned by NewClient
	cacheSet   bool

	// Bootstrap provides access to core FPL data like players, teams, and gameweeks.
	Bootstrap *endpoints.BootstrapService

	// Players provides methods for fetching player-specific data and history.
	Players *endpoints.PlayerService
	// Fixtures provides methods for fetching match fixtures and results.
	Fixtures *endpoints.FixtureService
	// Teams provides methods for fetching information about Premier League teams.
	Teams *endpoints.TeamService
	// Managers provides methods for fetching manager-specific data and history.
	Managers *endpoints.ManagerService
	// Leagues provides methods for fetching league standings and details.
	Leagues *endpoints.LeagueService
}

// NewClient creates and returns a new FPL API client.
// It accepts functional options to configure the client's behavior,
// such as custom HTTP clients, timeouts, and caching strategies.
func NewClient(opts ...Option) (*Client, error) {
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

	if c.cacheErr != nil {
		return nil, fmt.Errorf("client: cache configuration failed: %w", c.cacheErr)
	}

	if !c.cacheSet {
		if err := configureDefaultCache(); err != nil {
			return nil, fmt.Errorf("client: cache configuration failed: %w", err)
		}
	}

	// Bootstrap service
	c.Bootstrap = endpoints.NewBootstrapService(c)

	// Initialize domain-specific services
	c.Players = endpoints.NewPlayerService(c, c.Bootstrap)
	c.Teams = endpoints.NewTeamService(c, c.Bootstrap)
	c.Managers = endpoints.NewManagerService(c, c.Bootstrap)
	c.Fixtures = endpoints.NewFixtureService(c)
	c.Leagues = endpoints.NewLeagueService(c)

	return c, nil
}

// Get performs a rate-limited GET request to the specified endpoint relative to the baseURL.
func (c *Client) Get(endpoint string) (*http.Response, error) {
	c.rateLimit.Wait()
	url := c.baseURL + endpoint
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	return resp, nil
}

// GetContext performs a rate-limited GET request with a context to the specified endpoint.
func (c *Client) GetContext(ctx context.Context, endpoint string) (*http.Response, error) {
	c.rateLimit.Wait()
	url := c.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	return resp, nil
}
