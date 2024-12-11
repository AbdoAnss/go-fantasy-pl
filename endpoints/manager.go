// Package endpoints provides access to the Fantasy Premier League API
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
	managerDetailsEndpoint       = "/entry/%d/"
	managerHistoryEndpoint       = "/entry/%d/history"
	managerGameWeekPicksEndpoint = "/entry/%d/event/%d/picks/"
)

type ManagerService struct {
	client           api.Client
	bootstrapService *BootstrapService
}

func NewManagerService(client api.Client, bootstrap *BootstrapService) *ManagerService {
	return &ManagerService{
		client:           client,
		bootstrapService: bootstrap,
	}
}

func (ms *ManagerService) validateManager(manager *models.Manager) error {
	if manager == nil {
		return fmt.Errorf("received nil manager data")
	}
	if manager.ID == nil {
		return fmt.Errorf("manager ID is missing")
	}
	return nil
}

func (ms *ManagerService) GetManager(id int) (*models.Manager, error) {
	cacheKey := fmt.Sprintf("manager_%d", id)
	if cached, found := sharedCache.Get(cacheKey); found {
		if manager, ok := cached.(*models.Manager); ok {
			return manager, nil
		}
	}

	endpoint := fmt.Sprintf(managerDetailsEndpoint, id)
	resp, err := ms.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get manager data: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, fmt.Errorf("manager with ID %d not found", id)
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var manager models.Manager
	if err := json.Unmarshal(body, &manager); err != nil {
		return nil, fmt.Errorf("failed to decode manager data: %w", err)
	}

	if err := ms.validateManager(&manager); err != nil {
		return nil, err
	}

	sharedCache.Set(cacheKey, &manager, managerCacheTTL)

	return &manager, nil
}

func (ms *ManagerService) GetCurrentTeam(managerID int) (*models.ManagerTeam, error) {
	cacheKey := fmt.Sprintf("manager_team_%d", managerID)
	if cached, found := sharedCache.Get(cacheKey); found {
		if team, ok := cached.(*models.ManagerTeam); ok {
			return team, nil
		}
	}

	currentGameWeekID, err := ms.bootstrapService.GetCurrentGameWeek()
	if err != nil {
		return nil, fmt.Errorf("failed to get current game week: %w", err)
	}

	endpoint := fmt.Sprintf(managerGameWeekPicksEndpoint, managerID, currentGameWeekID)
	resp, err := ms.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get manager team: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get manager team: status %d", resp.StatusCode)
	}

	var team models.ManagerTeam
	if err := json.NewDecoder(resp.Body).Decode(&team); err != nil {
		return nil, fmt.Errorf("failed to decode manager team: %w", err)
	}

	sharedCache.Set(cacheKey, &team, managerCacheTTL)
	return &team, nil
}

func (ms *ManagerService) GetManagerHistory(id int) (*models.ManagerHistory, error) {
	cacheKey := fmt.Sprintf("manager_history_%d", id)
	if cached, found := sharedCache.Get(cacheKey); found {
		if ManagerHistory, ok := cached.(*models.ManagerHistory); ok {
			return ManagerHistory, nil
		}
	}

	endpoint := fmt.Sprintf(managerHistoryEndpoint, id)
	resp, err := ms.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get manager history data: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, fmt.Errorf("manager with ID %d not found", id)
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var managerHistory models.ManagerHistory
	if err := json.Unmarshal(body, &managerHistory); err != nil {
		return nil, fmt.Errorf("failed to decode manager data: %w", err)
	}

	sharedCache.Set(cacheKey, &managerHistory, managerCacheTTL)

	return &managerHistory, nil
}
