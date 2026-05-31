package endpoints_test

import (
	"testing"

	"github.com/AbdoAnss/go-fantasy-pl/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManagerEndpoints(t *testing.T) {
	skipUnlessLive(t)

	testManagerClient, err := client.NewClient(client.WithMemoryCache())
	require.NoError(t, err)

	t.Run("GetManager", func(t *testing.T) {
		managerID := 1387812

		manager, err := testManagerClient.Managers.GetManager(managerID)
		assert.NoError(t, err)
		assert.NotNil(t, manager)
	})

	t.Run("GetNonExistentManager", func(t *testing.T) {
		manager, err := testManagerClient.Managers.GetManager(99999999)
		assert.Error(t, err)
		assert.Nil(t, manager)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("GetCurrentTeam", func(t *testing.T) {
		managerID := 1387812

		team, err := testManagerClient.Managers.GetCurrentTeam(managerID)
		assert.NoError(t, err)
		assert.NotNil(t, team)
		assert.NotEmpty(t, team.Picks)
	})

	t.Run("GetManagerHistory", func(t *testing.T) {
		managerID := 1387812

		history, err := testManagerClient.Managers.GetManagerHistory(managerID)
		assert.NoError(t, err)
		assert.NotNil(t, history)
	})

	t.Run("CacheConsistency", func(t *testing.T) {
		managerID := 1387812

		manager1, err := testManagerClient.Managers.GetManager(managerID)
		require.NoError(t, err)

		manager2, err := testManagerClient.Managers.GetManager(managerID)
		require.NoError(t, err)

		assert.Equal(t, manager1.ID, manager2.ID)
		assert.Equal(t, manager1.Name, manager2.Name)
	})
}
