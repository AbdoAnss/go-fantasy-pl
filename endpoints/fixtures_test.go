package endpoints_test

import (
	"testing"

	"github.com/AbdoAnss/go-fantasy-pl/client"
	"github.com/AbdoAnss/go-fantasy-pl/endpoints"
	"github.com/AbdoAnss/go-fantasy-pl/models"
	"github.com/stretchr/testify/assert"
)

var fixtureID int

func setupFixtureTestService() *endpoints.FixtureService {
	c := client.NewClient()
	fixtureID = 8
	return endpoints.NewFixtureService(c)
}

// Helper function to get team map
func getTeamMap(t *testing.T, c *client.Client) map[int]*models.Team {
	teams, err := c.Bootstrap.GetTeams()
	if err != nil {
		t.Fatalf("Failed to get teams: %v", err)
	}

	teamMap := make(map[int]*models.Team)
	for _, team := range teams {
		// Create a new variable for each team to ensure we're not capturing the loop variable
		team := team
		teamMap[team.ID] = &team
	}
	return teamMap
}

func TestGetAllFixtures(t *testing.T) {
	c := client.NewClient()
	fs := endpoints.NewFixtureService(c)

	// Get team mappings for better readability
	teamMap := getTeamMap(t, c)

	fixtures, err := fs.GetAllFixtures()
	assert.NoError(t, err, "expected no error when getting all fixtures")
	assert.NotEmpty(t, fixtures, "expected fixtures to be returned from API")

	t.Logf("Retrieved %d fixtures from the API.", len(fixtures))

	for i, fixture := range fixtures {
		homeTeam := "Unknown"
		awayTeam := "Unknown"

		if home, ok := teamMap[fixture.TeamH]; ok {
			homeTeam = home.GetShortName()
		}
		if away, ok := teamMap[fixture.TeamA]; ok {
			awayTeam = away.GetShortName()
		}

		t.Logf("Fixture %d: ID: %d, %s vs %s",
			i+1,
			fixture.ID,
			homeTeam,
			awayTeam)

		if i >= 3 {
			break
		}
	}
}

func TestGetFixture(t *testing.T) {
	fs := setupFixtureTestService()
	fixture, err := fs.GetFixture(fixtureID)
	assert.NoError(t, err, "expected no error when getting fixture")
	assert.NotNil(t, fixture, "expected fixture to be returned, got nil")
	// Log fixture details
	t.Logf("Fixture ID: %d", fixture.ID)
	t.Logf("Team A: %d vs Team H: %d", fixture.TeamA, fixture.TeamH)
	t.Logf("Fixture Finished: %v", fixture.Finished)
	t.Logf("Kickoff Time: %v", fixture.KickoffTime)
}

func TestGetNonExistentFixture(t *testing.T) {
	fs := setupFixtureTestService()
	fixture, err := fs.GetFixture(999)
	assert.Error(t, err, "expected an error when getting a non-existent fixture")
	assert.Nil(t, fixture, "expected fixture to be nil for non-existent fixture")
	assert.Equal(t, "fixture with ID 999 not found", err.Error())
	t.Logf("Error encountered: %s", err.Error())
}

func TestGetGoalscorers(t *testing.T) {
	fs := setupFixtureTestService()
	fixture, err := fs.GetFixture(fixtureID)
	assert.NoError(t, err, "expected no error when getting fixture")
	assert.NotNil(t, fixture, "expected fixture to be returned, got nil")
	goalscorers, err := fixture.GetGoalscorers()
	assert.NoError(t, err, "expected no error when getting goalscorers")
	t.Logf("Goalscorers: %+v", goalscorers)
}

func TestGetAssisters(t *testing.T) {
	fs := setupFixtureTestService()
	fixture, err := fs.GetFixture(fixtureID)
	assert.NoError(t, err, "expected no error when getting fixture")
	assert.NotNil(t, fixture, "expected fixture to be returned, got nil")
	assisters, err := fixture.GetAssisters()
	assert.NoError(t, err, "expected no error when getting assisters")
	t.Logf("Assisters: %+v", assisters)
}

func TestGetBonus(t *testing.T) {
	fs := setupFixtureTestService()
	fixture, err := fs.GetFixture(fixtureID)
	assert.NoError(t, err, "expected no error when getting fixture")
	assert.NotNil(t, fixture, "expected fixture to be returned, got nil")
	bonus, err := fixture.GetBonus()
	assert.NoError(t, err, "expected no error when getting bonus points")
	t.Logf("Bonus Points: %+v", bonus)
}
