package endpoints_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllFixtures(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	fixtures, err := c.Fixtures.GetAllFixtures()
	require.NoError(t, err)
	require.Len(t, fixtures, 2)

	assert.Equal(t, 1001, fixtures[0].ID)
	require.NotNil(t, fixtures[0].Event)
	assert.Equal(t, 1, *fixtures[0].Event)
}

func TestGetFixture(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	fixture, err := c.Fixtures.GetFixture(1001)
	require.NoError(t, err)
	require.NotNil(t, fixture)

	assert.Equal(t, 1, fixture.TeamH)
	assert.Equal(t, 2, fixture.TeamA)
	assert.True(t, fixture.Started)
	assert.True(t, fixture.Finished)
	require.NotNil(t, fixture.KickoffTime)
	assert.Equal(t, time.Date(2026, time.August, 16, 19, 0, 0, 0, time.UTC), *fixture.KickoffTime)
}

func TestGetNonExistentFixture(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	fixture, err := c.Fixtures.GetFixture(999)
	require.Error(t, err)
	assert.Nil(t, fixture)
	assert.Equal(t, "fixture with ID 999 not found", err.Error())
}

func TestGetGoalscorers(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	fixture, err := c.Fixtures.GetFixture(1001)
	require.NoError(t, err)

	goalscorers, err := fixture.GetGoalscorers()
	require.NoError(t, err)
	assert.Len(t, goalscorers["h"], 1)
	assert.Len(t, goalscorers["a"], 1)
	assert.Equal(t, 101, goalscorers["h"][0].Element)
}

func TestGetAssisters(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	fixture, err := c.Fixtures.GetFixture(1001)
	require.NoError(t, err)

	assisters, err := fixture.GetAssisters()
	require.NoError(t, err)
	assert.Len(t, assisters["h"], 1)
	assert.Equal(t, 103, assisters["h"][0].Element)
}

func TestGetBonus(t *testing.T) {
	c, server := newEndpointTestClient(t)
	defer server.Close()

	fixture, err := c.Fixtures.GetFixture(1001)
	require.NoError(t, err)

	bonus, err := fixture.GetBonus()
	require.NoError(t, err)
	assert.Len(t, bonus["h"], 1)
	assert.Equal(t, 3, bonus["h"][0].Value)
}
