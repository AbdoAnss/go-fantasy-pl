# Fantasy Premier League API Wrapper for Go

## Overview
This project provides a Go wrapper for the Fantasy Premier League (FPL) API, offering a type-safe and idiomatic way to interact with FPL data. The wrapper focuses on read operations for essential FPL data.

## Installation
```bash
go get github.com/abdoanss/go-fantasy-pl
```

## Key Features
- Type-safe access to FPL data
- Rate limiting handling
- Response caching
- Easy-to-use client interface

## Available Data
The wrapper currently provides access to:
- Player statistics and information
- Team data
- Fixture information
- General game settings

## Project Structure
```
go-fantasy-pl/
├── client/             # Core HTTP client implementation
│   ├── client.go       # Main client struct and configuration
│   ├── options.go      # Configuration options for the client
│   └── rate_limiter.go # Rate limiting implementation
├── models/             # Data structures for API responses
│   ├── player.go       # Player-related structs
│   ├── team.go         # Team-related structs
│   └── fixture.go      # Fixture-related structs
├── endpoints/          # API endpoint implementations
│   ├── bootstrap.go    # General game data
│   ├── players.go      # Player-related endpoints
│   ├── teams.go        # Team-related endpoints
│   └── fixtures.go     # Match fixtures
├── internal/           # Internal packages
│   └── cache/          # Caching functionality
└── examples/           # Usage examples
```

## Quick Start
```go
package main

import (
    "fmt"
    "github.com/abdoanss/go-fantasy-pl/client"
)

func main() {
    // Initialize client
    fpl := client.NewClient()

    // Get all teams
    teams, err := fpl.Teams.GetAllTeams()
    if err != nil {
        panic(err)
    }

    // Print team names
    for _, team := range teams {
        fmt.Printf("Team: %s\n", team.GetFullName())
    }

    // Get all players
    players, err := fpl.Players.GetAllPlayers()
    if err != nil {
        panic(err)
    }

    // Print top 5 players by points
    fmt.Println("\nTop 5 Players:")
    for i, player := range players[:5] {
        fmt.Printf("%d. %s - Points: %d\n", 
            i+1, player.GetDisplayName(), player.TotalPoints)
    }
}
```

## Features

### Client
- Configurable HTTP client
- Built-in rate limiting
- Automatic request retries
- Response caching

### Available Data Types
- Teams: Full team information and statistics
- Players: Detailed player data and performance stats
- Fixtures: Match information and results
- Game Settings: General FPL game configuration

### Caching
The wrapper includes an optional caching system with configurable TTL for:
- Team data (24 hours)
- Player data (10 minutes)
- Fixture data (10 minutes)
- Game settings (24 hours)

### Rate Limiting
Automatic rate limiting is implemented to prevent exceeding FPL's API limits:
- Default: 50 requests per minute
- Configurable through client options
- Automatic request queuing

## Examples
Check the `examples/` directory for complete usage examples, including:
- Basic data retrieval
- Fixture analysis
- Player statistics
- Team information

## Contributing
Contributions are welcome! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

## Error Handling
The wrapper provides detailed error messages for:
- API errors
- Network issues
- Invalid responses
- Rate limiting

## License
[MIT License](./LICENSE)

## Version
Current version: v0.1.0

## Limitations
- Read-only access (no team management operations)
- Subject to FPL API rate limits
- Some data may be delayed based on FPL API updates
