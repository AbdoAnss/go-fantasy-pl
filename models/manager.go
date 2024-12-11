package models

import "time"

// TODO: Implement Manager leagues

type Manager struct {
	ID                         *int      `json:"id"`
	JoinedTime                 time.Time `json:"joined_time"`
	StartedEvent               int       `json:"started_event"`
	FavouriteTeam              int       `json:"favourite_team"`
	PlayerFirstName            string    `json:"player_first_name"`
	PlayerLastName             string    `json:"player_last_name"`
	PlayerRegionID             int       `json:"player_region_id"`
	PlayerRegionName           string    `json:"player_region_name"`
	PlayerRegionISOCodesShort  string    `json:"player_region_iso_code_short"`
	PlayerRegionISOCodesLong   string    `json:"player_region_iso_code_long"`
	YearsActive                int       `json:"years_active"`
	SummaryOverallPoints       int       `json:"summary_overall_points"`
	SummaryOverallRank         int       `json:"summary_overall_rank"`
	SummaryEventPoints         int       `json:"summary_event_points"`
	SummaryEventRank           int       `json:"summary_event_rank"`
	CurrentEvent               *int      `json:"current_event"`
	Name                       string    `json:"name"`
	NameChangeBlocked          bool      `json:"name_change_blocked"`
	EnteredEvents              []int     `json:"entered_events"`
	Kit                        *string   `json:"kit"`
	LastDeadlineBank           int       `json:"last_deadline_bank"`
	LastDeadlineValue          int       `json:"last_deadline_value"`
	LastDeadlineTotalTransfers int       `json:"last_deadline_total_transfers"`
}

func (m *Manager) GetFullName() string {
	return m.PlayerFirstName + " " + m.PlayerLastName
}
