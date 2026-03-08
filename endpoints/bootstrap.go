package endpoints

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/AbdoAnss/go-fantasy-pl/api"
	"github.com/AbdoAnss/go-fantasy-pl/internal/cache"
	"github.com/AbdoAnss/go-fantasy-pl/models"
)

const (
	bootstrapEndpoint = "/bootstrap-static/"
)

var (
	sharedCache       cache.Cache = cache.NewMemoryCache()
	teamsCacheTTL                 = 24 * time.Hour   // Teams rarely change
	playersCacheTTL               = 10 * time.Minute // Players update more frequently (injuries, etc)
	fixturesCacheTTL              = 10 * time.Minute
	gameweeksCacheTTL             = 3 * time.Minute // Gameweeks status might change more often
	settingsCacheTTL              = 24 * time.Hour  // Game settings rarely change
	managerCacheTTL               = 5 * time.Minute // Managers data updates frequently
	leagueCacheTTL                = 5 * time.Minute // Leagues update frequently
)

func init() {
	if mc, ok := sharedCache.(*cache.MemoryCache); ok {
		mc.StartCleanupTask(5 * time.Minute)
	}
}

// SetSharedCache replaces the cache used by all endpoint services.
// Call this before creating any client instances, typically via client.WithRedisCache.
func SetSharedCache(c cache.Cache) {
	sharedCache = c
}

type Response struct {
	Teams    []models.Team       `json:"teams"`
	Elements []models.Player     `json:"elements"`
	Events   []models.GameWeek   `json:"events"`
	Settings models.GameSettings `json:"game_settings"`
}

type BootstrapService struct {
	client api.Client
}

func NewBootstrapService(client api.Client) *BootstrapService {
	return &BootstrapService{
		client: client,
	}
}

func (bs *BootstrapService) GetTeams() ([]models.Team, error) {
	const cacheKey = "teams"
	var teams []models.Team
	if sharedCache.Get(cacheKey, &teams) {
		return teams, nil
	}

	data, err := bs.fetchBootstrapData()
	if err != nil {
		return nil, fmt.Errorf("failed to get teams: %w", err)
	}

	sharedCache.Set(cacheKey, data.Teams, teamsCacheTTL)
	return data.Teams, nil
}

func (bs *BootstrapService) GetPlayers() ([]models.Player, error) {
	const cacheKey = "players"
	var players []models.Player
	if sharedCache.Get(cacheKey, &players) {
		return players, nil
	}

	data, err := bs.fetchBootstrapData()
	if err != nil {
		return nil, fmt.Errorf("failed to get players: %w", err)
	}

	sharedCache.Set(cacheKey, data.Elements, playersCacheTTL)
	return data.Elements, nil
}

func (bs *BootstrapService) GetGameWeeks() ([]models.GameWeek, error) {
	const cacheKey = "gameweeks"
	var gw []models.GameWeek
	if sharedCache.Get(cacheKey, &gw) {
		return gw, nil
	}

	data, err := bs.fetchBootstrapData()
	if err != nil {
		return nil, fmt.Errorf("failed to get gameweeks: %w", err)
	}

	sharedCache.Set(cacheKey, data.Events, gameweeksCacheTTL)
	return data.Events, nil
}

func (bs *BootstrapService) GetCurrentGameWeek() (int, error) {
	const cacheKey = "current_gameweek"
	var gw int
	if sharedCache.Get(cacheKey, &gw) {
		return gw, nil
	}

	gameweeks, err := bs.GetGameWeeks()
	if err != nil {
		return 0, fmt.Errorf("failed to get gameweeks: %w", err)
	}

	for _, g := range gameweeks {
		if g.IsCurrent {
			sharedCache.Set(cacheKey, g.ID, gameweeksCacheTTL)
			return g.ID, nil
		}
	}

	return 0, fmt.Errorf("failed to find current gameweek")
}

func (bs *BootstrapService) GetSettings() (*models.GameSettings, error) {
	const cacheKey = "settings"
	var settings models.GameSettings
	if sharedCache.Get(cacheKey, &settings) {
		return &settings, nil
	}

	data, err := bs.fetchBootstrapData()
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	sharedCache.Set(cacheKey, data.Settings, settingsCacheTTL)
	return &data.Settings, nil
}

func (bs *BootstrapService) fetchBootstrapData() (*Response, error) {
	resp, err := bs.client.Get(bootstrapEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get bootstrap data: %w", err)
	}
	defer resp.Body.Close()

	var bootstrapResp Response
	if err := json.NewDecoder(resp.Body).Decode(&bootstrapResp); err != nil {
		return nil, fmt.Errorf("failed to decode bootstrap data: %w", err)
	}

	return &bootstrapResp, nil
}
