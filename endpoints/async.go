package endpoints

import (
	"context"
	"sync"

	"github.com/AbdoAnss/go-fantasy-pl/models"
)

// Result is a generic wrapper for a value and an error,
// used for returning results from asynchronous operations.
type Result[T any] struct {
	Value T
	Err   error
}

// PlayerHistoryResult specifically wraps a player's ID and their history,
// designed for use in batch fetch operations where identifying the player is essential.
type PlayerHistoryResult struct {
	PlayerID int
	History  *models.PlayerHistory
	Err      error
}

// GetAllPlayersAsync fetches all players concurrently and returns a channel
// that receives a single Result containing all players or an error.
func (ps *PlayerService) GetAllPlayersAsync(ctx context.Context) <-chan Result[[]models.Player] {
	ch := make(chan Result[[]models.Player], 1)
	go func() {
		defer close(ch)
		players, err := ps.GetAllPlayers()
		select {
		case ch <- Result[[]models.Player]{Value: players, Err: err}:
		case <-ctx.Done():
		}
	}()
	return ch
}

// GetPlayerHistoryAsync fetches the history for a single player asynchronously
// and returns a channel that receives the result.
func (ps *PlayerService) GetPlayerHistoryAsync(ctx context.Context, id int) <-chan Result[*models.PlayerHistory] {
	ch := make(chan Result[*models.PlayerHistory], 1)
	go func() {
		defer close(ch)
		history, err := ps.GetPlayerHistory(id)
		select {
		case ch <- Result[*models.PlayerHistory]{Value: history, Err: err}:
		case <-ctx.Done():
		}
	}()
	return ch
}

// GetPlayerHistoriesBatch fetches player histories concurrently for multiple player IDs.
// Results are sent to the returned channel as they complete.
func (ps *PlayerService) GetPlayerHistoriesBatch(ctx context.Context, ids []int) <-chan PlayerHistoryResult {
	ch := make(chan PlayerHistoryResult, len(ids))
	var wg sync.WaitGroup

	for _, id := range ids {
		wg.Add(1)
		go func(playerID int) {
			defer wg.Done()
			history, err := ps.GetPlayerHistory(playerID)
			select {
			case ch <- PlayerHistoryResult{PlayerID: playerID, History: history, Err: err}:
			case <-ctx.Done():
			}
		}(id)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

// GetAllFixturesAsync fetches all fixtures asynchronously and returns a channel
// that receives a single Result containing all fixtures or an error.
func (fs *FixtureService) GetAllFixturesAsync(ctx context.Context) <-chan Result[[]models.Fixture] {
	ch := make(chan Result[[]models.Fixture], 1)
	go func() {
		defer close(ch)
		fixtures, err := fs.GetAllFixtures()
		select {
		case ch <- Result[[]models.Fixture]{Value: fixtures, Err: err}:
		case <-ctx.Done():
		}
	}()
	return ch
}

// GetAllTeamsAsync fetches all teams asynchronously and returns a channel
// that receives a single Result containing all teams or an error.
func (ts *TeamService) GetAllTeamsAsync(ctx context.Context) <-chan Result[[]models.Team] {
	ch := make(chan Result[[]models.Team], 1)
	go func() {
		defer close(ch)
		teams, err := ts.GetAllTeams()
		select {
		case ch <- Result[[]models.Team]{Value: teams, Err: err}:
		case <-ctx.Done():
		}
	}()
	return ch
}
