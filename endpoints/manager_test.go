package endpoints_test

import (
	"testing"

	"github.com/AbdoAnss/go-fantasy-pl/client"
	"github.com/stretchr/testify/assert"
)

var testManagerClient *client.Client

func init() {
	testManagerClient = client.NewClient()
}

func TestManagerEndpoints(t *testing.T) {
	t.Run("GetManager", func(t *testing.T) {
		// Using a known manager ID (you can replace with a valid one)
		managerID := 1387812

		manager, err := testManagerClient.Managers.GetManager(managerID)
		assert.NoError(t, err, "expected no error when getting manager")
		assert.NotNil(t, manager, "expected manager to be returned")

		// Log manager details
		t.Logf("Manager Details:")
		t.Logf("ID: %d", *manager.ID)
		t.Logf("Name: %s", manager.GetFullName())
		t.Logf("Team Name: %s", manager.Name)
		t.Logf("Overall Points: %d", manager.SummaryOverallPoints)
		t.Logf("Overall Rank: %d", manager.SummaryOverallRank)
	})

	t.Run("GetNonExistentManager", func(t *testing.T) {
		manager, err := testManagerClient.Managers.GetManager(99999999)
		assert.Error(t, err, "expected error when getting non-existent manager")
		assert.Nil(t, manager, "expected nil manager for non-existent ID")
		assert.Contains(t, err.Error(), "not found", "expected 'not found' error message")
	})

	t.Run("GetCurrentTeam", func(t *testing.T) {
		// Using same manager ID
		managerID := 1387812

		team, err := testManagerClient.Managers.GetCurrentTeam(managerID)
		assert.NoError(t, err, "expected no error when getting current team")
		assert.NotNil(t, team, "expected team to be returned")

		// Log team details
		t.Logf("Current Team Details:")
		t.Logf("Number of Picks: %d", len(team.Picks))

		// Log starting XI
		t.Log("Starting XI:")
		for _, pick := range team.GetStartingXI() {
			t.Logf("Position %d: Player ID %d (Captain: %v)",
				pick.Position,
				pick.Element,
				pick.IsCaptain)
		}

		// Log bench
		t.Log("Bench:")
		for _, pick := range team.GetBench() {
			t.Logf("Position %d: Player ID %d",
				pick.Position,
				pick.Element)
		}

		t.Logf("Team Value: £%.1fm", team.GetTeamValueInMillions())
		t.Logf("Bank: £%.1fm", team.GetBankValueInMillions())
	})

	t.Run("GetManagerHistory", func(t *testing.T) {
		// Using same manager ID
		managerID := 1387812

		history, err := testManagerClient.Managers.GetManagerHistory(managerID)
		assert.NoError(t, err, "expected no error when getting manager history")
		assert.NotNil(t, history, "expected history to be returned")

		// Log history details
		t.Logf("Manager History Details:")

		if len(history.Current) > 0 {
			t.Log("Current Season Performance:")
			for _, gw := range history.Current[:3] { // Show first 3 gameweeks
				t.Logf("GW%d: Points: %d, Overall Rank: %d",
					gw.Event,
					gw.Points,
					gw.OverallRank)
			}
		}

		if len(history.Past) > 0 {
			t.Log("Past Seasons:")
			for _, season := range history.Past {
				t.Logf("Season %s: Points: %d, Overall Rank: %d",
					season.SeasonName,
					season.TotalPoints,
					season.Rank)
			}
		}
	})

	t.Run("CacheConsistency", func(t *testing.T) {
		managerID := 1387812

		// First call
		manager1, err := testManagerClient.Managers.GetManager(managerID)
		assert.NoError(t, err)

		// Second call (should be from cache)
		manager2, err := testManagerClient.Managers.GetManager(managerID)
		assert.NoError(t, err)

		// Compare results
		assert.Equal(t, manager1.ID, manager2.ID, "cached manager should match original")
		assert.Equal(t, manager1.Name, manager2.Name, "cached manager name should match original")
	})
}
