package models

type PlayerHistory struct {
	Fixtures    []FixtureStats     `json:"fixtures"`     // Add fixtures field
	History     []GameWeekStats    `json:"history"`      // Keep history field
	HistoryPast []PastHistoryStats `json:"history_past"` // Add history_past field
}

type GameWeekStats struct {
	Element                  int    `json:"element"`
	Fixture                  int    `json:"fixture"`
	OpponentTeam             int    `json:"opponent_team"` // Add opponent_team field
	TotalPoints              int    `json:"total_points"`
	WasHome                  bool   `json:"was_home"`     // Add was_home field
	KickoffTime              string `json:"kickoff_time"` // Add kickoff_time field
	TeamHScore               *int   `json:"team_h_score"` // Use pointer to handle null
	TeamAScore               *int   `json:"team_a_score"` // Use pointer to handle null
	Round                    int    `json:"round"`
	Minutes                  int    `json:"minutes"`
	GoalsScored              int    `json:"goals_scored"`
	Assists                  int    `json:"assists"`
	CleanSheets              int    `json:"clean_sheets"`
	GoalsConceded            int    `json:"goals_conceded"`             // Add goals_conceded field
	OwnGoals                 int    `json:"own_goals"`                  // Add own_goals field
	PenaltiesSaved           int    `json:"penalties_saved"`            // Add penalties_saved field
	PenaltiesMissed          int    `json:"penalties_missed"`           // Add penalties_missed field
	YellowCards              int    `json:"yellow_cards"`               // Add yellow_cards field
	RedCards                 int    `json:"red_cards"`                  // Add red_cards field
	Saves                    int    `json:"saves"`                      // Add saves field
	Bonus                    int    `json:"bonus"`                      // Add bonus field
	Bps                      int    `json:"bps"`                        // Add bps field
	Influence                string `json:"influence"`                  // Add influence field
	Creativity               string `json:"creativity"`                 // Add creativity field
	Threat                   string `json:"threat"`                     // Add threat field
	IctIndex                 string `json:"ict_index"`                  // Add ict_index field
	Starts                   int    `json:"starts"`                     // Add starts field
	ExpectedGoals            string `json:"expected_goals"`             // Add expected_goals field
	ExpectedAssists          string `json:"expected_assists"`           // Add expected_assists field
	ExpectedGoalInvolvements string `json:"expected_goal_involvements"` // Add expected_goal_involvements field
	ExpectedGoalsConceded    string `json:"expected_goals_conceded"`    // Add expected_goals_conceded field
}

type FixtureStats struct {
	ID                   int    `json:"id"`
	Code                 int    `json:"code"`
	TeamH                int    `json:"team_h"`
	TeamHScore           *int   `json:"team_h_score"` // Use pointer to handle null
	TeamA                int    `json:"team_a"`
	TeamAScore           *int   `json:"team_a_score"` // Use pointer to handle null
	Event                int    `json:"event"`
	Finished             bool   `json:"finished"`
	Minutes              int    `json:"minutes"`
	ProvisionalStartTime bool   `json:"provisional_start_time"`
	KickoffTime          string `json:"kickoff_time"`
	EventName            string `json:"event_name"`
	IsHome               bool   `json:"is_home"`
	Difficulty           int    `json:"difficulty"`
}

type PastHistoryStats struct {
	SeasonName               string `json:"season_name"`
	ElementCode              int    `json:"element_code"`
	StartCost                int    `json:"start_cost"`
	EndCost                  int    `json:"end_cost"`
	TotalPoints              int    `json:"total_points"`
	Minutes                  int    `json:"minutes"`
	GoalsScored              int    `json:"goals_scored"`
	Assists                  int    `json:"assists"`
	CleanSheets              int    `json:"clean_sheets"`
	GoalsConceded            int    `json:"goals_conceded"`             // Add goals_conceded field
	OwnGoals                 int    `json:"own_goals"`                  // Add own_goals field
	PenaltiesSaved           int    `json:"penalties_saved"`            // Add penalties_saved field
	PenaltiesMissed          int    `json:"penalties_missed"`           // Add penalties_missed field
	YellowCards              int    `json:"yellow_cards"`               // Add yellow_cards field
	RedCards                 int    `json:"red_cards"`                  // Add red_cards field
	Saves                    int    `json:"saves"`                      // Add saves field
	Bonus                    int    `json:"bonus"`                      // Add bonus field
	Bps                      int    `json:"bps"`                        // Add bps field
	Influence                string `json:"influence"`                  // Add influence field
	Creativity               string `json:"creativity"`                 // Add creativity field
	Threat                   string `json:"threat"`                     // Add threat field
	IctIndex                 string `json:"ict_index"`                  // Add ict_index field
	Starts                   int    `json:"starts"`                     // Add starts field
	ExpectedGoals            string `json:"expected_goals"`             // Add expected_goals field
	ExpectedAssists          string `json:"expected_assists"`           // Add expected_assists field
	ExpectedGoalInvolvements string `json:"expected_goal_involvements"` // Add expected_goal_involvements field
	ExpectedGoalsConceded    string `json:"expected_goals_conceded"`    // Add expected_goals_conceded field
}
