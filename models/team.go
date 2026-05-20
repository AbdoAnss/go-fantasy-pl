package models

// Team represents a Premier League team in the FPL system.
// It includes metadata, performance statistics, and strength ratings.
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

// GetShortName returns the 3-letter abbreviation of the team (e.g., "ARS").
func (t *Team) GetShortName() string {
	return t.ShortName
}

// GetFullName returns the full name of the team (e.g., "Arsenal").
func (t *Team) GetFullName() string {
	return t.Name
}

// GetWinRate calculates the percentage of games won by the team.
func (t *Team) GetWinRate() float64 {
	if t.Played == 0 {
		return 0.0
	}
	return float64(t.Win) / float64(t.Played) * 100
}

// GetDrawRate calculates the percentage of games drawn by the team.
func (t *Team) GetDrawRate() float64 {
	if t.Played == 0 {
		return 0.0
	}
	return float64(t.Draw) / float64(t.Played) * 100
}

// GetLossRate calculates the percentage of games lost by the team.
func (t *Team) GetLossRate() float64 {
	if t.Played == 0 {
		return 0.0
	}
	return float64(t.Loss) / float64(t.Played) * 100
}

// IsTopTeam checks if the team is currently ranked within the top N positions.
func (t *Team) IsTopTeam(topN int) bool {
	position := t.Position
	if position != 0 {
		return t.Position <= topN
	}
	return false
}
