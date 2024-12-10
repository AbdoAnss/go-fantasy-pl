package models

import (
	"time"
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
	Started              *bool      `json:"started"`
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
	return *f.TeamAScore
}

func (f *Fixture) GetTeamHScore() int {
	return *f.TeamHScore
}

func (f *Fixture) GetGoalscorers() (map[string][]StatDetail, error) {
	goals := make(map[string][]StatDetail)
	for _, stat := range f.Stats {
		if stat.Identifier == "goals_scored" {
			goals["a"] = stat.A
			goals["h"] = stat.H
			return goals, nil
		}
	}
	return goals, nil
}

func (f *Fixture) GetAssisters() (map[string][]StatDetail, error) {
	assists := make(map[string][]StatDetail)
	for _, stat := range f.Stats {
		if stat.Identifier == "assists" {
			assists["a"] = stat.A
			assists["h"] = stat.H
			return assists, nil
		}
	}
	return assists, nil
}

func (f *Fixture) GetOwnGoalscorers() (map[string][]StatDetail, error) {
	ownGoals := make(map[string][]StatDetail)
	for _, stat := range f.Stats {
		if stat.Identifier == "own_goals" {
			ownGoals["a"] = stat.A
			ownGoals["h"] = stat.H
			return ownGoals, nil
		}
	}
	return ownGoals, nil
}

func (f *Fixture) GetYellowCards() (map[string][]StatDetail, error) {
	yellowCards := make(map[string][]StatDetail)
	for _, stat := range f.Stats {
		if stat.Identifier == "yellow_cards" {
			yellowCards["a"] = stat.A
			yellowCards["h"] = stat.H
			return yellowCards, nil
		}
	}
	return yellowCards, nil
}

func (f *Fixture) GetRedCards() (map[string][]StatDetail, error) {
	redCards := make(map[string][]StatDetail)
	for _, stat := range f.Stats {
		if stat.Identifier == "red_cards" {
			redCards["a"] = stat.A
			redCards["h"] = stat.H
			return redCards, nil
		}
	}
	return redCards, nil
}

func (f *Fixture) GetPenaltySaves() (map[string][]StatDetail, error) {
	penaltySaves := make(map[string][]StatDetail)
	for _, stat := range f.Stats {
		if stat.Identifier == "penalties_saved" {
			penaltySaves["a"] = stat.A
			penaltySaves["h"] = stat.H
			return penaltySaves, nil
		}
	}
	return penaltySaves, nil
}

func (f *Fixture) GetPenaltyMisses() (map[string][]StatDetail, error) {
	penaltyMisses := make(map[string][]StatDetail)
	for _, stat := range f.Stats {
		if stat.Identifier == "penalties_missed" {
			penaltyMisses["a"] = stat.A
			penaltyMisses["h"] = stat.H
			return penaltyMisses, nil
		}
	}
	return penaltyMisses, nil
}

func (f *Fixture) GetSaves() (map[string][]StatDetail, error) {
	saves := make(map[string][]StatDetail)
	for _, stat := range f.Stats {
		if stat.Identifier == "saves" {
			saves["a"] = stat.A
			saves["h"] = stat.H
			return saves, nil
		}
	}
	return saves, nil
}

func (f *Fixture) GetBonus() (map[string][]StatDetail, error) {
	bonus := make(map[string][]StatDetail)
	for _, stat := range f.Stats {
		if stat.Identifier == "bonus" {
			bonus["a"] = stat.A
			bonus["h"] = stat.H
			return bonus, nil
		}
	}
	return bonus, nil
}
