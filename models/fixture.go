package models

import (
	"time"
)

const (
	StatGoalsScored     = "goals_scored"
	StatAssists         = "assists"
	StatOwnGoals        = "own_goals"
	StatYellowCards     = "yellow_cards"
	StatRedCards        = "red_cards"
	StatPenaltiesSaved  = "penalties_saved"
	StatPenaltiesMissed = "penalties_missed"
	StatSaves           = "saves"
	StatBonus           = "bonus"
)

type Stat struct {
	Identifier string       `json:"identifier"`
	A          []StatDetail `json:"a"`
	H          []StatDetail `json:"h"`
}

type StatDetail struct {
	Value   int `json:"value"`
	Element int `json:"element"`
}

type Fixture struct {
	Code                 int        `json:"code"`
	Event                *int       `json:"event"`
	Finished             bool       `json:"finished"`
	FinishedProvisional  bool       `json:"finished_provisional"`
	ID                   int        `json:"id"`
	KickoffTime          *time.Time `json:"kickoff_time"`
	Minutes              int        `json:"minutes"`
	ProvisionalStartTime bool       `json:"provisional_start_time"`
	Started              bool       `json:"started"`
	TeamA                int        `json:"team_a"`
	TeamAScore           *int       `json:"team_a_score"`
	TeamH                int        `json:"team_h"`
	TeamHScore           *int       `json:"team_h_score"`
	Stats                []Stat     `json:"stats"`
	TeamHDifficulty      int        `json:"team_h_difficulty"`
	TeamADifficulty      int        `json:"team_a_difficulty"`
	PulseID              int        `json:"pulse_id"`
}

func (f *Fixture) GetTeamAScore() int {
	if f.TeamAScore == nil {
		return 0
	}
	return *f.TeamAScore
}

func (f *Fixture) GetTeamHScore() int {
	if f.TeamHScore == nil {
		return 0
	}
	return *f.TeamHScore
}

func (f *Fixture) getStat(identifier string) (map[string][]StatDetail, error) {
	result := make(map[string][]StatDetail)
	for _, stat := range f.Stats {
		if stat.Identifier == identifier {
			result["a"] = stat.A
			result["h"] = stat.H
			return result, nil
		}
	}
	return result, nil
}

func (f *Fixture) GetGoalscorers() (map[string][]StatDetail, error) {
	return f.getStat(StatGoalsScored)
}

func (f *Fixture) GetAssisters() (map[string][]StatDetail, error) {
	return f.getStat(StatAssists)
}

func (f *Fixture) GetOwnGoalscorers() (map[string][]StatDetail, error) {
	return f.getStat(StatOwnGoals)
}

func (f *Fixture) GetYellowCards() (map[string][]StatDetail, error) {
	return f.getStat(StatYellowCards)
}

func (f *Fixture) GetRedCards() (map[string][]StatDetail, error) {
	return f.getStat(StatRedCards)
}

func (f *Fixture) GetPenaltySaves() (map[string][]StatDetail, error) {
	return f.getStat(StatPenaltiesSaved)
}

func (f *Fixture) GetPenaltyMisses() (map[string][]StatDetail, error) {
	return f.getStat(StatPenaltiesMissed)
}

func (f *Fixture) GetSaves() (map[string][]StatDetail, error) {
	return f.getStat(StatSaves)
}

func (f *Fixture) GetBonus() (map[string][]StatDetail, error) {
	return f.getStat(StatBonus)
}
