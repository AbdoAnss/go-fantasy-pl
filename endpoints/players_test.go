package endpoints_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllPlayers(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	players, err := c.Players.GetAllPlayers()
	require.NoError(t, err)
	require.Len(t, players, 3)

	assert.Equal(t, 101, players[0].ID)
	assert.Equal(t, "Bukayo Saka", players[0].GetDisplayName())
	assert.Equal(t, 105.0, players[0].NowCost)
}

func TestGetPlayer(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	player, err := c.Players.GetPlayer(102)
	require.NoError(t, err)
	require.NotNil(t, player)

	assert.Equal(t, "Erling Haaland", player.GetDisplayName())
	assert.Equal(t, 26.4, player.GetPriceInPounds())
	assert.Equal(t, 220, player.TotalPoints)
}

func TestGetPlayer_NotFound(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	player, err := c.Players.GetPlayer(9999)
	require.Error(t, err)
	assert.Nil(t, player)
	assert.Equal(t, "player with ID 9999 not found", err.Error())
}

func TestGetPlayerHistory(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	history, err := c.Players.GetPlayerHistory(101)
	require.NoError(t, err)
	require.NotNil(t, history)

	require.Len(t, history.History, 2)
	require.Len(t, history.HistoryPast, 1)
	assert.Equal(t, 1, history.History[0].Round)
	assert.Equal(t, 9, history.History[0].TotalPoints)
	assert.Equal(t, "2024/25", history.HistoryPast[0].SeasonName)
}
