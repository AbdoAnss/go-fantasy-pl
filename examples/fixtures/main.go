package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/AbdoAnss/go-fantasy-pl/client"
)

func main() {
	fpl := client.NewClient()

	// Get core data
	players, err := fpl.Players.GetAllPlayers()
	if err != nil {
		log.Fatalf("Failed to get players: %v", err)
	}

	teams, err := fpl.Teams.GetAllTeams()
	if err != nil {
		log.Fatalf("Failed to get teams: %v", err)
	}

	fixtures, err := fpl.Fixtures.GetAllFixtures()
	if err != nil {
		log.Fatalf("Failed to get fixtures: %v", err)
	}

	// Create lookup maps
	playerMap := make(map[int]struct {
		ID   int
		Name string
	})
	for _, p := range players {
		playerMap[p.ID] = struct {
			ID   int
			Name string
		}{p.ID, p.GetDisplayName()}
	}

	teamMap := make(map[int]struct {
		ID   int
		Name string
	})
	for _, t := range teams {
		teamMap[t.ID] = struct {
			ID   int
			Name string
		}{t.ID, t.GetFullName()}
	}

	// Sort fixtures with nil-safe comparison
	sort.Slice(fixtures, func(i, j int) bool {
		// Handle nil KickoffTime
		if fixtures[i].KickoffTime == nil {
			return false
		}
		if fixtures[j].KickoffTime == nil {
			return true
		}
		return fixtures[i].KickoffTime.Before(*fixtures[j].KickoffTime)
	})

	fmt.Println("Recent and Upcoming Fixtures")
	fmt.Println("============================")

	now := time.Now()
	analyzed := 0

	for _, fix := range fixtures {
		// Skip fixtures without kickoff time
		if fix.KickoffTime == nil {
			continue
		}

		// Skip old matches
		if fix.KickoffTime.Before(now.Add(-7 * 24 * time.Hour)) {
			continue
		}

		home := teamMap[fix.TeamH]
		away := teamMap[fix.TeamA]

		fmt.Printf("\n%s vs %s\n", home.Name, away.Name)
		fmt.Printf("Kickoff: %s\n", fix.KickoffTime.Format(time.RFC822))

		if fix.Started && fix.TeamHScore != nil && fix.TeamAScore != nil {
			fmt.Printf("Score: %d - %d\n", *fix.TeamHScore, *fix.TeamAScore)

			// Print scorers if available
			if goals, err := fix.GetGoalscorers(); err == nil && len(goals) > 0 {
				fmt.Println("Scorers:")
				for _, scorer := range goals["h"] {
					if player, ok := playerMap[scorer.Element]; ok {
						fmt.Printf("  %s (%d)\n", player.Name, scorer.Value)
					}
				}
				for _, scorer := range goals["a"] {
					if player, ok := playerMap[scorer.Element]; ok {
						fmt.Printf("  %s (%d)\n", player.Name, scorer.Value)
					}
				}
			}
		} else {
			fmt.Printf("Difficulty: Home %d - Away %d\n",
				fix.TeamHDifficulty, fix.TeamADifficulty)
		}

		fmt.Println("----------------------------------------")

		analyzed++
		if analyzed >= 5 {
			break
		}
	}
}
