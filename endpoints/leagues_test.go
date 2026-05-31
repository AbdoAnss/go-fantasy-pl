package endpoints_test

import (
	"testing"

	"github.com/AbdoAnss/go-fantasy-pl/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLeagueEndpoints(t *testing.T) {
	skipUnlessLive(t)

	testLeagueClient, err := client.NewClient(client.WithMemoryCache())
	require.NoError(t, err)

	t.Run("GetClassicLeague", func(t *testing.T) {
		leagueID := 1185652
		league, err := testLeagueClient.Leagues.GetClassicLeagueStandings(leagueID, 1)
		assert.NoError(t, err)
		assert.NotNil(t, league)
	})

	t.Run("ValidateLeagueData", func(t *testing.T) {
		leagueID := 1185652
		league, err := testLeagueClient.Leagues.GetClassicLeagueStandings(leagueID, 1)
		require.NoError(t, err)

		assert.Equal(t, leagueID, league.League.ID)
		assert.NotEmpty(t, league.League.Name)
		assert.NotEmpty(t, league.GetTopManagers(1))
	})

	t.Run("GetNonExistentLeague", func(t *testing.T) {
		league, err := testLeagueClient.Leagues.GetClassicLeagueStandings(99999999, 1)
		assert.Error(t, err)
		assert.Nil(t, league)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("CacheConsistency", func(t *testing.T) {
		leagueID := 1185652

		league1, err := testLeagueClient.Leagues.GetClassicLeagueStandings(leagueID, 1)
		require.NoError(t, err)

		league2, err := testLeagueClient.Leagues.GetClassicLeagueStandings(leagueID, 1)
		require.NoError(t, err)

		assert.Equal(t, league1.GetLeagueInfo(), league2.GetLeagueInfo())
	})
}
