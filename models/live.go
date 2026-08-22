package models

import (
	"fmt"
	"strconv"
)

// EventLive holds the live points data for every player in a single
// gameweek, as served by /event/{id}/live/.
//
// The payload is gameweek-scoped and manager-agnostic: one request covers
// all players, so a manager's live points are computed by joining these
// element stats with the manager's picks (see ManagerTeam).
//
// Freshness: during live matches the API reflects in-play stats, with bonus
// derived provisionally from BPS rankings. Bonus is finalized only after the
// fixture ends, and the upstream CDN may serve copies up to a few minutes
// old, so treat this as "live-ish" rather than realtime data.
type EventLive struct {
	Elements []LiveElement `json:"elements"`
}

// LiveElement is the live gameweek record for a single player. ID matches
// the player ID used across the API (bootstrap elements, picks, fixtures).
type LiveElement struct {
	ID       int              `json:"id"`
	Stats    LiveElementStats `json:"stats"`
	Explain  []LiveExplain    `json:"explain"`
	Modified bool             `json:"modified"`
}

// LiveElementStats carries the in-play counting stats for a player.
//
// Influence, Creativity, Threat, IctIndex, and the Expected* fields are
// decimal numbers serialized as strings upstream (e.g. "0.34"); they are
// kept as strings per the Player model convention, with Get* accessors
// returning parsed float64 values.
type LiveElementStats struct {
	Minutes         int `json:"minutes"`
	GoalsScored     int `json:"goals_scored"`
	Assists         int `json:"assists"`
	CleanSheets     int `json:"clean_sheets"`
	GoalsConceded   int `json:"goals_conceded"`
	OwnGoals        int `json:"own_goals"`
	PenaltiesSaved  int `json:"penalties_saved"`
	PenaltiesMissed int `json:"penalties_missed"`
	YellowCards     int `json:"yellow_cards"`
	RedCards        int `json:"red_cards"`
	Saves           int `json:"saves"`
	Bonus           int `json:"bonus"`
	Bps             int `json:"bps"`

	Influence  string `json:"influence"`
	Creativity string `json:"creativity"`
	Threat     string `json:"threat"`
	IctIndex   string `json:"ict_index"`

	ClearancesBlocksInterceptions int `json:"clearances_blocks_interceptions"`
	Recoveries                    int `json:"recoveries"`
	Tackles                       int `json:"tackles"`
	DefensiveContribution         int `json:"defensive_contribution"`
	Starts                        int `json:"starts"`

	ExpectedGoals            string `json:"expected_goals"`
	ExpectedAssists          string `json:"expected_assists"`
	ExpectedGoalInvolvements string `json:"expected_goal_involvements"`
	ExpectedGoalsConceded    string `json:"expected_goals_conceded"`

	TotalPoints int  `json:"total_points"`
	InDreamteam bool `json:"in_dreamteam"`
	Played      bool `json:"played"`
}

// LiveExplain breaks down how a player's points were earned in a single
// fixture. Fixture references the fixture ID used by the Fixtures service.
type LiveExplain struct {
	Fixture int               `json:"fixture"`
	Stats   []LiveExplainStat `json:"stats"`
}

// LiveExplainStat is one scoring rule applied in a fixture: the underlying
// stat Value and the Points it produced (PointsModification covers manual
// adjustments, which are rare).
type LiveExplainStat struct {
	Identifier         string  `json:"identifier"`
	Points             float64 `json:"points"`
	Value              float64 `json:"value"`
	PointsModification float64 `json:"points_modification"`
}

// StatsFor returns the live stats for a player ID, and whether the player
// is present in this gameweek's live data.
func (e *EventLive) StatsFor(playerID int) (*LiveElementStats, bool) {
	for i := range e.Elements {
		if e.Elements[i].ID == playerID {
			return &e.Elements[i].Stats, true
		}
	}
	return nil, false
}

// PointsFor returns the live total points for a player ID. The second
// return value is false when the player has no live record this gameweek.
func (e *EventLive) PointsFor(playerID int) (int, bool) {
	stats, ok := e.StatsFor(playerID)
	if !ok {
		return 0, false
	}
	return stats.TotalPoints, true
}

// ExplainFor returns the per-fixture point breakdown for a player ID, and
// whether the player is present in this gameweek's live data.
func (e *EventLive) ExplainFor(playerID int) ([]LiveExplain, bool) {
	for i := range e.Elements {
		if e.Elements[i].ID == playerID {
			return e.Elements[i].Explain, true
		}
	}
	return nil, false
}

func (s LiveElementStats) parseDecimal(field string) float64 {
	v, err := strconv.ParseFloat(field, 64)
	if err != nil {
		return 0
	}
	return v
}

// GetInfluence returns the Influence score as a float64 (0 on parse error).
func (s LiveElementStats) GetInfluence() float64 { return s.parseDecimal(s.Influence) }

// GetCreativity returns the Creativity score as a float64 (0 on parse error).
func (s LiveElementStats) GetCreativity() float64 { return s.parseDecimal(s.Creativity) }

// GetThreat returns the Threat score as a float64 (0 on parse error).
func (s LiveElementStats) GetThreat() float64 { return s.parseDecimal(s.Threat) }

// GetICTIndex returns the ICT index as a float64 (0 on parse error).
func (s LiveElementStats) GetICTIndex() float64 { return s.parseDecimal(s.IctIndex) }

// GetExpectedGoals returns expected goals as a float64 (0 on parse error).
func (s LiveElementStats) GetExpectedGoals() float64 { return s.parseDecimal(s.ExpectedGoals) }

// GetExpectedAssists returns expected assists as a float64 (0 on parse error).
func (s LiveElementStats) GetExpectedAssists() float64 { return s.parseDecimal(s.ExpectedAssists) }

// GetExpectedGoalInvolvements returns xG plus xA as a float64 (0 on parse error).
func (s LiveElementStats) GetExpectedGoalInvolvements() float64 {
	return s.parseDecimal(s.ExpectedGoalInvolvements)
}

// GetExpectedGoalsConceded returns expected goals conceded as a float64 (0 on parse error).
func (s LiveElementStats) GetExpectedGoalsConceded() float64 {
	return s.parseDecimal(s.ExpectedGoalsConceded)
}

// String implements fmt.Stringer with a compact summary of the stats.
func (s LiveElementStats) String() string {
	return fmt.Sprintf("pts=%d (min=%d goals=%d assists=%d bonus=%d bps=%d)",
		s.TotalPoints, s.Minutes, s.GoalsScored, s.Assists, s.Bonus, s.Bps)
}
