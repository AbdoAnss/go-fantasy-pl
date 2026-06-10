package endpoints_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AbdoAnss/go-fantasy-pl/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetH2HLeagueMatches(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/leagues-h2h-matches/league/1221170/", r.URL.Path)
			assert.Equal(t, "1", r.URL.Query().Get("page"))
			assert.Equal(t, "3", r.URL.Query().Get("event"))

			w.Header().Set("Content-Type", "application/json")
			writeTestdata(t, w, "h2h-league-matches.json")
		}))
		defer server.Close()

		c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithMemoryCache())
		require.NoError(t, err)

		resp, err := c.Leagues.GetH2HLeagueMatches(1221170, 1, 3)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Results, 1)

		assert.True(t, resp.HasNext)
		assert.Equal(t, 1, resp.Page)
		assert.Equal(t, "uzzifc", resp.Results[0].Entry1Name)
		assert.Equal(t, 1, resp.Results[0].Event)
		assert.Nil(t, resp.Results[0].Winner)
	})

	t.Run("invalid input", func(t *testing.T) {
		c, err := client.NewClient(client.WithMemoryCache())
		require.NoError(t, err)

		resp, err := c.Leagues.GetH2HLeagueMatches(0, 1, 0)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "league ID must be positive")

		resp, err = c.Leagues.GetH2HLeagueMatches(1221170, 0, 0)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "page must be positive")

		resp, err = c.Leagues.GetH2HLeagueMatches(1221170, 1, -1)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "event cannot be negative")
	})

	t.Run("not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}))
		defer server.Close()

		c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithMemoryCache())
		require.NoError(t, err)

		resp, err := c.Leagues.GetH2HLeagueMatches(1221170, 1, 0)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "league with ID 1221170 not found")
	})

	t.Run("malformed json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("{"))
		}))
		defer server.Close()

		c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithMemoryCache())
		require.NoError(t, err)

		resp, err := c.Leagues.GetH2HLeagueMatches(1221170, 1, 0)
		assert.Nil(t, resp)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode H2H league matches data")
	})

	t.Run("unexpected status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		}))
		defer server.Close()

		c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithMemoryCache())
		require.NoError(t, err)

		resp, err := c.Leagues.GetH2HLeagueMatches(1221170, 1, 0)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "unexpected status code: 418")
	})
}

func TestGetH2HLeagueStandings(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/leagues-h2h/1221170/standings/", r.URL.Path)
			assert.Equal(t, "1", r.URL.Query().Get("page_standings"))

			w.Header().Set("Content-Type", "application/json")
			writeTestdata(t, w, "h2h-league-standings.json")
		}))
		defer server.Close()

		c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithMemoryCache())
		require.NoError(t, err)

		resp, err := c.Leagues.GetH2HLeagueStandings(1221170, 1)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Standings.Results, 1)

		assert.Equal(t, 1221170, resp.League.ID)
		assert.Equal(t, 1, resp.Standings.Page)
		assert.Nil(t, resp.LastUpdatedData)
		assert.Equal(t, "Saige Fc", resp.Standings.Results[0].EntryName)
		assert.Equal(t, 72, resp.Standings.Results[0].Total)
	})

	t.Run("invalid input", func(t *testing.T) {
		c, err := client.NewClient(client.WithMemoryCache())
		require.NoError(t, err)

		resp, err := c.Leagues.GetH2HLeagueStandings(0, 1)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "league ID must be positive")

		resp, err = c.Leagues.GetH2HLeagueStandings(1221170, 0)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "page must be positive")
	})

	t.Run("not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}))
		defer server.Close()

		c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithMemoryCache())
		require.NoError(t, err)

		resp, err := c.Leagues.GetH2HLeagueStandings(1221170, 1)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "league with ID 1221170 not found")
	})

	t.Run("malformed json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("{"))
		}))
		defer server.Close()

		c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithMemoryCache())
		require.NoError(t, err)

		resp, err := c.Leagues.GetH2HLeagueStandings(1221170, 1)
		assert.Nil(t, resp)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode H2H league standings data")
	})

	t.Run("unexpected status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithMemoryCache())
		require.NoError(t, err)

		resp, err := c.Leagues.GetH2HLeagueStandings(1221170, 1)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "unexpected status code: 500")
	})
}
