package main

import (
	"fmt"
	"log"

	"github.com/AbdoAnss/go-fantasy-pl/client"
)

func main() {
	c := client.NewClient()

	fixtureID := 8 // Example fixture ID

	// Get fixture details
	fmt.Printf("Getting fixture details for ID %d...\n", fixtureID)
	fixture, err := c.Fixtures.GetFixture(fixtureID)
	if err != nil {
		log.Printf("Warning: Could not get fixture details: %v\n", err)
		return
	}

	fmt.Println("----------------------------------------")
	fmt.Printf("Fixture ID: %d\n", fixture.ID)
	fmt.Printf("Team A: %d vs Team H: %d\n", fixture.TeamA, fixture.TeamH)
	fmt.Printf("Kickoff Time: %v\n", fixture.KickoffTime)
	fmt.Printf("Finished: %v\n", fixture.Finished)
	fmt.Printf("Team A Score: %v, Team H Score: %v\n", fixture.GetTeamAScore(), fixture.GetTeamHScore())
	fmt.Println("----------------------------------------")

	// Get goalscorers
	goalscorers, err := fixture.GetGoalscorers()
	if err != nil {
		log.Printf("Warning: Could not get goalscorers: %v\n", err)
	} else {
		fmt.Println("Goalscorers:")
		for team, players := range goalscorers {
			fmt.Printf("Team %s:\n", team)
			for _, player := range players {
				fmt.Printf("Player ID: %d, Goals: %d\n", player.Element, player.Value)
			}
		}
	}

	// Get assisters
	assisters, err := fixture.GetAssisters()
	if err != nil {
		log.Printf("Warning: Could not get assisters: %v\n", err)
	} else {
		fmt.Println("Assisters:")
		for team, players := range assisters {
			fmt.Printf("Team %s:\n", team)
			for _, player := range players {
				fmt.Printf("Player ID: %d, Assists: %d\n", player.Element, player.Value)
			}
		}
	}

	// Get bonus points
	bonus, err := fixture.GetBonus()
	if err != nil {
		log.Printf("Warning: Could not get bonus points: %v\n", err)
	} else {
		fmt.Println("Bonus Points:")
		for team, players := range bonus {
			fmt.Printf("Team %s:\n", team)
			for _, player := range players {
				fmt.Printf("Player ID: %d, Bonus Points: %d\n", player.Element, player.Value)
			}
		}
	}
}
