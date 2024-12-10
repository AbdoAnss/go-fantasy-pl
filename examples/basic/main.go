package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/AbdoAnss/go-fantasy-pl/client"
	"github.com/AbdoAnss/go-fantasy-pl/models"
)

func main() {
	// Step 1: Initialize the FPL client
	fpl := client.NewClient()

	// Step 2: Get all teams and create a team map for quick lookups
	teams, err := fpl.Teams.GetAllTeams()
	if err != nil {
		log.Fatalf("Failed to get teams: %v", err)
	}
	teamMap := make(map[int]models.Team)
	for _, team := range teams {
		teamMap[team.ID] = team
	}

	// Step 3: Analyze team strengths
	fmt.Println("Top 5 Teams by Overall Strength:")
	sortedTeams := make([]models.Team, len(teams))
	copy(sortedTeams, teams)
	sort.Slice(sortedTeams, func(i, j int) bool {
		return sortedTeams[i].Strength > sortedTeams[j].Strength
	})
	for i, team := range sortedTeams[:5] {
		fmt.Printf("%d. %s (Strength: %d)\n", i+1, team.GetFullName(), team.Strength)
	}
	fmt.Println("----------------------------------------")

	// Step 4: Get all players and find top performers
	players, err := fpl.Players.GetAllPlayers()
	if err != nil {
		log.Fatalf("Failed to get players: %v", err)
	}

	// Sort players by total points
	sort.Slice(players, func(i, j int) bool {
		return players[i].TotalPoints > players[j].TotalPoints
	})

	// Step 5: Display top 5 players
	fmt.Println("Top 5 Players by Total Points:")
	for i, player := range players[:5] {
		team := teamMap[player.Team]
		fmt.Printf("%d. %s (%s) - £%.1fm - %d points\n",
			i+1,
			player.GetDisplayName(),
			team.GetShortName(),
			player.GetPriceInPounds(),
			player.TotalPoints)
	}
	fmt.Println("----------------------------------------")

	// Step 6: Get upcoming fixtures
	fixtures, err := fpl.Fixtures.GetAllFixtures()
	if err != nil {
		log.Fatalf("Failed to get fixtures: %v", err)
	}

	// Filter for upcoming fixtures
	var upcomingFixtures []models.Fixture
	for _, fix := range fixtures {
		if fix.Started && fix.Finished {
			upcomingFixtures = append(upcomingFixtures, fix)
		}
	}

	// Step 7: Display next 5 fixtures with team strengths
	fmt.Println("Next 5 Fixtures (with Team Strengths):")
	for i, fix := range upcomingFixtures[:5] {
		homeTeam := teamMap[fix.TeamH]
		awayTeam := teamMap[fix.TeamA]
		fmt.Printf("%d. %s (%d) vs %s (%d) - Gameweek %d\n",
			i+1,
			homeTeam.GetShortName(), homeTeam.Strength,
			awayTeam.GetShortName(), awayTeam.Strength,
			fix.Event)
	}
	fmt.Println("----------------------------------------")

	// Step 8: Get detailed history for a top player
	topPlayer := players[0]
	history, err := fpl.Players.GetPlayerHistory(topPlayer.ID)
	if err != nil {
		log.Fatalf("Failed to get player history: %v", err)
	}

	fmt.Printf("Recent Performance - %s:\n", topPlayer.GetDisplayName())
	sort.Slice(history.History, func(i, j int) bool {
		return history.History[i].Round > history.History[j].Round
	})

	for i, game := range history.History[:3] {
		fmt.Printf("%d. GW%d: %d pts (Goals: %d, Assists: %d)\n",
			i+1,
			game.Round,
			game.TotalPoints,
			game.GoalsScored,
			game.Assists)
	}
}
