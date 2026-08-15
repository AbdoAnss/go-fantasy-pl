package models

import "fmt"

// H2HMatchesFeed is the paginated response from the H2H league matches
// endpoint (/leagues-h2h-matches/league/<id>/).
//
// Semantics:
//   - Without an event filter, results contain a mixed feed of matches
//     spanning multiple gameweeks, ordered by event then match id.
//   - With event=<gw>, results are restricted to that gameweek (or knockout
//     round when gw is in the knockout range).
//   - has_next indicates whether another page exists for the same query.
type H2HMatchesFeed struct {
	HasNext bool       `json:"has_next"`
	Page    int        `json:"page"`
	Results []H2HMatch `json:"results"`
}

// H2HMatch is a single head-to-head fixture between two entries.
type H2HMatch struct {
	ID               string `json:"id"`
	Event            int    `json:"event"`
	Entry1Match      int    `json:"entry_1_match"`
	Entry1Points     int    `json:"entry_1_points"`
	Entry1Win        bool   `json:"entry_1_win"`
	Entry1Draw       bool   `json:"entry_1_draw"`
	Entry1Loss       bool   `json:"entry_1_loss"`
	Entry1Name       string `json:"entry_1_name"`
	Entry1PlayerName string `json:"entry_1_player_name"`
	Entry2Match      int    `json:"entry_2_match"`
	Entry2Points     int    `json:"entry_2_points"`
	Entry2Win        bool   `json:"entry_2_win"`
	Entry2Draw       bool   `json:"entry_2_draw"`
	Entry2Loss       bool   `json:"entry_2_loss"`
	Entry2Name       string `json:"entry_2_name"`
	Entry2PlayerName string `json:"entry_2_player_name"`

	// Knockout fields: is_knockout is true for knockout-round matches, in
	// which case winner holds the winning entry's match id (0/null for a
	// draw before penalties resolution) and knockout_name names the round
	// (e.g. "Round 1", "Quarter Final").
	IsKnockout   bool   `json:"is_knockout"`
	Winner       int    `json:"winner"`
	KnockoutName string `json:"knockout_name"`
}

// IsDraw reports whether the match ended in a draw. For knockout matches a
// draw means the winner was decided by penalties (winner is 0 in that case
// only when unresolved).
func (m *H2HMatch) IsDraw() bool {
	return m.Entry1Draw && m.Entry2Draw
}

// WinnerName returns the entry name of the winning side, or "" for a draw.
func (m *H2HMatch) WinnerName() string {
	switch {
	case m.Entry1Win:
		return m.Entry1Name
	case m.Entry2Win:
		return m.Entry2Name
	default:
		return ""
	}
}

// Score returns the match score as "entry1Points - entry2Points".
func (m *H2HMatch) Score() string {
	return fmt.Sprintf("%d - %d", m.Entry1Points, m.Entry2Points)
}
