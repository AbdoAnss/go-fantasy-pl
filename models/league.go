package models

import (
	"fmt"
	"time"
)

type ClassicLeague struct {
	NewEntries      NewEntries `json:"new_entries"`
	LastUpdatedData time.Time  `json:"last_updated_data"`
	League          League     `json:"league"`
	Standings       Standings  `json:"standings"`
}

// H2HLeague represents head-to-head league standings and match results.
type H2HLeague struct {
	NewEntries  H2HNewEntries `json:"new_entries"`
	League      H2HLeagueInfo `json:"league"`
	Standings   H2HStandings  `json:"standings"`
	MatchesNext H2HMatches    `json:"matches_next"`
}

type NewEntries struct {
	HasNext bool          `json:"has_next"`
	Page    int           `json:"page"`
	Results []interface{} `json:"results"` // Assuming results can be of various types
}

// H2HNewEntries represents newly joined entries for a head-to-head league.
type H2HNewEntries struct {
	HasNext bool          `json:"has_next"`
	Page    int           `json:"page"`
	Number  int           `json:"number"`
	Results []interface{} `json:"results"`
}

// League represents the league details.
type League struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Created     time.Time `json:"created"`
	Closed      bool      `json:"closed"`
	MaxEntries  *int      `json:"max_entries"` // Pointer to allow null
	LeagueType  string    `json:"league_type"`
	Scoring     string    `json:"scoring"`
	AdminEntry  *int      `json:"admin_entry"` // Pointer to allow null
	StartEvent  int       `json:"start_event"`
	CodePrivacy string    `json:"code_privacy"`
	HasCup      bool      `json:"has_cup"`
	CupLeague   int       `json:"cup_league"`
	Rank        *int      `json:"rank"` // Pointer to allow null
}

// H2HLeagueInfo represents the league details returned by the H2H endpoint.
type H2HLeagueInfo struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Created    time.Time `json:"created"`
	Closed     bool      `json:"closed"`
	LeagueType string    `json:"league_type"`
	Scoring    string    `json:"_scoring"`
	AdminEntry *int      `json:"admin_entry"`
	StartEvent int       `json:"start_event"`
	HasStarted bool      `json:"has_started"`
	ShortName  *string   `json:"short_name"`
	Rank       *int      `json:"rank"`
	Size       *int      `json:"size"`
	KoRounds   int       `json:"ko_rounds"`
}

// Standings represents the standings section of the Classic League.
type Standings struct {
	HasNext bool            `json:"has_next"`
	Page    int             `json:"page"`
	Results []LeagueManager `json:"results"`
}

// H2HStandings represents the standings section of a head-to-head league.
type H2HStandings struct {
	HasNext bool             `json:"has_next"`
	Page    int              `json:"page"`
	Number  int              `json:"number"`
	Results []H2HLeagueEntry `json:"results"`
}

// Player represents an individual player's standings information.
type LeagueManager struct {
	ID         int    `json:"id"`
	EventTotal int    `json:"event_total"`
	PlayerName string `json:"player_name"`
	Rank       int    `json:"rank"`
	LastRank   int    `json:"last_rank"`
	RankSort   int    `json:"rank_sort"`
	Total      int    `json:"total"`
	Entry      int    `json:"entry"`
	EntryName  string `json:"entry_name"`
	HasPlayed  bool   `json:"has_played"`
}

// H2HLeagueEntry represents an individual manager's H2H standings row.
type H2HLeagueEntry struct {
	ID            int    `json:"id"`
	EntryName     string `json:"entry_name"`
	PlayerName    string `json:"player_name"`
	Movement      string `json:"movement"`
	OwnEntry      bool   `json:"own_entry"`
	Rank          int    `json:"rank"`
	LastRank      int    `json:"last_rank"`
	RankSort      int    `json:"rank_sort"`
	Total         int    `json:"total"`
	MatchesPlayed int    `json:"matches_played"`
	MatchesWon    int    `json:"matches_won"`
	MatchesDrawn  int    `json:"matches_drawn"`
	MatchesLost   int    `json:"matches_lost"`
	PointsFor     int    `json:"points_for"`
	PointsAgainst int    `json:"points_against"`
	PointsTotal   int    `json:"points_total"`
	Division      int    `json:"division"`
	Entry         int    `json:"entry"`
}

