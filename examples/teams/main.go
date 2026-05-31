package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AbdoAnss/go-fantasy-pl/client"
)

const separator = "----------------------------------------"

func main() {
	c, err := client.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get team details through the async teams list endpoint
	teamsResult := <-c.Teams.GetAllTeamsAsync(ctx)
	if teamsResult.Err != nil {
		log.Printf("Warning: Could not get teams: %v\n", teamsResult.Err)
		return
	}

	if len(teamsResult.Value) == 0 {
		log.Print("Warning: no teams returned by the API")
		return
	}

	team := teamsResult.Value[0]
	for _, candidate := range teamsResult.Value[1:] {
		if candidate.Strength > team.Strength {
			team = candidate
		}
	}

	fmt.Printf("Showing current strongest team: %s (ID %d)\n", team.GetFullName(), team.ID)

	fmt.Println(separator)
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

	fmt.Println(separator)
}
