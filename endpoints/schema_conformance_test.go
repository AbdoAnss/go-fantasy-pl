package endpoints_test

import (
	"encoding/json"
	"testing"

	"github.com/AbdoAnss/go-fantasy-pl/endpoints"
	"github.com/AbdoAnss/go-fantasy-pl/internal/conformance"
	"github.com/AbdoAnss/go-fantasy-pl/models"
	"github.com/stretchr/testify/require"
)

// Hermetic schema conformance: every committed capture in testdata is decoded
// into the models and checked with conformance.Check. These tests are the
// first line of defense against API schema drift — they run in CI, need no
// network, and never assert on real-world values, so refreshing the captures
// via `make recapture` is a no-op unless the schema actually changed.

func TestBootstrapStaticSchema(t *testing.T) {
	raw := readTestdata(t, "bootstrap-static.json")

	var resp endpoints.Response
	require.NoError(t, json.Unmarshal(raw, &resp))

	conformance.Check(t, raw, conformance.Spec{Model: &resp, Allowlist: bootstrapAllowlist})

	require.NotEmpty(t, resp.Teams)
	conformance.Check(t,
		conformance.Extract(t, raw, "teams", 0),
		conformance.Spec{Model: &resp.Teams[0], Allowlist: teamAllowlist})

	require.NotEmpty(t, resp.Elements)
	conformance.Check(t,
		conformance.Extract(t, raw, "elements", 0),
		conformance.Spec{Model: &resp.Elements[0], Allowlist: playerAllowlist})

	require.NotEmpty(t, resp.Events)
	conformance.Check(t,
		conformance.Extract(t, raw, "events", 0),
		conformance.Spec{Model: &resp.Events[0], Allowlist: gameWeekAllowlist})

	conformance.Check(t,
		conformance.Extract(t, raw, "game_settings"),
		conformance.Spec{Model: &resp.Settings, Allowlist: gameSettingsAllowlist})
}

func TestFixturesSchema(t *testing.T) {
	raw := readTestdata(t, "fixtures.json")

	var fixtures []models.Fixture
	require.NoError(t, json.Unmarshal(raw, &fixtures))
	require.NotEmpty(t, fixtures)

	conformance.Check(t, raw, conformance.Spec{Model: fixtures})
}

func TestElementSummarySchema(t *testing.T) {
	raw := readTestdata(t, "element-summary-101.json")

	var history models.PlayerHistory
	require.NoError(t, json.Unmarshal(raw, &history))

	conformance.Check(t, raw, conformance.Spec{Model: &history})

	if len(history.History) > 0 {
		conformance.Check(t,
			conformance.Extract(t, raw, "history", 0),
			conformance.Spec{Model: &history.History[0]})
	}

	require.NotEmpty(t, history.HistoryPast)
	conformance.Check(t,
		conformance.Extract(t, raw, "history_past", 0),
		conformance.Spec{Model: &history.HistoryPast[0], Allowlist: pastHistoryStatsAllowlist})

	require.NotEmpty(t, history.Fixtures)
	conformance.Check(t,
		conformance.Extract(t, raw, "fixtures", 0),
		conformance.Spec{Model: &history.Fixtures[0]})
}

func TestClassicLeagueStandingsSchema(t *testing.T) {
	raw := readTestdata(t, "leagues-classic-standings.json")

	var league models.ClassicLeague
	require.NoError(t, json.Unmarshal(raw, &league))

	conformance.Check(t, raw, conformance.Spec{Model: &league})
}

func TestH2HLeagueMatchesSchema(t *testing.T) {
	raw := readTestdata(t, "h2h-league-matches.json")

	var page models.H2HLeagueMatchesPage
	require.NoError(t, json.Unmarshal(raw, &page))

	conformance.Check(t, raw, conformance.Spec{Model: &page})

	require.NotEmpty(t, page.Results)
	conformance.Check(t,
		conformance.Extract(t, raw, "results", 0),
		conformance.Spec{Model: &page.Results[0]})
}

func TestH2HLeagueStandingsSchema(t *testing.T) {
	raw := readTestdata(t, "h2h-league-standings.json")

	var standings models.H2HLeagueStandings
	require.NoError(t, json.Unmarshal(raw, &standings))

	conformance.Check(t, raw, conformance.Spec{Model: &standings})
}
