package models

type Team struct {
	Code                int     `json:"code"`
	Draw                int     `json:"draw"`
	Form                *string `json:"form"` // Using pointer to handle null values
	ID                  int     `json:"id"`
	Loss                int     `json:"loss"`
	Name                string  `json:"name"`
	Played              int     `json:"played"`
	Points              int     `json:"points"`
	Position            int     `json:"position"`
	ShortName           string  `json:"short_name"`
	Strength            int     `json:"strength"`
	TeamDivision        *string `json:"team_division"` // Using pointer to handle null values
	Unavailable         bool    `json:"unavailable"`
	Win                 int     `json:"win"`
	StrengthOverallHome int     `json:"strength_overall_home"`
	StrengthOverallAway int     `json:"strength_overall_away"`
	StrengthAttackHome  int     `json:"strength_attack_home"`
	StrengthAttackAway  int     `json:"strength_attack_away"`
	StrengthDefenceHome int     `json:"strength_defence_home"`
	StrengthDefenceAway int     `json:"strength_defence_away"`
	PulseID             int     `json:"pulse_id"`
}

func (t *Team) GetShortName() string {
	return t.ShortName
}

func (t *Team) GetFullName() string {
	return t.Name
}

func (t *Team) GetWinRate() float64 {
	if t.Played == 0 {
		return 0.0
	}
	return float64(t.Win) / float64(t.Played) * 100
}

func (t *Team) GetDrawRate() float64 {
	if t.Played == 0 {
		return 0.0
	}
	return float64(t.Draw) / float64(t.Played) * 100
}

func (t *Team) GetLossRate() float64 {
	if t.Played == 0 {
		return 0.0
	}
	return float64(t.Loss) / float64(t.Played) * 100
}

// This can be used to check if a team is top 4 for example.
func (t *Team) IsTopTeam(topN int) bool {
	position := t.Position
	if position != 0 {
		return t.Position <= topN
	}
	return false
}
