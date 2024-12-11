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
	sharedCache       = cache.NewCache()
	teamsCacheTTL     = 24 * time.Hour   // Teams rarely change
	playersCacheTTL   = 10 * time.Minute // Players update more frequently (injuries, etc)
	fixturesCacheTTL  = 10 * time.Minute
	gameweeksCacheTTL = 3 * time.Minute // Gameweeks status might change more often
	settingsCacheTTL  = 24 * time.Hour  // Game settings rarely change
	managerCacheTTL   = 5 * time.Minute // Managers data updates frequently
)

func init() {
	sharedCache.StartCleanupTask(5 * time.Minute)
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
	if cached, found := sharedCache.Get(cacheKey); found {
		if teams, ok := cached.([]models.Team); ok {
			return teams, nil
		}
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
	if cached, found := sharedCache.Get(cacheKey); found {
		if players, ok := cached.([]models.Player); ok {
			return players, nil
		}
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
	if cached, found := sharedCache.Get(cacheKey); found {
		if gw, ok := cached.([]models.GameWeek); ok {
			return gw, nil
		}
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
	if cached, found := sharedCache.Get(cacheKey); found {
		if gw, ok := cached.(int); ok {
			return gw, nil
		}
	}
	gameweeks, err := bs.GetGameWeeks()
	if err != nil {
		return 0, fmt.Errorf("failed to get gameweeks: %w", err)
	}

	for _, gw := range gameweeks {
		if gw.IsCurrent {
			sharedCache.Set(cacheKey, gw.ID, gameweeksCacheTTL)
			return gw.ID, nil
		}
	}

	return 0, fmt.Errorf("failed to find current gameweek")
}

func (bs *BootstrapService) GetSettings() (*models.GameSettings, error) {
	const cacheKey = "settings"
	if cached, found := sharedCache.Get(cacheKey); found {
		if settings, ok := cached.(*models.GameSettings); ok {
			return settings, nil
		}
	}

	data, err := bs.fetchBootstrapData()
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	settings := &data.Settings
	sharedCache.Set(cacheKey, settings, settingsCacheTTL)
	return settings, nil
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
