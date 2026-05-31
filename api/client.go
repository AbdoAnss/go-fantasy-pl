// Package api defines the core interfaces for the FPL client.
// This package is used to decouple services from the concrete client implementation,
// preventing circular dependencies and facilitating easier testing through mocking.
package api

import (
	"context"
	"net/http"
)

// Client is the interface that wraps the basic HTTP methods used by FPL services.
// It allows for different client implementations, such as the standard client
// or a mock client for testing.
type Client interface {
	// Get performs a GET request to the specified endpoint using the default client settings.
	Get(endpoint string) (*http.Response, error)

	// GetContext performs a GET request to the specified endpoint with the provided context.
	GetContext(ctx context.Context, endpoint string) (*http.Response, error)
}
