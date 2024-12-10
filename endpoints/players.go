package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/AbdoAnss/go-fantasy-pl/api"
	"github.com/AbdoAnss/go-fantasy-pl/internal/cache"
	"github.com/AbdoAnss/go-fantasy-pl/models"
)

const (
	playersEndpoint       = "/bootstrap-static/"
	playerDetailsEndpoint = "/element-summary/%d/"
)

type PlayerService struct {
	client api.Client
}

func NewPlayerService(client api.Client) *PlayerService {
	return &PlayerService{
		client: client,
	}
}

// TODO:
// Centralized Cache with Namespacing:
// Use a single cache instance and differentiate keys using endpoint-specific prefixes.

var (
	defaultCacheTTL = 10 * time.Minute
	playersCache    = cache.NewCache()
)

func init() {
	// Start cleanup task to run every 5 minute
	playersCache.StartCleanupTask(5 * time.Minute)
}

func (ps *PlayerService) GetAllPlayers() ([]models.Player, error) {
	if cached, found := playersCache.Get("players"); found {
		if players, ok := cached.([]models.Player); ok {
			return players, nil
		}
	}

	resp, err := ps.client.Get(playersEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get players: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		Elements []models.Player `json:"elements"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode players: %w", err)
	}

	playersCache.Set("players", response.Elements, defaultCacheTTL)

	return response.Elements, nil
}

func (ps *PlayerService) GetPlayer(id int) (*models.Player, error) {
	players, err := ps.GetAllPlayers()
	if err != nil {
		return nil, err
	}

	for _, p := range players {
		if p.ID == id {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("player with ID %d not found", id)
}

func (ps *PlayerService) GetPlayerHistory(id int) (*models.PlayerHistory, error) {
	cacheKey := fmt.Sprintf("player_history_%d", id)
	if cached, found := playersCache.Get(cacheKey); found {
		if history, ok := cached.(*models.PlayerHistory); ok {
			return history, nil
		}
	}

	endpoint := fmt.Sprintf(playerDetailsEndpoint, id)
	resp, err := ps.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error fetching player history: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("player not found: %d", id)
		}
		return nil, fmt.Errorf("error fetching player history: received status code %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var history models.PlayerHistory
	if err := json.Unmarshal(bodyBytes, &history); err != nil {
		return nil, fmt.Errorf("error decoding player history: %w", err)
	}

	if history.History == nil {
		return nil, fmt.Errorf("history is nil in response for player ID %d", id)
	}

	playersCache.Set(cacheKey, &history, defaultCacheTTL)

	return &history, nil
}
