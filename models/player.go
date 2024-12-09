package models

type Player struct {
	ID                       int     `json:"id"`
	FirstName                string  `json:"first_name"`
	SecondName               string  `json:"second_name"`
	WebName                  string  `json:"web_name"`
	Team                     int     `json:"team"`
	TeamCode                 int     `json:"team_code"`
	TotalPoints              int     `json:"total_points"`
	NowCost                  float64 `json:"now_cost"`
	SelectedByPercent        string  `json:"selected_by_percent"`
	Form                     string  `json:"form"`
	InDreamteam              bool    `json:"in_dreamteam"`
	Minutes                  int     `json:"minutes"`
	GoalsScored              int     `json:"goals_scored"`
	Assists                  int     `json:"assists"`
	CleanSheets              int     `json:"clean_sheets"`
	YellowCards              int     `json:"yellow_cards"`
	RedCards                 int     `json:"red_cards"`
	Status                   string  `json:"status"`
	ChanceOfPlaying          float64 `json:"chance_of_playing_next_round"`
	Code                     int     `json:"code"`
	CostChangeEvent          int     `json:"cost_change_event"`
	CostChangeEventFall      int     `json:"cost_change_event_fall"`
	CostChangeStart          int     `json:"cost_change_start"`
	CostChangeStartFall      int     `json:"cost_change_start_fall"`
	DreamteamCount           int     `json:"dreamteam_count"`
	ElementType              int     `json:"element_type"`
	EpNext                   string  `json:"ep_next"`
	EpThis                   string  `json:"ep_this"`
	EventPoints              int     `json:"event_points"`
	News                     string  `json:"news"`
	NewsAdded                string  `json:"news_added"`
	PointsPerGame            string  `json:"points_per_game"`
	Special                  bool    `json:"special"`
	SquadNumber              *int    `json:"squad_number"`
	TransfersIn              int     `json:"transfers_in"`
	TransfersInEvent         int     `json:"transfers_in_event"`
	TransfersOut             int     `json:"transfers_out"`
	TransfersOutEvent        int     `json:"transfers_out_event"`
	ValueForm                string  `json:"value_form"`
	ValueSeason              string  `json:"value_season"`
	Influence                string  `json:"influence"`
	Creativity               string  `json:"creativity"`
	Threat                   string  `json:"threat"`
	IctIndex                 string  `json:"ict_index"`
	Starts                   int     `json:"starts"`
	ExpectedGoals            string  `json:"expected_goals"`
	ExpectedAssists          string  `json:"expected_assists"`
	ExpectedGoalInvolvements string  `json:"expected_goal_involvements"`
	ExpectedGoalsConceded    string  `json:"expected_goals_conceded"`
	NowCostRank              int     `json:"now_cost_rank"`
	NowCostRankType          int     `json:"now_cost_rank_type"`
	FormRank                 int     `json:"form_rank"`
	FormRankType             int     `json:"form_rank_type"`
	PointsPerGameRank        int     `json:"points_per_game_rank"`
	PointsPerGameRankType    int     `json:"points_per_game_rank_type"`
	SelectedRank             int     `json:"selected_rank"`
	SelectedRankType         int     `json:"selected_rank_type"`
}

func (p *Player) GetDisplayName() string {
	return p.FirstName + " " + p.SecondName
}

func (p *Player) GetPriceInPounds() float64 {
	return p.NowCost / 10
}
