package endpoints_test

import (
	"testing"

	"github.com/AbdoAnss/go-fantasy-pl/client"
	"github.com/AbdoAnss/go-fantasy-pl/endpoints"
	"github.com/stretchr/testify/assert"
)

var fixtureID int

func setupFixtureTestService() *endpoints.FixtureService {
	c := client.NewClient()
	fixtureID = 8

	return endpoints.NewFixtureService(c)
}

func TestGetAllFixtures(t *testing.T) {
	fs := setupFixtureTestService()
	fixtures, err := fs.GetAllFixtures()

	assert.NoError(t, err, "expected no error when getting all fixtures")
	assert.NotEmpty(t, fixtures, "expected fixtures to be returned from API")

	t.Logf("Retrieved %d fixtures from the API.", len(fixtures))

	// TODO: tests to improve when adding teams logic
	// convert ID to team ShortName for better readability

	for i, fixture := range fixtures {
		t.Logf("Fixture %d: ID: %d, Team A: %d, Team H: %d",
			i+1,
			fixture.ID,
			fixture.TeamA,
			fixture.TeamH)

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
