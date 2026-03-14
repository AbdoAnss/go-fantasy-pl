# go-fantasy-pl

Simple Go client for the Fantasy Premier League API.

It supports players, teams, fixtures, managers, leagues, caching, and rate limiting. For workloads that need multiple independent datasets, the async helpers are now the best option because they let you fetch data concurrently.

## Install

```bash
go get github.com/AbdoAnss/go-fantasy-pl
```

## What you get

- Typed models for FPL responses
- Simple client setup
- Built-in rate limiting
- In-memory cache by default
- Optional Redis cache
- Async helpers for concurrent reads

## Quick start

```go
package main

import (
    "fmt"
    "log"

    "github.com/AbdoAnss/go-fantasy-pl/client"
)

func main() {
    c, err := client.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    teams, err := c.Teams.GetAllTeams()
    if err != nil {
        log.Fatal(err)
    }

    players, err := c.Players.GetAllPlayers()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("teams: %d\n", len(teams))
    fmt.Printf("players: %d\n", len(players))
}
```

## Async usage

If you want the fastest path when fetching players, teams, and fixtures together, use the async methods and wait on the result channels.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/AbdoAnss/go-fantasy-pl/client"
)

func main() {
    c, err := client.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    playersCh := c.Players.GetAllPlayersAsync(ctx)
    teamsCh := c.Teams.GetAllTeamsAsync(ctx)
    fixturesCh := c.Fixtures.GetAllFixturesAsync(ctx)

    players := <-playersCh
    teams := <-teamsCh
    fixtures := <-fixturesCh

    if players.Err != nil {
        log.Fatal(players.Err)
    }
    if teams.Err != nil {
        log.Fatal(teams.Err)
    }
    if fixtures.Err != nil {
        log.Fatal(fixtures.Err)
    }

    fmt.Printf("players: %d, teams: %d, fixtures: %d\n",
        len(players.Value), len(teams.Value), len(fixtures.Value))
}
```

Available async helpers:

- `GetAllPlayersAsync(ctx)`
- `GetPlayerHistoryAsync(ctx, playerID)`
- `GetPlayerHistoriesBatch(ctx, ids)`
- `GetAllTeamsAsync(ctx)`
- `GetAllFixturesAsync(ctx)`

## Cache

Default cache is in-memory. Redis-backed caching is also supported when you want shared cache across multiple instances.

## Examples

See `examples/` for working examples.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)
