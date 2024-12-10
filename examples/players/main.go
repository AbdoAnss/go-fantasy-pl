package main

import (
	"fmt"
	"log"

	"github.com/AbdoAnss/go-fantasy-pl/client"
)

func main() {
	c := client.NewClient()
	playerID := 328 // Salah's ID

	fmt.Printf("Getting player details for ID %d...\n", playerID)
	player, err := c.Players.GetPlayer(playerID)
	if err != nil {
		log.Printf("Warning: Could not get player details: %v\n", err)
	} else {
		fmt.Println("----------------------------------------")
		fmt.Printf("Mo Salah Details:\n")
		fmt.Printf("Full Name: %s\n", player.GetDisplayName())
		fmt.Printf("Current Price: £%.1f\n", player.GetPriceInPounds())
		fmt.Printf("Total Points: %d\n", player.TotalPoints)
		fmt.Printf("Form: %s\n", player.Form)
		fmt.Println("----------------------------------------")
	}

	// Get history
	fmt.Printf("Getting player history...\n")
	history, err := c.Players.GetPlayerHistory(playerID)
	if err != nil {
		log.Fatalf("Failed to get player history: %v", err)
	}
	// Current season games
	fmt.Printf("\nLast 3 games:\n")
	fmt.Println("----------------------------------------")
	recentGames := history.History
	if len(recentGames) > 3 {
		recentGames = recentGames[len(recentGames)-3:]
	}
	for _, game := range recentGames {
		fmt.Printf("Gameweek %d:\n", game.Round)
		fmt.Printf("Against: Team %d (Home: %v)\n", game.OpponentTeam, game.WasHome)
		fmt.Printf("Points: %d\n", game.TotalPoints)
		fmt.Printf("Goals: %d, Assists: %d\n", game.GoalsScored, game.Assists)
		fmt.Printf("Minutes: %d\n", game.Minutes)
		fmt.Printf("Bonus Points: %d\n", game.Bonus)
		fmt.Println("----------------------------------------")
	}

	// Past seasons
	if len(history.HistoryPast) > 0 {
		fmt.Printf("\nPast Seasons:\n")
		fmt.Println("----------------------------------------")
		for _, season := range history.HistoryPast {
			fmt.Printf("Season: %s\n", season.SeasonName)
			fmt.Printf("Total Points: %d\n", season.TotalPoints)
			fmt.Printf("Goals: %d, Assists: %d\n", season.GoalsScored, season.Assists)
			fmt.Println("----------------------------------------")
		}
	}
}
