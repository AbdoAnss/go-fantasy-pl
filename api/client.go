package api

import (
	"context"
	"net/http"
)

type Client interface {
	Get(endpoint string) (*http.Response, error)
	GetContext(ctx context.Context, endpoint string) (*http.Response, error)
}
