package models

type ManagerTeam struct {
	ActiveChip    *string        `json:"active_chip"`    // Active chip (e.g., "wildcard", "benchboost")
	AutomaticSubs []AutomaticSub `json:"automatic_subs"` // List of automatic substitutions
	EntryHistory  EntryHistory   `json:"entry_history"`  // Entry history details
	Picks         []Pick         `json:"picks"`          // Picks in the team
}

type AutomaticSub struct {
	Entry      int `json:"entry"`       // Entry ID of the player
	ElementIn  int `json:"element_in"`  // Player ID coming in
	ElementOut int `json:"element_out"` // Player ID going out
	Event      int `json:"event"`       // Gameweek ID
}

type EntryHistory struct {
	Event              int `json:"event"`                // Gameweek ID
	Points             int `json:"points"`               // Points scored
	TotalPoints        int `json:"total_points"`         // Total points
	Rank               int `json:"rank"`                 // Current rank
	RankSort           int `json:"rank_sort"`            // Rank sort
	OverallRank        int `json:"overall_rank"`         // Overall rank
	PercentileRank     int `json:"percentile_rank"`      // Percentile rank
	Bank               int `json:"bank"`                 // Money in bank
	Value              int `json:"value"`                // Team value
	EventTransfers     int `json:"event_transfers"`      // Transfers made this week
	EventTransfersCost int `json:"event_transfers_cost"` // Cost of transfers
	PointsOnBench      int `json:"points_on_bench"`      // Points scored by bench players
}

type Pick struct {
	Element       int  `json:"element"`         // Player ID
	Position      int  `json:"position"`        // Position in team (1-15)
	Multiplier    int  `json:"multiplier"`      // 2 for captain, 3 for triple captain, 0 for benched
	IsCaptain     bool `json:"is_captain"`      // Is this player captain?
	IsViceCaptain bool `json:"is_vice_captain"` // Is this player vice-captain?
	ElementType   int  `json:"element_type"`    // Type of the player (e.g., defender, midfielder)
}

func (mt *ManagerTeam) GetStartingXI() []Pick {
	starters := make([]Pick, 0, 11)
	for _, pick := range mt.Picks {
		if pick.Position <= 11 {
			starters = append(starters, pick)
		}
	}
	return starters
}

func (mt *ManagerTeam) GetBench() []Pick {
	bench := make([]Pick, 0, 4)
	for _, pick := range mt.Picks {
		if pick.Position > 11 {
			bench = append(bench, pick)
		}
	}
	return bench
}

func (mt *ManagerTeam) GetTeamValueInMillions() float64 {
	return float64(mt.EntryHistory.Value) / 10
}

func (mt *ManagerTeam) GetBankValueInMillions() float64 {
	return float64(mt.EntryHistory.Bank) / 10
}
