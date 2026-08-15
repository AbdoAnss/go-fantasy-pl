package endpoints_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/AbdoAnss/go-fantasy-pl/client"
	"github.com/AbdoAnss/go-fantasy-pl/endpoints"
	"github.com/AbdoAnss/go-fantasy-pl/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newH2HTestServer serves canned H2H matches responses matching the live
// FPL API semantics for league 1221170.
func newH2HTestServer(t *testing.T) (*client.Client, *httptest.Server, *[]url.Values) {
	t.Helper()

	var queries []url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		queries = append(queries, r.URL.Query())

		event := r.URL.Query().Get("event")
		page := r.URL.Query().Get("page")

		if !strings.Contains(r.URL.Path, "/league/1221170/") {
			http.NotFound(w, r)
			return
		}

		var fixture string
		switch {
		case event == "" && page == "":
			fixture = "h2h-matches-mixed.json"
		case event == "" && page == "1":
			fixture = "h2h-matches-page1.json"
		case event == "1" && page == "":
			fixture = "h2h-matches-event1.json"
		case event == "1" && page == "2":
			fixture = "h2h-matches-event1-page2.json"
		case event == "36":
			fixture = "h2h-matches-event36.json"
		case event == "37":
			fixture = "h2h-matches-event37-knockout.json"
		case event == "999":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"detail": "Invalid event."}`))
			return
		case event == "-1" || event == "abc":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"detail": "Invalid event."}`))
			return
		default:
			http.NotFound(w, r)
			return
		}
		writeTestdata(t, w, fixture)
	}))

	c, err := client.NewClient(
		client.WithBaseURL(server.URL),
		client.WithMemoryCache(),
	)
	require.NoError(t, err)

	return c, server, &queries
}

func TestGetH2HLeagueMatches_MixedFeedWithoutEvent(t *testing.T) {
	c, server, queries := newH2HTestServer(t)
	defer server.Close()

	feed, err := c.Leagues.GetH2HLeagueMatches(1221170)
	require.NoError(t, err)

	// Omitting event returns a paginated mixed feed spanning gameweeks.
	assert.True(t, feed.HasNext, "mixed feed should have more pages")
	assert.Equal(t, 1, feed.Page)
	assert.NotEmpty(t, feed.Results)

	events := map[int]bool{}
	for _, m := range feed.Results {
		events[m.Event] = true
	}
	assert.Greater(t, len(events), 1, "mixed feed should span multiple gameweeks, got %v", events)
	assert.False(t, feed.Results[0].IsKnockout)

	sent := (*queries)[0]
	assert.Empty(t, sent.Get("event"), "no event param should be sent")
	assert.Empty(t, sent.Get("page"))
}

func TestGetH2HLeagueMatches_ExplicitPageWithoutEvent(t *testing.T) {
	c, server, queries := newH2HTestServer(t)
	defer server.Close()

	feed, err := c.Leagues.GetH2HLeagueMatches(1221170, endpoints.WithH2HPage(1))
	require.NoError(t, err)

	assert.Equal(t, 1, feed.Page)
	assert.True(t, feed.HasNext)
	assert.Len(t, feed.Results, 4)

	sent := (*queries)[0]
	assert.Equal(t, "1", sent.Get("page"))
	assert.Empty(t, sent.Get("event"))
}

func TestGetH2HLeagueMatches_EventFilterSingleGameweek(t *testing.T) {
	c, server, queries := newH2HTestServer(t)
	defer server.Close()

	feed, err := c.Leagues.GetH2HLeagueMatches(1221170, endpoints.WithH2HEvent(1))
	require.NoError(t, err)

	// event=<gw> restricts matches to one gameweek and can flip has_next.
	assert.False(t, feed.HasNext, "filtered gameweek should have no next page")
	for _, m := range feed.Results {
		assert.Equal(t, 1, m.Event)
	}

	assert.Equal(t, "1", (*queries)[0].Get("event"))
}

func TestGetH2HLeagueMatches_EmptyFilteredPage(t *testing.T) {
	c, server, _ := newH2HTestServer(t)
	defer server.Close()

	feed, err := c.Leagues.GetH2HLeagueMatches(1221170, endpoints.WithH2HEvent(1), endpoints.WithH2HPage(2))
	require.NoError(t, err)

	assert.Empty(t, feed.Results)
	assert.False(t, feed.HasNext)
	assert.Equal(t, 2, feed.Page)
}

