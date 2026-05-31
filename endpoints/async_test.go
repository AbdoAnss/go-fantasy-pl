package endpoints_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllPlayersAsync(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := <-c.Players.GetAllPlayersAsync(ctx)
	require.NoError(t, result.Err)
	assert.Len(t, result.Value, 3)
}

func TestGetPlayerHistoryAsync(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := <-c.Players.GetPlayerHistoryAsync(ctx, 101)
	require.NoError(t, result.Err)
	require.NotNil(t, result.Value)
	assert.Len(t, result.Value.History, 2)
}

func TestGetPlayerHistoryAsyncContextCancel(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	resultCh := c.Players.GetPlayerHistoryAsync(ctx, 101)
	cancel()

	select {
	case _, ok := <-resultCh:
		t.Logf("channel closed after context cancellation (ok=%v)", ok)
	case <-time.After(5 * time.Second):
		t.Fatal("expected channel to close after context cancellation")
	}
}

func TestGetPlayerHistoriesBatch(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ids := []int{101, 102, 103}
	resultCh := c.Players.GetPlayerHistoriesBatch(ctx, ids)

	results := make(map[int]int, len(ids))
	for result := range resultCh {
		require.NoError(t, result.Err)
		require.NotNil(t, result.History)
		results[result.PlayerID] = len(result.History.History)
	}

	assert.Len(t, results, len(ids))
	assert.Equal(t, 2, results[101])
	assert.Equal(t, 2, results[102])
	assert.Equal(t, 1, results[103])
}

func TestGetAllFixturesAsync(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := <-c.Fixtures.GetAllFixturesAsync(ctx)
	require.NoError(t, result.Err)
	assert.Len(t, result.Value, 2)
}

func TestGetAllTeamsAsync(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := <-c.Teams.GetAllTeamsAsync(ctx)
	require.NoError(t, result.Err)
	assert.Len(t, result.Value, 3)
}

func TestConcurrentAsyncRequests(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	playersCh := c.Players.GetAllPlayersAsync(ctx)
	fixturesCh := c.Fixtures.GetAllFixturesAsync(ctx)
	teamsCh := c.Teams.GetAllTeamsAsync(ctx)

	playersResult := <-playersCh
	fixturesResult := <-fixturesCh
	teamsResult := <-teamsCh

	require.NoError(t, playersResult.Err)
	require.NoError(t, fixturesResult.Err)
	require.NoError(t, teamsResult.Err)

	assert.Len(t, playersResult.Value, 3)
	assert.Len(t, fixturesResult.Value, 2)
	assert.Len(t, teamsResult.Value, 3)
}
