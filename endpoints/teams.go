package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AbdoAnss/go-fantasy-pl/api"
	"github.com/AbdoAnss/go-fantasy-pl/internal/cache"
	"github.com/AbdoAnss/go-fantasy-pl/models"
)

/*
**
* Team endpoint is useful for matching team IDs with their corresponding details.
* This endpoint provides information about various teams in the league, including attributes such as:
* - code: A unique identifier for the team.
* - id: The team's ID used for referencing.
* - name: The full name of the team.
* - short_name: A shortened version of the team's name.
* - points: The total points accumulated by the team.
* - played: The number of matches played by the team.
* - win, draw, loss: The counts of wins, draws, and losses respectively.
* - strength: A general strength rating of the team.
* - strength_overall_home and strength_overall_away: Strength ratings for home and away matches.
* - strength_attack_home and strength_attack_away: Attack strength ratings for home and away matches.
* - strength_defence_home and strength_defence_away: Defense strength ratings for home and away matches.
* - pulse_id: A unique identifier used in the FPL system for the team.
*
* Upon inspecting the JSON response, it is observed that some attributes (such as points, played, win, draw, and loss)
* are always 0, especially at the beginning of the season or during certain periods. This makes these attributes
* less interesting for analysis, as they do not provide meaningful insights during those times.
* However, strength-related attributes can still offer valuable insights into the team's potential performance.
*
* This endpoint is essential for applications that require team-specific data for analysis,
* fantasy league management, or displaying team information to users.
 */

const (
	teamsEndpoint = "/bootstrap-static/"
)

type TeamService struct {
	client api.Client
}

func NewTeamService(client api.Client) *TeamService {
	return &TeamService{
		client: client,
	}
}

// TODO:
// Centralized Cache with Namespacing:
// Use a single cache instance and differentiate keys using endpoint-specific prefixes.

var (
	teamCacheTTL = 5 * time.Hour // team infos are rarely modified
	teamsCache   = cache.NewCache()
)

func init() {
	teamsCache.StartCleanupTask(5 * time.Minute)
}

func (ts *TeamService) GetAllTeams() ([]models.Team, error) {
	if cached, found := teamsCache.Get("teams"); found {
		if teams, ok := cached.([]models.Team); ok {
			return teams, nil
		}
	}

	resp, err := ts.client.Get(teamsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get teams: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		Elements []models.Team `json:"teams"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode teams: %w", err)
	}

	teamsCache.Set("teams", response.Elements, teamCacheTTL)

	return response.Elements, nil
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

	return nil, fmt.Errorf("player with ID %d not found", id)
}
