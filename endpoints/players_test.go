package endpoints_test

import (
	"encoding/json"
	"testing"

	"github.com/AbdoAnss/go-fantasy-pl/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fixtureElements mirrors the parts of the bootstrap capture the player
// tests need, so expectations are derived from the capture itself instead of
// hardcoded values.
func fixtureElements(t *testing.T) []struct {
	ID      int     `json:"id"`
	NowCost float64 `json:"now_cost"`
} {
	t.Helper()

	var boot struct {
		Elements []struct {
			ID      int     `json:"id"`
			NowCost float64 `json:"now_cost"`
		} `json:"elements"`
	}
	require.NoError(t, json.Unmarshal(readTestdata(t, "bootstrap-static.json"), &boot))
	require.NotEmpty(t, boot.Elements)
	return boot.Elements
}

func TestGetAllPlayers(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	players, err := c.Players.GetAllPlayers()
	require.NoError(t, err)

	want := fixtureElements(t)
	require.Len(t, players, len(want))

	ids := make(map[int]bool, len(players))
	for _, p := range players {
		assert.Positive(t, p.ID)
		assert.NotEmpty(t, p.WebName)
		ids[p.ID] = true
	}
	assert.Len(t, ids, len(players), "player IDs should be unique")
}

func TestGetPlayer(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	sample := fixtureElements(t)[0]

	player, err := c.Players.GetPlayer(sample.ID)
	require.NoError(t, err)
	require.NotNil(t, player)

	// Echo and self-consistency invariants: no real-world values.
	assert.Equal(t, sample.ID, player.ID, "returned player should echo the requested ID")
	assert.Equal(t, sample.NowCost, player.NowCost)
	assert.Equal(t, sample.NowCost/10, player.GetPriceInPounds())
	assert.NotEmpty(t, player.WebName)
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

	var want struct {
		Fixtures    []any `json:"fixtures"`
		History     []any `json:"history"`
		HistoryPast []any `json:"history_past"`
	}
	require.NoError(t, json.Unmarshal(readTestdata(t, "element-summary-101.json"), &want))
	require.Len(t, history.History, len(want.History))
	require.Len(t, history.HistoryPast, len(want.HistoryPast))
	require.Len(t, history.Fixtures, len(want.Fixtures))

	for _, gw := range history.History {
		assert.Positive(t, gw.Round)
		assert.Positive(t, gw.Element)
	}
	for _, past := range history.HistoryPast {
		assert.NotEmpty(t, past.SeasonName)
	}
}

func TestPlayerHelpers(t *testing.T) {
	// Pure unit tests on a synthetic player: our formatting, not API data.
	player := &models.Player{FirstName: "Test", SecondName: "Player", NowCost: 55}

	assert.Equal(t, "Test Player", player.GetDisplayName())
	assert.Equal(t, 5.5, player.GetPriceInPounds())
}
