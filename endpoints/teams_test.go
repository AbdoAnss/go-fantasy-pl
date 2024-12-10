package endpoints_test

import (
	"testing"

	"github.com/AbdoAnss/go-fantasy-pl/client"
	"github.com/stretchr/testify/assert"
)

var teamID = 13 // Example team ID

func init() {
	testClient = client.NewClient()
}

func TestGetAllTeams(t *testing.T) {
	teams, err := testClient.Teams.GetAllTeams()
	assert.NoError(t, err, "expected no error when getting all teams")
	assert.NotEmpty(t, teams, "expected teams to be returned from API")

	// Note: Stats like points, wins, draws, losses are typically 0 at season start
	// or during pre-season. Use strength ratings for consistent team comparisons.
	for i, team := range teams {
		t.Logf("Team %d: %s, Short name: %s, Code: %d, Strength: %d",
			i+1,
			team.GetFullName(),
			team.GetShortName(),
			team.Code,
			team.Strength)

		assert.NotEmpty(t, team.Name, "expected team name to be non-empty")
		assert.GreaterOrEqual(t, team.Strength, 0, "expected team strength to be non-negative")

		if i >= 3 {
			break
		}
	}
}

func TestGetTeam(t *testing.T) {
	team, err := testClient.Teams.GetTeam(teamID)
	assert.NoError(t, err, "expected no error when getting team")
	assert.NotNil(t, team, "expected team to be returned, got nil")

	t.Logf("----------------------------------------")
	t.Logf("Team: %s", team.GetShortName())
	t.Logf("Team ID: %d", team.ID)
	t.Logf("Strength Ratings:")
	t.Logf("  Overall: %d", team.Strength)
	t.Logf("  Home Attack: %d", team.StrengthAttackHome)
	t.Logf("  Home Defense: %d", team.StrengthDefenceHome)
	t.Logf("  Away Attack: %d", team.StrengthAttackAway)
	t.Logf("  Away Defense: %d", team.StrengthDefenceAway)

	// Performance stats may be 0 depending on season timing
	t.Logf("\nSeason Stats (may be 0 during pre-season or start of season):")
	t.Logf("  Points: %d", team.Points)
	t.Logf("  Position: %d", team.Position)
	t.Logf("  W/D/L: %d/%d/%d", team.Win, team.Draw, team.Loss)
}
