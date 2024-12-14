package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AbdoAnss/go-fantasy-pl/api"
	"github.com/AbdoAnss/go-fantasy-pl/models"
)

const (
	classicLeagueEndpoint = "/leagues-classic/%d/standings/?page_standings=%d"
	h2hLeagueEndpoint     = "/leagues-h2h-matches/league/%d/"
	maxPageCache          = 3 // Only cache first 3 pages
)

type LeagueService struct {
	client api.Client
}

func NewLeagueService(client api.Client) *LeagueService {
	return &LeagueService{
		client: client,
	}
}

func (ls *LeagueService) GetClassicLeagueStandings(id, page int) (*models.ClassicLeague, error) {
	// Only cache first few pages to prevent memory bloat
	useCache := page <= maxPageCache

	if useCache {
		cacheKey := fmt.Sprintf("classic_league_%d_page_%d", id, page)
		if cached, found := sharedCache.Get(cacheKey); found {
			if league, ok := cached.(*models.ClassicLeague); ok {
				return league, nil
			}
		}
	}

	endpoint := fmt.Sprintf(classicLeagueEndpoint, id, page)
	resp, err := ls.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get league standings: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, fmt.Errorf("league with ID %d not found", id)
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var league models.ClassicLeague
	if err := json.Unmarshal(body, &league); err != nil {
		return nil, fmt.Errorf("failed to decode league data: %w", err)
	}

	if err := ls.validateLeague(&league); err != nil {
		return nil, err
	}

	if useCache {
		cacheKey := fmt.Sprintf("classic_league_%d_page_%d", id, page)
		sharedCache.Set(cacheKey, &league, leagueCacheTTL)
	}

	return &league, nil
}

/*
func (ls *LeagueService) GetH2HLeague(id int) (*models.H2HLeague, error) {
    cacheKey := fmt.Sprintf("h2h_league_%d", id)
    if cached, found := sharedCache.Get(cacheKey); found {
        if league, ok := cached.(*models.H2HLeague); ok {
            return league, nil
        }
    }

    endpoint := fmt.Sprintf(h2hLeagueEndpoint, id)
    resp, err := ls.client.Get(endpoint)
    if err != nil {
        return nil, fmt.Errorf("failed to get H2H league: %w", err)
    }
    defer resp.Body.Close()

    switch resp.StatusCode {
    case http.StatusOK:
        // Continue processing
    case http.StatusNotFound:
        return nil, fmt.Errorf("H2H league with ID %d not found", id)
    default:
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %w", err)
    }

    var league models.H2HLeague
    if err := json.Unmarshal(body, &league); err != nil {
        return nil, fmt.Errorf("failed to decode H2H league data: %w", err)
    }

    sharedCache.Set(cacheKey, &league, leagueCacheTTL)
    return &league, nil
}
*/

func (ls *LeagueService) validateLeague(league *models.ClassicLeague) error {
	if league == nil {
		return fmt.Errorf("received nil league data")
	}
	if league.League.ID == 0 {
		return fmt.Errorf("invalid league ID")
	}
	return nil
}

func (ls *LeagueService) GetTotalPages(league *models.ClassicLeague) int {
	if league == nil || len(league.Standings.Results) == 0 {
		return 0
	}

	totalEntries := len(league.Standings.Results)
	if league.League.GetMaxEntries() > 0 {
		totalEntries = league.League.GetMaxEntries()
	}

	entriesPerPage := 50 // FPL default
	return (totalEntries + entriesPerPage - 1) / entriesPerPage
}
