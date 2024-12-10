package endpoints_test

import (
	"testing"

	"github.com/AbdoAnss/go-fantasy-pl/client"
	"github.com/AbdoAnss/go-fantasy-pl/endpoints"
	"github.com/stretchr/testify/assert"
)

var playerID int

func setupPlayersTestService() *endpoints.PlayerService {
	c := client.NewClient()
	playerID = 328 // Example player ID for Mohamed Salah

	return endpoints.NewPlayerService(c)
}

func TestGetAllPlayers(t *testing.T) {
	ps := setupPlayersTestService()
	players, err := ps.GetAllPlayers()

	assert.NoError(t, err, "expected no error when getting all players")

	assert.NotEmpty(t, players, "expected players to be returned from API")

	for i, player := range players {
		t.Logf("Player %d: %s, Team: %d, Points: %d",
			i+1,
			player.GetDisplayName(),
			player.Team,
			player.TotalPoints)

		if i >= 3 {
			break
		}
	}
}

func TestGetPlayer(t *testing.T) {
	ps := setupPlayersTestService()

	player, err := ps.GetPlayer(playerID)

	// Assert no error occurred
	assert.NoError(t, err, "expected no error when getting player")

	// Assert player is returned
	assert.NotNil(t, player, "expected player to be returned, got nil")

	// Log player details
	t.Logf("Player: %s", player.GetDisplayName())
	t.Logf("Price in pound £: %.2f", player.GetPriceInPounds())
	t.Logf("Total Points: %d", player.TotalPoints)
	t.Logf("Form: %v", player.Form)
}

func TestGetPlayerHistory(t *testing.T) {
	ps := setupPlayersTestService()

	history, err := ps.GetPlayerHistory(playerID)

	// Assert no error occurred
	assert.NoError(t, err, "expected no error when getting player history")

	// Assert history is returned
	assert.NotNil(t, history, "expected player history to be returned, got nil")

	// Log history details
	t.Logf("Current season appearances: %d", len(history.History))
	t.Logf("Years in Premier League: %d", len(history.HistoryPast))
}
