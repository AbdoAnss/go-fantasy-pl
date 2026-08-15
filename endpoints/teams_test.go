package endpoints_test

import (
	"encoding/json"
	"testing"

	"github.com/AbdoAnss/go-fantasy-pl/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fixtureTeamIDs(t *testing.T) []int {
	t.Helper()

	var boot struct {
		Teams []struct {
			ID int `json:"id"`
		} `json:"teams"`
	}
	require.NoError(t, json.Unmarshal(readTestdata(t, "bootstrap-static.json"), &boot))
	require.NotEmpty(t, boot.Teams)
	ids := make([]int, len(boot.Teams))
	for i, tm := range boot.Teams {
		ids[i] = tm.ID
	}
	return ids
}

func TestGetAllTeams(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	teams, err := c.Teams.GetAllTeams()
	require.NoError(t, err)

	want := fixtureTeamIDs(t)
	require.Len(t, teams, len(want))

	ids := make(map[int]bool, len(teams))
	for _, tm := range teams {
		assert.Positive(t, tm.ID)
		assert.NotEmpty(t, tm.Name)
		assert.NotEmpty(t, tm.ShortName)
		// Strength is null before the season starts and 1-5 once set.
		if tm.Strength != 0 {
			assert.GreaterOrEqual(t, tm.Strength, 1)
			assert.LessOrEqual(t, tm.Strength, 5)
		}
		ids[tm.ID] = true
	}
	assert.Len(t, ids, len(teams), "team IDs should be unique")
}

func TestGetTeam(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	id := fixtureTeamIDs(t)[1]

	team, err := c.Teams.GetTeam(id)
	require.NoError(t, err)
	require.NotNil(t, team)

	assert.Equal(t, id, team.ID, "returned team should echo the requested ID")
	assert.NotEmpty(t, team.GetFullName())
	assert.NotEmpty(t, team.GetShortName())

	// Self-consistency: the rate helpers must agree with the raw counters.
	if team.Played > 0 {
		assert.InDelta(t, float64(team.Win)/float64(team.Played)*100, team.GetWinRate(), 1e-9)
		assert.InDelta(t, float64(team.Draw)/float64(team.Played)*100, team.GetDrawRate(), 1e-9)
		assert.InDelta(t, float64(team.Loss)/float64(team.Played)*100, team.GetLossRate(), 1e-9)
	} else {
		assert.Zero(t, team.GetWinRate())
	}
}

func TestGetTeam_NotFound(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	ids := fixtureTeamIDs(t)
	maxID := ids[0]
	for _, id := range ids[1:] {
		if id > maxID {
			maxID = id
		}
	}

	team, err := c.Teams.GetTeam(maxID + 100)
	require.Error(t, err)
	assert.Nil(t, team)
}

func TestTeamRateHelpers(t *testing.T) {
	// Pure unit test on a synthetic team: our arithmetic, not API data.
	team := &models.Team{Win: 3, Draw: 1, Loss: 0, Played: 4}

	assert.Equal(t, 75.0, team.GetWinRate())
	assert.Equal(t, 25.0, team.GetDrawRate())
	assert.Equal(t, 0.0, team.GetLossRate())

	team.Played = 0
	assert.Equal(t, 0.0, team.GetWinRate())
}