// H2HMatches represents the current or next H2H match page.
type H2HMatches struct {
	HasNext bool              `json:"has_next"`
	Page    int               `json:"page"`
	Number  int               `json:"number"`
	Results []H2HLeagueResult `json:"results"`
}

// H2HLeagueResult represents one head-to-head match result.
type H2HLeagueResult struct {
	ID               int         `json:"id"`
	Entry1Entry      int         `json:"entry_1_entry"`
	Entry1Name       string      `json:"entry_1_name"`
	Entry1PlayerName string      `json:"entry_1_player_name"`
	Entry2Entry      int         `json:"entry_2_entry"`
	Entry2Name       string      `json:"entry_2_name"`
	Entry2PlayerName string      `json:"entry_2_player_name"`
	IsKnockout       bool        `json:"is_knockout"`
	Winner           *int        `json:"winner"`
	Tiebreak         interface{} `json:"tiebreak"`
	OwnEntry         bool        `json:"own_entry"`
	Entry1Points     int         `json:"entry_1_points"`
	Entry1Win        int         `json:"entry_1_win"`
	Entry1Draw       int         `json:"entry_1_draw"`
	Entry1Loss       int         `json:"entry_1_loss"`
	Entry2Points     int         `json:"entry_2_points"`
	Entry2Win        int         `json:"entry_2_win"`
	Entry2Draw       int         `json:"entry_2_draw"`
	Entry2Loss       int         `json:"entry_2_loss"`
	Entry1Total      int         `json:"entry_1_total"`
	Entry2Total      int         `json:"entry_2_total"`
	SeedValue        interface{} `json:"seed_value"`
	Event            int         `json:"event"`
}

func (cl *ClassicLeague) GetLeagueInfo() string {
	return fmt.Sprintf("%s (ID: %d)", cl.League.Name, cl.League.ID)
}

func (cl *ClassicLeague) GetUpdateTime() string {
	return cl.LastUpdatedData.Format(time.RFC822)
}

func (cl *ClassicLeague) GetTopManagers(n int) []LeagueManager {
	if n > len(cl.Standings.Results) {
		n = len(cl.Standings.Results)
	}
	return cl.Standings.Results[:n]
}

// League methods
func (l *League) GetMaxEntries() int {
	if l.MaxEntries == nil {
		return 0
	}
	return *l.MaxEntries
}

func (l *League) GetAdminEntry() int {
	if l.AdminEntry == nil {
		return 0
	}
	return *l.AdminEntry
}

func (l *League) GetRank() int {
	if l.Rank == nil {
		return 0
	}
	return *l.Rank
}

func (l *League) GetCreationDate() string {
	return l.Created.Format("2006-01-02")
}

// LeagueManager methods
func (lm *LeagueManager) GetManagerInfo() string {
	return fmt.Sprintf("%s (%s)", lm.EntryName, lm.PlayerName)
}

func (lm *LeagueManager) GetRankChange() int {
	return lm.LastRank - lm.Rank
}

func (lm *LeagueManager) GetRankChangeString() string {
	change := lm.GetRankChange()
	if change > 0 {
		return fmt.Sprintf("↑%d", change)
	} else if change < 0 {
		return fmt.Sprintf("↓%d", -change)
	}
	return "→"
}

// Standings methods
func (s *Standings) GetPageInfo() string {
	return fmt.Sprintf("Page %d", s.Page)
}

func (s *Standings) HasPreviousPage() bool {
	return s.Page > 1
}
