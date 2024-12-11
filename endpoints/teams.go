package endpoints

import (
	"fmt"

	"github.com/AbdoAnss/go-fantasy-pl/api"
	"github.com/AbdoAnss/go-fantasy-pl/models"
)

// TeamService provides access to team-related data from the FPL API.
// Teams have attributes including ID, name, strength ratings (overall/attack/defense)
// for both home and away matches. Note that some stats (points, played, wins, etc.)
// may be zero, especially early in the season.

type TeamService struct {
	client           api.Client
	bootstrapService *BootstrapService
}

func NewTeamService(client api.Client, bootstrap *BootstrapService) *TeamService {
	return &TeamService{
		client:           client,
		bootstrapService: bootstrap,
	}
}

func (ts *TeamService) GetAllTeams() ([]models.Team, error) {
	return ts.bootstrapService.GetTeams()
}

func (ts *TeamService) GetTeam(id int) (*models.Team, error) {
	teams, err := ts.GetAllTeams()
	if err != nil {
		return nil, err
	}

	for _, t := range teams {
		if t.ID == id {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("team with ID %d not found", id)
}
