package models

import (
	"time"
)

type ChipPlay struct {
	ChipName  string `json:"chip_name"`
	NumPlayed int    `json:"num_played"`
}

type TopElementInfo struct {
	ID     int `json:"id"`
	Points int `json:"points"`
}

type GameWeek struct {
	ID                     int            `json:"id"`
	Name                   string         `json:"name"`
	DeadlineTime           time.Time      `json:"deadline_time"`
	ReleaseTime            *time.Time     `json:"release_time"` // Pointer to handle null
	AverageEntryScore      int            `json:"average_entry_score"`
	Finished               bool           `json:"finished"`
	DataChecked            bool           `json:"data_checked"`
	HighestScoringEntry    int            `json:"highest_scoring_entry"`
	DeadlineTimeEpoch      int64          `json:"deadline_time_epoch"`
	DeadlineTimeGameOffset int            `json:"deadline_time_game_offset"`
	HighestScore           int            `json:"highest_score"`
	IsPrevious             bool           `json:"is_previous"`
	IsCurrent              bool           `json:"is_current"`
	IsNext                 bool           `json:"is_next"`
	CupLeaguesCreated      bool           `json:"cup_leagues_created"`
	H2hKoMatchesCreated    bool           `json:"h2h_ko_matches_created"`
	RankedCount            int            `json:"ranked_count"`
	ChipPlays              []ChipPlay     `json:"chip_plays"`
	MostSelected           int            `json:"most_selected"`
	MostTransferredIn      int            `json:"most_transferred_in"`
	TopElement             int            `json:"top_element"`
	TopElementInfo         TopElementInfo `json:"top_element_info"`
	TransfersMade          int            `json:"transfers_made"`
	MostCaptained          int            `json:"most_captained"`
	MostViceCaptained      int            `json:"most_vice_captained"`
}

func (gw *GameWeek) GetChipPlayCount(chipName string) int {
	for _, chipPlay := range gw.ChipPlays {
		if chipPlay.ChipName == chipName {
			return chipPlay.NumPlayed
		}
	}
	return 0
}

func (gw *GameWeek) IsFinished() bool {
	return gw.Finished
}

func (gw *GameWeek) GetTopElementInfo() TopElementInfo {
	return gw.TopElementInfo
}
