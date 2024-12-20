package endpoints

import (
	"encoding/json"
	"fmt"

	"github.com/AbdoAnss/go-fantasy-pl/api"
	"github.com/AbdoAnss/go-fantasy-pl/models"
)

const (
	fixturesEndpoint = "/fixtures/"
)

type FixtureService struct {
	client api.Client
}

type FixtureNotFoundError struct {
	ID int
}

func (e *FixtureNotFoundError) Error() string {
	return fmt.Sprintf("fixture with ID %d not found", e.ID)
}

func NewFixtureService(client api.Client) *FixtureService {
	return &FixtureService{
		client: client,
	}
}

func (fs *FixtureService) GetAllFixtures() ([]models.Fixture, error) {
	if cached, found := sharedCache.Get("fixtures"); found {
		if fixtures, ok := cached.([]models.Fixture); ok {
			return fixtures, nil
		}
	}

	resp, err := fs.client.Get(fixturesEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get fixtures: %w", err)
	}
	defer resp.Body.Close()

	var fixtures []models.Fixture
	if err := json.NewDecoder(resp.Body).Decode(&fixtures); err != nil {
		return nil, fmt.Errorf("failed to decode fixtures: %w", err)
	}

	sharedCache.Set("fixtures", fixtures, fixturesCacheTTL)

	return fixtures, nil
}

func (fs *FixtureService) GetFixture(id int) (*models.Fixture, error) {
	if cached, found := sharedCache.Get(fmt.Sprintf("fixture_%d", id)); found {
		if fixture, ok := cached.(*models.Fixture); ok {
			return fixture, nil
		}
	}

	fixtures, err := fs.GetAllFixtures()
	if err != nil {
		return nil, err
	}

	for _, f := range fixtures {
		if f.ID == id {
			sharedCache.Set(fmt.Sprintf("fixture_%d", id), &f, fixturesCacheTTL)
			return &f, nil
		}
	}

	return nil, &FixtureNotFoundError{ID: id}
}
