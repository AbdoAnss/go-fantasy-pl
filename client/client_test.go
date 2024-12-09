package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewClient(
		WithTimeout(20*time.Second),
		WithRateLimit(30, time.Minute),
	)

	// Use Testify's assert package for better readability
	assert.NotNil(t, client, "Expected non-nil client")
	assert.Equal(t, baseURL, client.baseURL, "Expected baseURL to be %s, got %s", baseURL, client.baseURL)
}
