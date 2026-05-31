package endpoints_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllTeams(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	teams, err := c.Teams.GetAllTeams()
	require.NoError(t, err)
	require.Len(t, teams, 3)

	assert.Equal(t, "Arsenal", teams[0].GetFullName())
	assert.Equal(t, "ARS", teams[0].GetShortName())
	assert.Equal(t, 5, teams[0].Strength)
}

func TestGetTeam(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	team, err := c.Teams.GetTeam(2)
	require.NoError(t, err)
	require.NotNil(t, team)

	assert.Equal(t, "Manchester City", team.GetFullName())
	assert.Equal(t, 4, team.Strength)
	assert.Equal(t, 75.0, team.GetWinRate())
}

func TestGetTeam_NotFound(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	team, err := c.Teams.GetTeam(9999)
	require.Error(t, err)
	assert.Nil(t, team)
	assert.Equal(t, "team with ID 9999 not found", err.Error())
}
