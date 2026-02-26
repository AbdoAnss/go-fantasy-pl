package endpoints_test

import (
	"context"
	"testing"
	"time"

	"github.com/AbdoAnss/go-fantasy-pl/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllPlayersAsync(t *testing.T) {
	c := client.NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resultCh := c.Players.GetAllPlayersAsync(ctx)
	result := <-resultCh

	require.NoError(t, result.Err, "expected no error when getting all players asynchronously")
	assert.NotEmpty(t, result.Value, "expected players to be returned")
	t.Logf("Retrieved %d players asynchronously", len(result.Value))
}

func TestGetPlayerHistoryAsync(t *testing.T) {
	c := client.NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resultCh := c.Players.GetPlayerHistoryAsync(ctx, playerID)
	result := <-resultCh

	require.NoError(t, result.Err, "expected no error when getting player history asynchronously")
	assert.NotNil(t, result.Value, "expected player history to be returned")
	t.Logf("Retrieved history for player %d: %d gameweeks", playerID, len(result.Value.History))
}

func TestGetPlayerHistoryAsyncContextCancel(t *testing.T) {
	c := client.NewClient()
	ctx, cancel := context.WithCancel(context.Background())

	resultCh := c.Players.GetPlayerHistoryAsync(ctx, playerID)

	// Cancel context immediately
	cancel()

	// The channel should either return a result or be closed without sending
	select {
	case _, ok := <-resultCh:
		// Either got a result before cancellation or channel was closed; both are valid
		t.Logf("Channel closed (ok=%v) after context cancellation", ok)
	case <-time.After(5 * time.Second):
		t.Error("expected channel to be closed after context cancellation")
	}
}

func TestGetPlayerHistoriesBatch(t *testing.T) {
	c := client.NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Fetch histories for a few players concurrently
	ids := []int{328, 233, 427} // Salah, Haaland, Palmer (example IDs)
	resultCh := c.Players.GetPlayerHistoriesBatch(ctx, ids)

	received := 0
	for result := range resultCh {
		if result.Err != nil {
			t.Logf("Error fetching history for player %d: %v", result.PlayerID, result.Err)
		} else {
			t.Logf("Retrieved history for player %d: %d gameweeks", result.PlayerID, len(result.History.History))
		}
		received++
	}

	assert.Equal(t, len(ids), received, "expected results for all requested player IDs")
}

func TestGetAllFixturesAsync(t *testing.T) {
	c := client.NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resultCh := c.Fixtures.GetAllFixturesAsync(ctx)
	result := <-resultCh

	require.NoError(t, result.Err, "expected no error when getting all fixtures asynchronously")
	assert.NotEmpty(t, result.Value, "expected fixtures to be returned")
	t.Logf("Retrieved %d fixtures asynchronously", len(result.Value))
}

func TestGetAllTeamsAsync(t *testing.T) {
	c := client.NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resultCh := c.Teams.GetAllTeamsAsync(ctx)
	result := <-resultCh

	require.NoError(t, result.Err, "expected no error when getting all teams asynchronously")
	assert.NotEmpty(t, result.Value, "expected teams to be returned")
	t.Logf("Retrieved %d teams asynchronously", len(result.Value))
}

func TestConcurrentAsyncRequests(t *testing.T) {
	c := client.NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Launch players, fixtures, and teams fetches concurrently
	playersCh := c.Players.GetAllPlayersAsync(ctx)
	fixturesCh := c.Fixtures.GetAllFixturesAsync(ctx)
	teamsCh := c.Teams.GetAllTeamsAsync(ctx)

	playersResult := <-playersCh
	fixturesResult := <-fixturesCh
	teamsResult := <-teamsCh

	require.NoError(t, playersResult.Err, "expected no error for concurrent players fetch")
	require.NoError(t, fixturesResult.Err, "expected no error for concurrent fixtures fetch")
	require.NoError(t, teamsResult.Err, "expected no error for concurrent teams fetch")

	assert.NotEmpty(t, playersResult.Value, "expected players from concurrent fetch")
	assert.NotEmpty(t, fixturesResult.Value, "expected fixtures from concurrent fetch")
	assert.NotEmpty(t, teamsResult.Value, "expected teams from concurrent fetch")

	t.Logf("Concurrent fetch: %d players, %d fixtures, %d teams",
		len(playersResult.Value), len(fixturesResult.Value), len(teamsResult.Value))
}
