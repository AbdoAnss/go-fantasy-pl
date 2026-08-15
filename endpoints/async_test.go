package endpoints_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fixtureCounts derives the expected element/history lengths from the
// committed captures, so these tests survive fixture refreshes.
func fixtureCounts(t *testing.T) (players, teams, fixtures int) {
	t.Helper()

	var boot struct {
		Elements []any `json:"elements"`
		Teams    []any `json:"teams"`
	}
	require.NoError(t, json.Unmarshal(readTestdata(t, "bootstrap-static.json"), &boot))

	var fx []any
	require.NoError(t, json.Unmarshal(readTestdata(t, "fixtures.json"), &fx))

	return len(boot.Elements), len(boot.Teams), len(fx)
}

func fixtureHistoryLen(t *testing.T, playerID int) int {
	t.Helper()

	var payload struct {
		History []any `json:"history"`
	}
	require.NoError(t, json.Unmarshal(readTestdata(t, jsonName("element-summary-%d.json", playerID)), &payload))
	return len(payload.History)
}

func TestGetAllPlayersAsync(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := <-c.Players.GetAllPlayersAsync(ctx)
	require.NoError(t, result.Err)

	wantPlayers, _, _ := fixtureCounts(t)
	assert.Len(t, result.Value, wantPlayers)
}

func TestGetPlayerHistoryAsync(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := <-c.Players.GetPlayerHistoryAsync(ctx, 101)
	require.NoError(t, result.Err)
	require.NotNil(t, result.Value)
	assert.Len(t, result.Value.History, fixtureHistoryLen(t, 101))
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

	require.Len(t, results, len(ids), "every requested player should report exactly once")
	for _, id := range ids {
		assert.Equal(t, fixtureHistoryLen(t, id), results[id],
			"player %d history length should match the capture", id)
	}
}

func TestGetAllFixturesAsync(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := <-c.Fixtures.GetAllFixturesAsync(ctx)
	require.NoError(t, result.Err)

	_, _, wantFixtures := fixtureCounts(t)
	assert.Len(t, result.Value, wantFixtures)
}

func TestGetAllTeamsAsync(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := <-c.Teams.GetAllTeamsAsync(ctx)
	require.NoError(t, result.Err)

	_, wantTeams, _ := fixtureCounts(t)
	assert.Len(t, result.Value, wantTeams)
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

	wantPlayers, wantTeams, wantFixtures := fixtureCounts(t)
	assert.Len(t, playersResult.Value, wantPlayers)
	assert.Len(t, fixturesResult.Value, wantFixtures)
	assert.Len(t, teamsResult.Value, wantTeams)
}
