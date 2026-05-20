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

// LeagueService provides methods for fetching league standings and details,
// supporting both classic and head-to-head (H2H) leagues.
type LeagueService struct {
	client api.Client
}

// NewLeagueService creates a new instance of the LeagueService.
func NewLeagueService(client api.Client) *LeagueService {
	return &LeagueService{
		client: client,
	}
}

// GetClassicLeagueStandings returns the standings for a classic league by its unique ID.
// The page parameter allows for paginated access to large leagues (50 entries per page).
func (ls *LeagueService) GetClassicLeagueStandings(id, page int) (*models.ClassicLeague, error) {
	// Only cache first few pages to prevent memory bloat
	useCache := page <= maxPageCache

	if useCache {
		cacheKey := fmt.Sprintf("classic_league_%d_page_%d", id, page)
		var league models.ClassicLeague
		if sharedCache.Get(cacheKey, &league) {
			return &league, nil
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
		if err := sharedCache.Set(cacheKey, &league, leagueCacheTTL); err != nil {
			return nil, fmt.Errorf("failed to cache league standings: %w", err)
		}
	}

	return &league, nil
}

func (ls *LeagueService) validateLeague(league *models.ClassicLeague) error {
	if league == nil {
		return fmt.Errorf("received nil league data")
	}
	if league.League.ID == 0 {
		return fmt.Errorf("invalid league ID")
	}
	return nil
}

// GetTotalPages calculates the total number of pages in a classic league.
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
