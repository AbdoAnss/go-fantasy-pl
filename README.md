# Fantasy Premier League API Wrapper for Go

## Overview
This project provides a Go wrapper for the Fantasy Premier League (FPL) API, offering a type-safe and idiomatic way to interact with FPL data. The wrapper focuses on read operations only, as team management operations are restricted to the official FPL application.

## Key Features
- Type-safe access to FPL data
- Comprehensive coverage of public API endpoints
- Rate limiting handling
- Optional response caching
- Concurrent request support for batch operations
- Extensive documentation and examples

## Scope
The wrapper provides access to the following FPL data:
- Player statistics and information
- Team data
- Gameweek details
- League standings
- Fixture information
- Historical performance data
- General game settings and rules

## Architecture
The project follows a monolithic package architecture, organized into logical components:

```
go-fantasy-pl/
├── client/          # Core HTTP client implementation
│   ├── client.go    # Main client struct and configuration
│   └── rate.go      # Rate limiting implementation
├── models/          # Data structures for API responses
│   ├── player.go    # Player-related structs
│   ├── team.go      # Team-related structs
│   └── gameweek.go  # Gameweek-related structs
├── endpoints/       # API endpoint implementations
│   ├── bootstrap.go # General game data
│   ├── fixtures.go  # Match fixtures
│   └── leagues.go   # League standings
├── utils/          # Helper functions
│   ├── cache.go    # Caching implementation
│   └── helpers.go  # General utility functions
└── examples/       # Usage examples
    └── basic.go    # Basic usage patterns
```

## Main Components

### Client
- Handles HTTP communication with the FPL API
- Manages rate limiting and request throttling
- Implements retry logic for failed requests
- Provides configuration options for timeouts and concurrency

### Models
- Defines Go structs that match FPL API response structures
- Implements custom JSON unmarshaling where needed
- Provides helper methods for common data operations
- Includes validation logic for response data

### Endpoints
- Implements methods for each supported API endpoint
- Groups related endpoints into logical packages
- Handles parameter validation and request formation
- Returns strongly-typed responses

### Utils
- Provides caching mechanisms to reduce API load
- Implements helper functions for common operations
- Handles error wrapping and context management

## Usage Examples

```go
// Initialize client
client := fpl.NewClient()

// Get all players
players, err := client.GetPlayers()

// Get specific gameweek
gameweek, err := client.GetGameweek(1)

// Get league standings
league, err := client.GetLeague(123456)
```

## Error Handling
The wrapper provides detailed error types for different failure scenarios:
- API errors (rate limiting, server errors)
- Validation errors (invalid parameters)
- Network errors (timeout, connection issues)
- Parsing errors (invalid response format)

## Rate Limiting
The wrapper implements automatic rate limiting to prevent exceeding FPL's API limits:
- Configurable request rates
- Automatic request queuing
- Backoff strategies for rate limit errors

## Caching
Optional response caching is available to improve performance:
- In-memory cache with configurable TTL
- Cache invalidation on gameweek boundaries
- Selective caching for specific endpoint responses

## Contributing
Contributions are welcome! Please refer to our contributing guidelines for more information.

## Future Considerations
- Support for additional statistical analysis
- Enhanced caching strategies
- Additional convenience methods for common queries
- Performance optimizations for large data sets

## Limitations
- No support for team management operations (transfers, captain selection, etc.)
- Rate limiting may affect real-time data access
- Some data may be delayed based on FPL API updates
- Cache invalidation might need manual handling in some cases

## License
[MIT License](./LICENSE)


## Versioning
This project follows semantic versioning. Major version changes may include breaking changes to the API.
