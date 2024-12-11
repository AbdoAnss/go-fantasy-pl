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
	playerDetailsEndpoint = "/element-summary/%d/"
)

type PlayerService struct {
	client           api.Client
	bootstrapService *BootstrapService
}

func NewPlayerService(client api.Client, bootstrap *BootstrapService) *PlayerService {
	return &PlayerService{
		client:           client,
		bootstrapService: bootstrap,
	}
}

func (ps *PlayerService) GetAllPlayers() ([]models.Player, error) {
	return ps.bootstrapService.GetPlayers()
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
	if cached, found := sharedCache.Get(cacheKey); found {
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

	sharedCache.Set(cacheKey, &history, playersCacheTTL)
	return &history, nil
}