func TestGetH2HLeagueMatches_LaterRegularGameweek(t *testing.T) {
	c, server, _ := newH2HTestServer(t)
	defer server.Close()

	feed, err := c.Leagues.GetH2HLeagueMatches(1221170, endpoints.WithH2HEvent(36))
	require.NoError(t, err)

	require.Len(t, feed.Results, 2)
	for _, m := range feed.Results {
		assert.Equal(t, 36, m.Event)
		assert.False(t, m.IsKnockout)
		assert.Empty(t, m.KnockoutName)
		assert.Zero(t, m.Winner)
	}

	m := feed.Results[0]
	assert.Equal(t, "Pep's Chess Club", m.WinnerName())
	assert.Equal(t, "55 - 62", m.Score())
	assert.False(t, m.IsDraw())
}

func TestGetH2HLeagueMatches_KnockoutRound(t *testing.T) {
	c, server, _ := newH2HTestServer(t)
	defer server.Close()

	feed, err := c.Leagues.GetH2HLeagueMatches(1221170, endpoints.WithH2HEvent(37))
	require.NoError(t, err)

	require.NotEmpty(t, feed.Results)
	for _, m := range feed.Results {
		assert.True(t, m.IsKnockout, "event 37 should be a knockout round")
		assert.NotZero(t, m.Winner, "knockout winner should be populated")
		assert.Equal(t, "Round 1", m.KnockoutName)
		assert.Equal(t, 37, m.Event)
	}

	// winner holds the winning entry's match id.
	m := feed.Results[0]
	assert.Equal(t, m.Entry2Match, m.Winner)
	assert.Equal(t, "Pep's Chess Club", m.WinnerName())
}

func TestGetH2HLeagueMatches_InvalidEventReturnsDomainError(t *testing.T) {
	c, server, _ := newH2HTestServer(t)
	defer server.Close()

	for _, event := range []int{999, -1} {
		feed, err := c.Leagues.GetH2HLeagueMatches(1221170, endpoints.WithH2HEvent(event))
		assert.Nil(t, feed)
		require.Error(t, err)

		var invalidQuery *endpoints.ErrInvalidH2HQuery
		require.True(t, errors.As(err, &invalidQuery),
			"expected *ErrInvalidH2HQuery for event=%d, got: %v", event, err)
		assert.Equal(t, 1221170, invalidQuery.LeagueID)
		assert.Contains(t, err.Error(), "400")
		if event == 999 {
			assert.Contains(t, err.Error(), "event=999")
		}
	}
}

func TestGetH2HLeagueMatches_LeagueNotFound(t *testing.T) {
	c, server, _ := newH2HTestServer(t)
	defer server.Close()

	feed, err := c.Leagues.GetH2HLeagueMatches(99999999)
	assert.Nil(t, feed)
	require.Error(t, err)
	assert.ErrorIs(t, err, endpoints.ErrLeagueNotFound)
}

func TestH2HMatchHelpers(t *testing.T) {
	draw := models.H2HMatch{
		Entry1Points: 40, Entry1Draw: true,
		Entry2Points: 40, Entry2Draw: true,
		Entry1Name: "A", Entry2Name: "B",
	}
	assert.True(t, draw.IsDraw())
	assert.Empty(t, draw.WinnerName())
	assert.Equal(t, "40 - 40", draw.Score())
}

// Live API tests for the H2H matches endpoint. The FPL API purges leagues
// between seasons, so a 404 is treated as a skip rather than a failure.
func TestGetH2HLeagueMatches_Live(t *testing.T) {
	skipUnlessLive(t)

	testClient, err := client.NewClient(client.WithMemoryCache())
	require.NoError(t, err)

	leagueID := 1221170

	t.Run("MixedFeed", func(t *testing.T) {
		feed, err := testClient.Leagues.GetH2HLeagueMatches(leagueID)
		if errors.Is(err, endpoints.ErrLeagueNotFound) {
			t.Skip("league no longer exists on the live API")
		}
		require.NoError(t, err)
		assert.NotEmpty(t, feed.Results)
	})

	t.Run("EventFilter", func(t *testing.T) {
		feed, err := testClient.Leagues.GetH2HLeagueMatches(leagueID, endpoints.WithH2HEvent(1))
		if errors.Is(err, endpoints.ErrLeagueNotFound) {
			t.Skip("league no longer exists on the live API")
		}
		require.NoError(t, err)
		for _, m := range feed.Results {
			assert.Equal(t, 1, m.Event)
		}
	})

	t.Run("InvalidEvent", func(t *testing.T) {
		_, err := testClient.Leagues.GetH2HLeagueMatches(leagueID, endpoints.WithH2HEvent(999))
		if errors.Is(err, endpoints.ErrLeagueNotFound) {
			t.Skip("league no longer exists on the live API")
		}
		var invalidQuery *endpoints.ErrInvalidH2HQuery
		require.True(t, errors.As(err, &invalidQuery),
			"expected *ErrInvalidH2HQuery for event=999, got: %v", err)
	})
}
