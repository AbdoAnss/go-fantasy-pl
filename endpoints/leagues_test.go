package endpoints_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestGetH2HLeagueStandings(t *testing.T) {
	testLeagueClient, server := newLeagueTestClient(t)
	defer server.Close()

	league, err := testLeagueClient.Leagues.GetH2HLeagueStandings(42)
	require.NoError(t, err)
	require.NotNil(t, league)

	assert.Equal(t, 42, league.League.ID)
	assert.Equal(t, "Test H2H League", league.League.Name)
	require.Len(t, league.Standings.Results, 1)
	assert.Equal(t, "North Stand FC", league.Standings.Results[0].EntryName)
	assert.Equal(t, 7, league.Standings.Results[0].PointsTotal)
	require.Len(t, league.MatchesNext.Results, 1)
	assert.Equal(t, 9001, league.MatchesNext.Results[0].Entry1Entry)
	require.NotNil(t, league.MatchesNext.Results[0].Winner)
	assert.Equal(t, 9001, *league.MatchesNext.Results[0].Winner)
}

func TestGetH2HLeagueStandingsRejectsInvalidLeagueID(t *testing.T) {
	testLeagueClient, server := newLeagueTestClient(t)
	defer server.Close()

	league, err := testLeagueClient.Leagues.GetH2HLeagueStandings(0)
	require.Error(t, err)
	assert.Nil(t, league)
	assert.Contains(t, err.Error(), "league ID must be positive")
}

func TestGetH2HLeagueStandingsNotFound(t *testing.T) {
	testLeagueClient, server := newLeagueTestClient(t)
	defer server.Close()

	league, err := testLeagueClient.Leagues.GetH2HLeagueStandings(404)
	require.Error(t, err)
	assert.Nil(t, league)
	assert.Contains(t, err.Error(), "not found")
}

func TestGetH2HLeagueStandingsMalformedJSON(t *testing.T) {
	testLeagueClient, server := newLeagueTestClient(t)
	defer server.Close()

	league, err := testLeagueClient.Leagues.GetH2HLeagueStandings(500)
	require.Error(t, err)
	assert.Nil(t, league)
	assert.Contains(t, err.Error(), "failed to decode h2h league data")
}

func TestGetH2HLeagueStandingsUnexpectedStatus(t *testing.T) {
	testLeagueClient, server := newLeagueTestClient(t)
	defer server.Close()

	league, err := testLeagueClient.Leagues.GetH2HLeagueStandings(503)
	require.Error(t, err)
	assert.Nil(t, league)
	assert.Contains(t, err.Error(), "unexpected status code: 503")
}

func newLeagueTestClient(t *testing.T) (*client.Client, *httptest.Server) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/leagues-h2h-matches/league/42/":
			writeTestdata(t, w, "h2h-league.json")
		case "/leagues-h2h-matches/league/500/":
			_, err := fmt.Fprint(w, `{"league":`)
			require.NoError(t, err)
		case "/leagues-h2h-matches/league/503/":
			w.WriteHeader(http.StatusServiceUnavailable)
			_, err := fmt.Fprint(w, `{"detail":"service unavailable"}`)
			require.NoError(t, err)
		default:
			http.NotFound(w, r)
		}
	}))

	c, err := client.NewClient(
		client.WithBaseURL(server.URL),
		client.WithMemoryCache(),
	)
	require.NoError(t, err)

	return c, server
}
