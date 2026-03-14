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
	playerID := 328 // Salah's ID

	fmt.Printf("Getting player details and history for ID %d...\n", playerID)
	playersResult := <-c.Players.GetAllPlayersAsync(ctx)
	historyResult := <-c.Players.GetPlayerHistoryAsync(ctx, playerID)

	if playersResult.Err != nil {
		log.Fatalf("Failed to get players: %v", playersResult.Err)
	}

	var found bool
	for _, player := range playersResult.Value {
		if player.ID != playerID {
			continue
		}
		found = true
		fmt.Println(separator)
		fmt.Printf("Mo Salah Details:\n")
		fmt.Printf("Full Name: %s\n", player.GetDisplayName())
		fmt.Printf("Current Price: £%.1f\n", player.GetPriceInPounds())
		fmt.Printf("Total Points: %d\n", player.TotalPoints)
		fmt.Printf("Form: %s\n", player.Form)
		fmt.Println(separator)
		break
	}
	if !found {
		log.Printf("Warning: Could not find player with ID %d\n", playerID)
	}

	if historyResult.Err != nil {
		log.Fatalf("Failed to get player history: %v", historyResult.Err)
	}
	history := historyResult.Value
	// Current season games
	fmt.Printf("\nLast 3 games:\n")
	fmt.Println(separator)
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
		fmt.Println(separator)
	}

	// Past seasons
	if len(history.HistoryPast) > 0 {
		fmt.Printf("\nPast Seasons:\n")
		fmt.Println(separator)
		for _, season := range history.HistoryPast {
			fmt.Printf("Season: %s\n", season.SeasonName)
			fmt.Printf("Total Points: %d\n", season.TotalPoints)
			fmt.Printf("Goals: %d, Assists: %d\n", season.GoalsScored, season.Assists)
			fmt.Println(separator)
		}
	}
}
