package endpoints_test

import (
	"encoding/json"
	"testing"

	"github.com/AbdoAnss/go-fantasy-pl/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fixtureFixtureIDs(t *testing.T) []int {
	t.Helper()

	var fixtures []struct {
		ID int `json:"id"`
	}
	require.NoError(t, json.Unmarshal(readTestdata(t, "fixtures.json"), &fixtures))
	require.NotEmpty(t, fixtures)
	ids := make([]int, len(fixtures))
	for i, f := range fixtures {
		ids[i] = f.ID
	}
	return ids
}

func TestGetAllFixtures(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	fixtures, err := c.Fixtures.GetAllFixtures()
	require.NoError(t, err)

	want := fixtureFixtureIDs(t)
	require.Len(t, fixtures, len(want))

	for _, f := range fixtures {
		assert.Positive(t, f.ID)
		assert.NotNil(t, f.KickoffTime, "scheduled fixtures must have a kickoff time")
		assert.GreaterOrEqual(t, f.TeamH, 1)
		assert.LessOrEqual(t, f.TeamH, 20)
		assert.GreaterOrEqual(t, f.TeamA, 1)
		assert.LessOrEqual(t, f.TeamA, 20)
		assert.NotEqual(t, f.TeamH, f.TeamA, "a team cannot play itself")
	}
}

func TestGetFixture(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	id := fixtureFixtureIDs(t)[0]

	fixture, err := c.Fixtures.GetFixture(id)
	require.NoError(t, err)
	require.NotNil(t, fixture)

	assert.Equal(t, id, fixture.ID, "returned fixture should echo the requested ID")
	assert.NotNil(t, fixture.KickoffTime)

	// Finished matches carry scores; scheduled ones do not. Fixtures that
	// started but have not finished may already report live scores.
	if fixture.Finished {
		require.NotNil(t, fixture.TeamHScore)
		require.NotNil(t, fixture.TeamAScore)
		assert.GreaterOrEqual(t, fixture.GetTeamHScore(), 0)
		assert.GreaterOrEqual(t, fixture.GetTeamAScore(), 0)
	}
	if !fixture.Started {
		assert.Nil(t, fixture.TeamHScore)
		assert.Nil(t, fixture.TeamAScore)
	}
}

func TestGetNonExistentFixture(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	ids := fixtureFixtureIDs(t)
	maxID := ids[0]
	for _, id := range ids[1:] {
		if id > maxID {
			maxID = id
		}
	}

	fixture, err := c.Fixtures.GetFixture(maxID + 1000)
	require.Error(t, err)
	assert.Nil(t, fixture)
}

// statFixture builds a synthetic match with one contributed stat per side,
// exercising our stat-extraction helpers independently of API data.
func statFixture() *models.Fixture {
	return &models.Fixture{
		ID:    1,
		TeamH: 1,
		TeamA: 2,
		Stats: []models.Stat{
			{
				Identifier: models.StatGoalsScored,
				H:          []models.StatDetail{{Value: 2, Element: 101}},
				A:          []models.StatDetail{{Value: 1, Element: 201}},
			},
			{
				Identifier: models.StatAssists,
				H:          []models.StatDetail{{Value: 1, Element: 103}},
				A:          []models.StatDetail{},
			},
			{
				Identifier: models.StatBonus,
				H:          []models.StatDetail{{Value: 3, Element: 101}},
				A:          []models.StatDetail{},
			},
		},
	}
}

func TestGetGoalscorers(t *testing.T) {
	goalscorers, err := statFixture().GetGoalscorers()
	require.NoError(t, err)

	assert.Len(t, goalscorers["h"], 1)
	assert.Len(t, goalscorers["a"], 1)
	assert.Equal(t, 101, goalscorers["h"][0].Element)
	assert.Equal(t, 2, goalscorers["h"][0].Value)
}

func TestGetAssisters(t *testing.T) {
	assisters, err := statFixture().GetAssisters()
	require.NoError(t, err)

	assert.Len(t, assisters["h"], 1)
	assert.Equal(t, 103, assisters["h"][0].Element)
	assert.Empty(t, assisters["a"])
}

func TestGetBonus(t *testing.T) {
	bonus, err := statFixture().GetBonus()
	require.NoError(t, err)

	assert.Len(t, bonus["h"], 1)
	assert.Equal(t, 3, bonus["h"][0].Value)
	assert.Empty(t, bonus["a"])
}
