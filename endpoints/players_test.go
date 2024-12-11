package endpoints_test

import (
	"testing"

	"github.com/AbdoAnss/go-fantasy-pl/client"
	"github.com/stretchr/testify/assert"
)

var (
	testClient *client.Client
	playerID   = 328 // Example player ID for Mohamed Salah
)

func init() {
	testClient = client.NewClient()
}

func TestGetAllPlayers(t *testing.T) {
	players, err := testClient.Players.GetAllPlayers()
	assert.NoError(t, err, "expected no error when getting all players")
	assert.NotEmpty(t, players, "expected players to be returned from API")

	for i, player := range players {
		t.Logf("Player %d: %s, Team: %d, Points: %d",
			i+1, player.GetDisplayName(), player.Team, player.TotalPoints)
		if i >= 3 {
			break
		}
	}
}

func TestGetPlayer(t *testing.T) {
	player, err := testClient.Players.GetPlayer(playerID)
	assert.NoError(t, err, "expected no error when getting player")
	assert.NotNil(t, player, "expected player to be returned, got nil")

	t.Logf("Player: %s", player.GetDisplayName())
	t.Logf("Price in pound £: %.2f", player.GetPriceInPounds())
	t.Logf("Total Points: %d", player.TotalPoints)
	t.Logf("Form: %v", player.Form)
}

func TestGetPlayerHistory(t *testing.T) {
	history, err := testClient.Players.GetPlayerHistory(playerID)
	assert.NoError(t, err, "expected no error when getting player history")
	assert.NotNil(t, history, "expected player history to be returned, got nil")

	t.Logf("Current season appearances: %d", len(history.History))
	t.Logf("Years in Premier League: %d", len(history.HistoryPast))
}
