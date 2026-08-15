package endpoints_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AbdoAnss/go-fantasy-pl/client"
	"github.com/AbdoAnss/go-fantasy-pl/endpoints"
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
		assert.ErrorContains(t, err, "league with ID 1221170 not found")
		assert.ErrorIs(t, err, endpoints.ErrLeagueNotFound)
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
		assert.Contains(t, err.Error(), "failed to decode H2H matches data")
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

// newH2HSemanticsServer routes requests to fixtures that exercise the live
// API's event/page semantics. The mixed-feed and event=1 fixtures start from
// a real captured response (h2h-league-matches.json); the knockout fixture
// is synthetic because no live knockout capture was available.
func newH2HSemanticsServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		event := r.URL.Query().Get("event")
		page := r.URL.Query().Get("page")

		var fixture string
		switch {
		case event == "" && page == "1":
			fixture = "h2h-matches-mixed.json"
		case event == "1" && page == "1":
			fixture = "h2h-matches-event1.json"
		case event == "1" && page == "2":
			fixture = "h2h-matches-event1-page2.json"
		case event == "37":
			fixture = "h2h-matches-knockout.json"
		case event == "999":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"detail": "Invalid event."}`))
			return
		default:
			http.NotFound(w, r)
			return
		}
		writeTestdata(t, w, fixture)
	}))
}

// Omitting event returns a paginated mixed feed spanning multiple gameweeks.
func TestGetH2HLeagueMatches_MixedFeedWithoutEvent(t *testing.T) {
	server := newH2HSemanticsServer(t)
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithMemoryCache())
	require.NoError(t, err)

	feed, err := c.Leagues.GetH2HLeagueMatches(1221170, 1, 0)
	require.NoError(t, err)

	assert.True(t, feed.HasNext, "mixed feed should have more pages")
	assert.Equal(t, 1, feed.Page)
	require.Len(t, feed.Results, 2)

	events := map[int]bool{}
	for _, m := range feed.Results {
		events[m.Event] = true
	}
	assert.Greater(t, len(events), 1, "mixed feed should span multiple gameweeks, got %v", events)

	// First result matches the real captured match.
	first := feed.Results[0]
	assert.Equal(t, 21353851, first.ID)
	assert.Equal(t, "uzzifc", first.Entry1Name)
	assert.Equal(t, "Saige Fc", first.Entry2Name)
	assert.False(t, first.IsKnockout)
	assert.Nil(t, first.Winner)
}

// event=<gw> filters matches to a single gameweek and can change has_next.
func TestGetH2HLeagueMatches_EventFilterSingleGameweek(t *testing.T) {
	server := newH2HSemanticsServer(t)
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithMemoryCache())
	require.NoError(t, err)

	feed, err := c.Leagues.GetH2HLeagueMatches(1221170, 1, 1)
	require.NoError(t, err)

	assert.False(t, feed.HasNext, "filtered gameweek should have no next page")
	require.Len(t, feed.Results, 1)
	assert.Equal(t, 1, feed.Results[0].Event)
}

// An event-filtered page past the last result is empty with has_next=false.
func TestGetH2HLeagueMatches_EmptyFilteredPage(t *testing.T) {
	server := newH2HSemanticsServer(t)
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithMemoryCache())
	require.NoError(t, err)

	feed, err := c.Leagues.GetH2HLeagueMatches(1221170, 2, 1)
	require.NoError(t, err)

	assert.Empty(t, feed.Results)
	assert.False(t, feed.HasNext)
	assert.Equal(t, 2, feed.Page)
}

// Knockout rounds return is_knockout=true, a populated winner, and a
// knockout_name. The fixture is synthetic (see newH2HSemanticsServer).
func TestGetH2HLeagueMatches_KnockoutRound(t *testing.T) {
	server := newH2HSemanticsServer(t)
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithMemoryCache())
	require.NoError(t, err)

	feed, err := c.Leagues.GetH2HLeagueMatches(1221170, 1, 37)
	require.NoError(t, err)

	require.Len(t, feed.Results, 1)
	m := feed.Results[0]
	assert.True(t, m.IsKnockout)
	require.NotNil(t, m.Winner)
	assert.Contains(t, []int{m.Entry1Entry, m.Entry2Entry}, *m.Winner,
		"knockout winner should be one of the two entries")
	assert.Equal(t, "Round 1", m.KnockoutName)
	assert.Equal(t, 37, m.Event)
}

// An event outside the valid gameweek range yields HTTP 400, surfaced as
// *ErrInvalidH2HQuery instead of a generic status-code error.
func TestGetH2HLeagueMatches_InvalidEventReturnsDomainError(t *testing.T) {
	server := newH2HSemanticsServer(t)
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithMemoryCache())
	require.NoError(t, err)

	feed, err := c.Leagues.GetH2HLeagueMatches(1221170, 1, 999)
	assert.Nil(t, feed)
	require.Error(t, err)

	var invalidQuery *endpoints.ErrInvalidH2HQuery
	require.True(t, errors.As(err, &invalidQuery),
		"expected *ErrInvalidH2HQuery for event=999, got: %v", err)
	assert.Equal(t, 1221170, invalidQuery.LeagueID)
	assert.Equal(t, 999, invalidQuery.Event)
	assert.Equal(t, "Invalid event.", invalidQuery.Detail)
}

// Live API tests for the H2H matches endpoint. The FPL API purges leagues
// between seasons, so a 404 is treated as a skip rather than a failure.
func TestGetH2HLeagueMatches_Live(t *testing.T) {
	skipUnlessLive(t)

	testClient, err := client.NewClient(client.WithMemoryCache())
	require.NoError(t, err)

	leagueID := 1221170

	t.Run("MixedFeed", func(t *testing.T) {
		feed, err := testClient.Leagues.GetH2HLeagueMatches(leagueID, 1, 0)
		if errors.Is(err, endpoints.ErrLeagueNotFound) {
			t.Skip("league no longer exists on the live API")
		}
		require.NoError(t, err)
		assert.NotEmpty(t, feed.Results)
	})

	t.Run("EventFilter", func(t *testing.T) {
		feed, err := testClient.Leagues.GetH2HLeagueMatches(leagueID, 1, 1)
		if errors.Is(err, endpoints.ErrLeagueNotFound) {
			t.Skip("league no longer exists on the live API")
		}
		require.NoError(t, err)
		for _, m := range feed.Results {
			assert.Equal(t, 1, m.Event)
		}
	})

	t.Run("InvalidEvent", func(t *testing.T) {
		_, err := testClient.Leagues.GetH2HLeagueMatches(leagueID, 1, 999)
		if errors.Is(err, endpoints.ErrLeagueNotFound) {
			t.Skip("league no longer exists on the live API")
		}
		var invalidQuery *endpoints.ErrInvalidH2HQuery
		require.True(t, errors.As(err, &invalidQuery),
			"expected *ErrInvalidH2HQuery for event=999, got: %v", err)
	})
}
