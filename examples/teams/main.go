package main

import (
	"fmt"
	"log"

	"github.com/AbdoAnss/go-fantasy-pl/client"
)

func main() {
	c := client.NewClient()

	teamID := 8 // Example team ID : Everton

	// Get team details
	fmt.Printf("Getting team details for ID %d...\n", teamID)
	team, err := c.Teams.GetTeam(teamID)
	if err != nil {
		log.Printf("Warning: Could not get team details: %v\n", err)
		return
	}

	fmt.Println("----------------------------------------")
	fmt.Printf("Team ID: %d\n", team.ID)
	fmt.Printf("Team Name: %s\n", team.GetFullName())
	fmt.Printf("Short Name: %s\n", team.GetShortName())
	fmt.Printf("Points: %d\n", team.Points)
	fmt.Printf("Played: %d\n", team.Played)
	fmt.Printf("Wins: %d\n", team.Win)
	fmt.Printf("Draws: %d\n", team.Draw)
	fmt.Printf("Losses: %d\n", team.Loss)
	fmt.Printf("Position: %d\n", team.Position)
	fmt.Printf("Strength: %d\n", team.Strength)
	fmt.Printf("Win Rate: %.2f%%\n", team.GetWinRate())
	fmt.Printf("Draw Rate: %.2f%%\n", team.GetDrawRate())
	fmt.Printf("Loss Rate: %.2f%%\n", team.GetLossRate())

	// Check if the team is in the top 4
	topN := 4
	if team.IsTopTeam(topN) {
		fmt.Printf("The team is in the top %d.\n", topN)
	} else {
		fmt.Printf("The team is not in the top %d.\n", topN)
	}

	fmt.Println("----------------------------------------")
}
