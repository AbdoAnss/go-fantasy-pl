package endpoints_test

import (
	"testing"

	"github.com/AbdoAnss/go-fantasy-pl/client"
	"github.com/stretchr/testify/assert"
)

var testLeagueClient *client.Client

func init() {
	testLeagueClient = client.NewClient()
}

func TestLeagueEndpoints(t *testing.T) {
	t.Run("GetClassicLeague", func(t *testing.T) {
		// INPT Fantasy LeagueID
		leagueID := 1185652
		page := 1

		league, err := testLeagueClient.Leagues.GetClassicLeagueStandings(leagueID, page)
		assert.NoError(t, err, "expected no error when getting classic league")
		assert.NotNil(t, league, "expected league to be returned")

		// Log league details
		t.Logf("\nLeague Details:")
		t.Logf("League: %s", league.GetLeagueInfo())
		t.Logf("Created: %s", league.League.GetCreationDate())
		t.Logf("Type: %s", league.League.LeagueType)
		t.Logf("Last Updated: %s", league.GetUpdateTime())
		t.Logf("Max Entries: %d", league.League.GetMaxEntries())

		// Log standings
		t.Logf("\nTop 4 Managers:")
		for _, manager := range league.GetTopManagers(4) {
			t.Logf("%s", manager.GetManagerInfo())
			t.Logf("  Points: %d (GW: %d)", manager.Total, manager.EventTotal)
			t.Logf("  Rank: %d %s", manager.Rank, manager.GetRankChangeString())
		}

		// Log pagination
		t.Logf("\nPagination:")
		t.Logf("Current: %s", league.Standings.GetPageInfo())
		t.Logf("Has Previous: %v", league.Standings.HasPreviousPage())
		t.Logf("Has Next: %v", league.Standings.HasNext)
	})

	t.Run("ValidateLeagueData", func(t *testing.T) {
		leagueID := 1185652
		league, err := testLeagueClient.Leagues.GetClassicLeagueStandings(leagueID, 1)
		assert.NoError(t, err)

		// Validate league structure
		assert.Equal(t, leagueID, league.League.ID)
		assert.NotEmpty(t, league.League.Name)
		assert.NotZero(t, league.League.Created)

		// Validate standings
		topManagers := league.GetTopManagers(1)
		assert.NotEmpty(t, topManagers)
		assert.Greater(t, topManagers[0].Entry, 0)
		assert.NotEmpty(t, topManagers[0].GetManagerInfo())
		assert.GreaterOrEqual(t, topManagers[0].Total, 0)
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
		assert.NoError(t, err)

		league2, err := testLeagueClient.Leagues.GetClassicLeagueStandings(leagueID, 1)
		assert.NoError(t, err)

		assert.Equal(t, league1.GetLeagueInfo(), league2.GetLeagueInfo())
		assert.Equal(t,
			league1.GetTopManagers(1)[0].GetManagerInfo(),
			league2.GetTopManagers(1)[0].GetManagerInfo())
	})
}
