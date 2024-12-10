package endpoints_test

import (
	"testing"

	"github.com/AbdoAnss/go-fantasy-pl/client"
	"github.com/AbdoAnss/go-fantasy-pl/endpoints"
	"github.com/stretchr/testify/assert"
)

var teamID int

func setupTeamsTestService() *endpoints.TeamService {
	c := client.NewClient()
	teamID = 13

	return endpoints.NewTeamService(c)
}

func TestGetAllTeams(t *testing.T) {
	ts := setupTeamsTestService()
	teams, err := ts.GetAllTeams()

	assert.NoError(t, err, "expected no error when getting all teams")

	assert.NotEmpty(t, teams, "expected teams to be returned from API")

	for i, team := range teams {
		t.Logf("Team %d: %s, Short name: %s, Code: %d, Points: %d",
			i+1,
			team.GetFullName(),
			team.GetShortName(),
			team.Code,
			team.Points) // somehow points are always 0

		assert.NotEmpty(t, team.Name, "expected team name to be non-empty")
		assert.GreaterOrEqual(t, team.Points, 0, "expected team points to be non-negative")

		if i >= 3 {
			break
		}
	}
}

func TestGetTeam(t *testing.T) {
	ts := setupTeamsTestService()

	team, err := ts.GetTeam(teamID)

	assert.NoError(t, err, "expected no error when getting team")

	assert.NotNil(t, team, "expected team to be returned, got nil")

	t.Logf("----------------------------------------")
	t.Logf("Team: %s", team.GetShortName())
	t.Logf("Team ID: %d", team.ID)
	t.Logf("Points: %d", team.Points)     // always 0 ?
	t.Logf("Wins: %d", team.Win)          // always 0 ?
	t.Logf("Draws: %d", team.Draw)        // always 0 ?
	t.Logf("Losses: %d", team.Loss)       // always 0 ?
	t.Logf("Position: %d", team.Position) // always 0 ?
	t.Logf("Strength: %d", team.Strength)
}
