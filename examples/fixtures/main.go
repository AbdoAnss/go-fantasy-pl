package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/AbdoAnss/go-fantasy-pl/client"
	"github.com/AbdoAnss/go-fantasy-pl/models"
)

const separator = "----------------------------------------"

type playerInfo struct {
	Name string
}

type teamInfo struct {
	Name string
}

func buildPlayerMap(players []models.Player) map[int]playerInfo {
	playerMap := make(map[int]playerInfo, len(players))
	for _, player := range players {
		playerMap[player.ID] = playerInfo{Name: player.GetDisplayName()}
	}

	return playerMap
}

func buildTeamMap(teams []models.Team) map[int]teamInfo {
	teamMap := make(map[int]teamInfo, len(teams))
	for _, team := range teams {
		teamMap[team.ID] = teamInfo{Name: team.GetFullName()}
	}

	return teamMap
}

func sortFixtures(fixtures []models.Fixture) {
	sort.Slice(fixtures, func(i, j int) bool {
		if fixtures[i].KickoffTime == nil {
			return false
		}
		if fixtures[j].KickoffTime == nil {
			return true
		}
		return fixtures[i].KickoffTime.Before(*fixtures[j].KickoffTime)
	})
}

func isRelevantFixture(fixture models.Fixture, now time.Time) bool {
	if fixture.KickoffTime == nil {
		return false
	}

	return !fixture.KickoffTime.Before(now.Add(-7 * 24 * time.Hour))
}

func printScorers(side string, scorers []models.StatDetail, playerMap map[int]playerInfo) {
	for _, scorer := range scorers {
		if player, ok := playerMap[scorer.Element]; ok {
			fmt.Printf("  %s %s (%d)\n", side, player.Name, scorer.Value)
		}
	}
}

func printFixtureDetails(fixture models.Fixture, teamMap map[int]teamInfo, playerMap map[int]playerInfo) {
	home := teamMap[fixture.TeamH]
	away := teamMap[fixture.TeamA]

	fmt.Printf("\n%s vs %s\n", home.Name, away.Name)
	fmt.Printf("Kickoff: %s\n", fixture.KickoffTime.Format(time.RFC822))

	if fixture.Started && fixture.TeamHScore != nil && fixture.TeamAScore != nil {
		fmt.Printf("Score: %d - %d\n", *fixture.TeamHScore, *fixture.TeamAScore)

		if goals, err := fixture.GetGoalscorers(); err == nil && len(goals) > 0 {
			fmt.Println("Scorers:")
			printScorers("H", goals["h"], playerMap)
			printScorers("A", goals["a"], playerMap)
		}
	} else {
		fmt.Printf("Difficulty: Home %d - Away %d\n",
			fixture.TeamHDifficulty, fixture.TeamADifficulty)
	}

	fmt.Println(separator)
}

func main() {
	fpl, err := client.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get core data concurrently
	playersCh := fpl.Players.GetAllPlayersAsync(ctx)
	teamsCh := fpl.Teams.GetAllTeamsAsync(ctx)
	fixturesCh := fpl.Fixtures.GetAllFixturesAsync(ctx)

	playersResult := <-playersCh
	teamsResult := <-teamsCh
	fixturesResult := <-fixturesCh

	if playersResult.Err != nil {
		log.Fatalf("Failed to get players: %v", playersResult.Err)
	}
	if teamsResult.Err != nil {
		log.Fatalf("Failed to get teams: %v", teamsResult.Err)
	}
	if fixturesResult.Err != nil {
		log.Fatalf("Failed to get fixtures: %v", fixturesResult.Err)
	}

	playerMap := buildPlayerMap(playersResult.Value)
	teamMap := buildTeamMap(teamsResult.Value)
	fixtures := fixturesResult.Value
	sortFixtures(fixtures)

	fmt.Println("Recent and Upcoming Fixtures")
	fmt.Println("============================")

	now := time.Now()
	analyzed := 0

	for _, fix := range fixtures {
		if !isRelevantFixture(fix, now) {
			continue
		}

		printFixtureDetails(fix, teamMap, playerMap)

		analyzed++
		if analyzed >= 5 {
			break
		}
	}
}
